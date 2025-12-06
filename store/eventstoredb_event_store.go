package store

import (
	"context"
	"encoding/json"
	"fmt"
	"net/url"
	"time"

	client "github.com/EventStore/EventStore-Client-Go/v3/esdb"
	"github.com/gofrs/uuid"
	"github.com/kegazani/metachat-event-sourcing/events"
)

type EventStoreDBEventStore struct {
	client *client.Client
	streamPrefix string
}

func NewEventStoreDBEventStore(connectionString string, streamPrefix string) (*EventStoreDBEventStore, error) {
	config, err := client.ParseConnectionString(connectionString)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeConnectionFailed, "failed to parse connection string", err)
	}

	dbClient, err := client.NewClient(config)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeConnectionFailed, "failed to create EventStoreDB client", err)
	}

	if streamPrefix == "" {
		streamPrefix = "metachat"
	}

	return &EventStoreDBEventStore{
		client: dbClient,
		streamPrefix: streamPrefix,
	}, nil
}

func NewEventStoreDBEventStoreFromConfig(eventStoreURL, username, password string, streamPrefix string) (*EventStoreDBEventStore, error) {
	parsedURL, err := url.Parse(eventStoreURL)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeConnectionFailed, "failed to parse URL", err)
	}

	connectionString := fmt.Sprintf("esdb://%s:%s@%s?tls=false", username, password, parsedURL.Host)
	return NewEventStoreDBEventStore(connectionString, streamPrefix)
}

func (e *EventStoreDBEventStore) getStreamName(aggregateID string) string {
	return fmt.Sprintf("%s-%s", e.streamPrefix, aggregateID)
}

func (e *EventStoreDBEventStore) SaveEvents(ctx context.Context, eventList []*events.Event) error {
	if len(eventList) == 0 {
		return nil
	}

	aggregateID := eventList[0].AggregateID
	if aggregateID == "" {
		return NewEventStoreError(ErrCodeSerialization, "aggregate ID cannot be empty", nil)
	}

	streamName := e.getStreamName(aggregateID)

	proposedEvents := make([]client.EventData, 0, len(eventList))
	
	for _, event := range eventList {
		if event.AggregateID != aggregateID {
			return NewEventStoreError(ErrCodeSerialization, "all events must have the same aggregate ID", nil)
		}

			payloadBytes := event.Payload

		metadata := map[string]interface{}{
			"type":        string(event.Type),
			"aggregate_id": event.AggregateID,
			"version":     event.Version,
			"timestamp":   event.Timestamp.Format(time.RFC3339),
		}

		if event.Metadata != nil {
			for k, v := range event.Metadata {
				metadata[k] = v
			}
		}

		metadataBytes, err := json.Marshal(metadata)
		if err != nil {
			return NewEventStoreError(ErrCodeSerialization, "failed to marshal metadata", err)
		}

		eventUUID, err := uuid.FromString(event.ID)
		if err != nil {
			return NewEventStoreError(ErrCodeSerialization, "failed to parse event ID as UUID", err)
		}

		eventData := client.EventData{
			EventID:     eventUUID,
			ContentType: client.ContentTypeJson,
			EventType:   string(event.Type),
			Data:        payloadBytes,
			Metadata:    metadataBytes,
		}

		proposedEvents = append(proposedEvents, eventData)
	}

	existingEvents, err := e.GetEventsByAggregateID(ctx, aggregateID)
	if err != nil {
		return NewEventStoreError(ErrCodeStorage, "failed to get existing events", err)
	}

	firstEventVersion := eventList[0].Version

	var expectedRevision client.ExpectedRevision
	if len(existingEvents) == 0 {
		if firstEventVersion != 1 {
			return ErrVersionConflict
		}
		expectedRevision = client.NoStream{}
	} else {
		lastVersion := existingEvents[len(existingEvents)-1].Version
		expectedLastVersion := firstEventVersion - len(eventList)
		if lastVersion != expectedLastVersion {
			return ErrVersionConflict
		}
		expectedRevision = client.Revision(uint64(lastVersion))
	}

	opts := client.AppendToStreamOptions{
		ExpectedRevision: expectedRevision,
	}

	_, err = e.client.AppendToStream(ctx, streamName, opts, proposedEvents...)
	if err != nil {
		if err.Error() == "wrong expected stream revision" {
			return ErrVersionConflict
		}
		return NewEventStoreError(ErrCodeStorage, "failed to append events to stream", err)
	}

	return nil
}

func (e *EventStoreDBEventStore) GetEventsByAggregateID(ctx context.Context, aggregateID string) ([]*events.Event, error) {
	if aggregateID == "" {
		return nil, NewEventStoreError(ErrCodeSerialization, "aggregate ID cannot be empty", nil)
	}

	streamName := e.getStreamName(aggregateID)

	stream, err := e.client.ReadStream(ctx, streamName, client.ReadStreamOptions{
		Direction: client.Forwards,
		From:      client.Start{},
	}, ^uint64(0))
	if err != nil {
		if err.Error() == "stream not found" {
			return []*events.Event{}, nil
		}
		return nil, NewEventStoreError(ErrCodeStorage, "failed to read stream", err)
	}
	defer stream.Close()

	var result []*events.Event
	for {
		event, err := stream.Recv()
		if err != nil {
			break
		}

		if event.Event == nil {
			continue
		}

		domainEvent, err := e.convertFromEventStoreEvent(event.Event)
		if err != nil {
			return nil, NewEventStoreError(ErrCodeSerialization, "failed to convert event", err)
		}

		result = append(result, domainEvent)
	}

	return result, nil
}

func (e *EventStoreDBEventStore) GetEventsByType(ctx context.Context, eventType events.EventType) ([]*events.Event, error) {
	stream, err := e.client.ReadAll(ctx, client.ReadAllOptions{
		Direction: client.Forwards,
		From:      client.Start{},
	}, ^uint64(0))
	if err != nil {
		return nil, NewEventStoreError(ErrCodeStorage, "failed to read all events", err)
	}
	defer stream.Close()

	var result []*events.Event
	for {
		event, err := stream.Recv()
		if err != nil {
			break
		}

		if event.Event == nil || event.Event.EventType != string(eventType) {
			continue
		}

		domainEvent, err := e.convertFromEventStoreEvent(event.Event)
		if err != nil {
			return nil, NewEventStoreError(ErrCodeSerialization, "failed to convert event", err)
		}

		result = append(result, domainEvent)
	}

	return result, nil
}

func (e *EventStoreDBEventStore) GetEventsByAggregateIDAndVersion(ctx context.Context, aggregateID string, version int) ([]*events.Event, error) {
	if aggregateID == "" {
		return nil, NewEventStoreError(ErrCodeSerialization, "aggregate ID cannot be empty", nil)
	}

	streamName := e.getStreamName(aggregateID)

	readOpts := client.ReadStreamOptions{
		Direction: client.Forwards,
		From:      client.Start{},
	}
	stream, err := e.client.ReadStream(ctx, streamName, readOpts, uint64(version+1))
	if err != nil {
		if err.Error() == "stream not found" {
			return []*events.Event{}, nil
		}
		return nil, NewEventStoreError(ErrCodeStorage, "failed to read stream", err)
	}
	defer stream.Close()

	var result []*events.Event
	for {
		event, err := stream.Recv()
		if err != nil {
			break
		}

		if event.Event == nil {
			continue
		}

		domainEvent, err := e.convertFromEventStoreEvent(event.Event)
		if err != nil {
			return nil, NewEventStoreError(ErrCodeSerialization, "failed to convert event", err)
		}

		if domainEvent.Version <= version {
			result = append(result, domainEvent)
		}
	}

	return result, nil
}

func (e *EventStoreDBEventStore) GetEventsByTimeRange(ctx context.Context, startTime, endTime string) ([]*events.Event, error) {
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeSerialization, "invalid start time", err)
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeSerialization, "invalid end time", err)
	}

	stream, err := e.client.ReadAll(ctx, client.ReadAllOptions{
		Direction: client.Forwards,
		From:      client.Start{},
	}, ^uint64(0))
	if err != nil {
		return nil, NewEventStoreError(ErrCodeStorage, "failed to read all events", err)
	}
	defer stream.Close()

	var result []*events.Event
	for {
		event, err := stream.Recv()
		if err != nil {
			break
		}

			if event.Event == nil {
			continue
		}

		domainEvent, err := e.convertFromEventStoreEvent(event.Event)
		if err != nil {
			return nil, NewEventStoreError(ErrCodeSerialization, "failed to convert event", err)
		}

		if domainEvent.Timestamp.After(start) && domainEvent.Timestamp.Before(end) {
			result = append(result, domainEvent)
		}
	}

	return result, nil
}

func (e *EventStoreDBEventStore) convertFromEventStoreEvent(event *client.RecordedEvent) (*events.Event, error) {
	var metadata map[string]interface{}
	if len(event.UserMetadata) > 0 {
		if err := json.Unmarshal(event.UserMetadata, &metadata); err != nil {
			return nil, err
		}
	}

	aggregateID, ok := metadata["aggregate_id"].(string)
	if !ok {
		return nil, fmt.Errorf("aggregate_id not found in metadata")
	}

	version, ok := metadata["version"].(float64)
	if !ok {
		version = float64(event.EventNumber)
	}

	timestampStr, ok := metadata["timestamp"].(string)
	var timestamp time.Time
	if ok {
		var err error
		timestamp, err = time.Parse(time.RFC3339, timestampStr)
		if err != nil {
			timestamp = event.CreatedDate
		}
	} else {
		timestamp = event.CreatedDate
	}

	eventType := events.EventType(event.EventType)

	delete(metadata, "type")
	delete(metadata, "aggregate_id")
	delete(metadata, "version")
	delete(metadata, "timestamp")

	return &events.Event{
		ID:          event.EventID.String(),
		Type:        eventType,
		AggregateID: aggregateID,
		Version:     int(version),
		Timestamp:   timestamp,
		Payload:     event.Data,
		Metadata:    metadata,
	}, nil
}

func (e *EventStoreDBEventStore) Close() error {
	return e.client.Close()
}

