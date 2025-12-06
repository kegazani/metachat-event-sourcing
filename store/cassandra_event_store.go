package store

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/gocql/gocql"
	"github.com/google/uuid"
	"github.com/kegazani/metachat-event-sourcing/events"
)

type CassandraEventStore struct {
	session *gocql.Session
}

func NewCassandraEventStore(session *gocql.Session) *CassandraEventStore {
	return &CassandraEventStore{
		session: session,
	}
}

func (c *CassandraEventStore) InitializeSchema(keyspace string) error {
	queries := []string{
		fmt.Sprintf(`CREATE KEYSPACE IF NOT EXISTS %s WITH replication = {'class': 'SimpleStrategy', 'replication_factor': 1}`, keyspace),
		fmt.Sprintf(`USE %s`, keyspace),
		`CREATE TABLE IF NOT EXISTS events (
			aggregate_type text,
			aggregate_id uuid,
			version int,
			event_id uuid,
			event_type text,
			payload text,
			metadata text,
			created_at timestamp,
			PRIMARY KEY ((aggregate_type, aggregate_id), version)
		) WITH CLUSTERING ORDER BY (version ASC)`,
		`CREATE INDEX IF NOT EXISTS ON events (event_type)`,
		`CREATE INDEX IF NOT EXISTS ON events (created_at)`,
	}

	for _, query := range queries {
		if err := c.session.Query(query).Exec(); err != nil {
			return fmt.Errorf("failed to execute schema query: %w", err)
		}
	}

	return nil
}

func (c *CassandraEventStore) SaveEvents(ctx context.Context, eventList []*events.Event) error {
	batch := c.session.NewBatch(gocql.LoggedBatch)

	for _, event := range eventList {
		aggregateType := c.getAggregateType(event.AggregateID)

		payloadJSON, err := json.Marshal(event.Payload)
		if err != nil {
			return NewEventStoreError(ErrCodeSerialization, "failed to marshal payload", err)
		}

		metadataJSON, err := json.Marshal(event.Metadata)
		if err != nil {
			return NewEventStoreError(ErrCodeSerialization, "failed to marshal metadata", err)
		}

		eventID, err := uuid.Parse(event.ID)
		if err != nil {
			eventID = uuid.New()
		}

		aggregateID, err := uuid.Parse(event.AggregateID)
		if err != nil {
			return NewEventStoreError(ErrCodeSerialization, "invalid aggregate ID", err)
		}

		batch.Query(
			`INSERT INTO events (aggregate_type, aggregate_id, version, event_id, event_type, payload, metadata, created_at)
			 VALUES (?, ?, ?, ?, ?, ?, ?, ?)`,
			aggregateType,
			aggregateID,
			event.Version,
			eventID,
			string(event.Type),
			string(payloadJSON),
			string(metadataJSON),
			event.Timestamp,
		)
	}

	if err := c.session.ExecuteBatch(batch); err != nil {
		return NewEventStoreError(ErrCodeStorage, "failed to save events", err)
	}

	return nil
}

func (c *CassandraEventStore) GetEventsByAggregateID(ctx context.Context, aggregateID string) ([]*events.Event, error) {
	aggregateType := c.getAggregateType(aggregateID)

	aggregateUUID, err := uuid.Parse(aggregateID)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeSerialization, "invalid aggregate ID", err)
	}

	iter := c.session.Query(
		`SELECT event_id, event_type, payload, metadata, created_at, version
		 FROM events
		 WHERE aggregate_type = ? AND aggregate_id = ?`,
		aggregateType,
		aggregateUUID,
	).Iter()

	var eventList []*events.Event
	var eventID, eventType, payload, metadata string
	var createdAt time.Time
	var version int

	for iter.Scan(&eventID, &eventType, &payload, &metadata, &createdAt, &version) {
		var metadataMap map[string]interface{}
		if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
			metadataMap = make(map[string]interface{})
		}

		event := &events.Event{
			ID:          eventID,
			Type:        events.EventType(eventType),
			AggregateID: aggregateID,
			Version:     version,
			Timestamp:   createdAt,
			Payload:     json.RawMessage(payload),
			Metadata:    metadataMap,
		}

		eventList = append(eventList, event)
	}

	if err := iter.Close(); err != nil {
		return nil, NewEventStoreError(ErrCodeStorage, "failed to retrieve events", err)
	}

	return eventList, nil
}

func (c *CassandraEventStore) GetEventsByType(ctx context.Context, eventType events.EventType) ([]*events.Event, error) {
	iter := c.session.Query(
		`SELECT aggregate_type, aggregate_id, event_id, event_type, payload, metadata, created_at, version
		 FROM events
		 WHERE event_type = ?`,
		string(eventType),
	).Iter()

	var eventList []*events.Event
	var aggregateType, aggregateID, eventID, eventTypeStr, payload, metadata string
	var createdAt time.Time
	var version int

	for iter.Scan(&aggregateType, &aggregateID, &eventID, &eventTypeStr, &payload, &metadata, &createdAt, &version) {
		var metadataMap map[string]interface{}
		if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
			metadataMap = make(map[string]interface{})
		}

		event := &events.Event{
			ID:          eventID,
			Type:        events.EventType(eventTypeStr),
			AggregateID: aggregateID,
			Version:     version,
			Timestamp:   createdAt,
			Payload:     json.RawMessage(payload),
			Metadata:    metadataMap,
		}

		eventList = append(eventList, event)
	}

	if err := iter.Close(); err != nil {
		return nil, NewEventStoreError(ErrCodeStorage, "failed to retrieve events", err)
	}

	return eventList, nil
}

func (c *CassandraEventStore) GetEventsByAggregateIDAndVersion(ctx context.Context, aggregateID string, version int) ([]*events.Event, error) {
	aggregateType := c.getAggregateType(aggregateID)

	aggregateUUID, err := uuid.Parse(aggregateID)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeSerialization, "invalid aggregate ID", err)
	}

	iter := c.session.Query(
		`SELECT event_id, event_type, payload, metadata, created_at, version
		 FROM events
		 WHERE aggregate_type = ? AND aggregate_id = ? AND version <= ?`,
		aggregateType,
		aggregateUUID,
		version,
	).Iter()

	var eventList []*events.Event
	var eventID, eventType, payload, metadata string
	var createdAt time.Time
	var eventVersion int

	for iter.Scan(&eventID, &eventType, &payload, &metadata, &createdAt, &eventVersion) {
		var metadataMap map[string]interface{}
		if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
			metadataMap = make(map[string]interface{})
		}

		event := &events.Event{
			ID:          eventID,
			Type:        events.EventType(eventType),
			AggregateID: aggregateID,
			Version:     eventVersion,
			Timestamp:   createdAt,
			Payload:     json.RawMessage(payload),
			Metadata:    metadataMap,
		}

		eventList = append(eventList, event)
	}

	if err := iter.Close(); err != nil {
		return nil, NewEventStoreError(ErrCodeStorage, "failed to retrieve events", err)
	}

	return eventList, nil
}

func (c *CassandraEventStore) GetEventsByTimeRange(ctx context.Context, startTime, endTime string) ([]*events.Event, error) {
	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeSerialization, "invalid start time", err)
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, NewEventStoreError(ErrCodeSerialization, "invalid end time", err)
	}

	iter := c.session.Query(
		`SELECT aggregate_type, aggregate_id, event_id, event_type, payload, metadata, created_at, version
		 FROM events
		 WHERE created_at >= ? AND created_at <= ?`,
		start,
		end,
	).Iter()

	var eventList []*events.Event
	var aggregateType, aggregateID, eventID, eventType, payload, metadata string
	var createdAt time.Time
	var version int

	for iter.Scan(&aggregateType, &aggregateID, &eventID, &eventType, &payload, &metadata, &createdAt, &version) {
		var metadataMap map[string]interface{}
		if err := json.Unmarshal([]byte(metadata), &metadataMap); err != nil {
			metadataMap = make(map[string]interface{})
		}

		event := &events.Event{
			ID:          eventID,
			Type:        events.EventType(eventType),
			AggregateID: aggregateID,
			Version:     version,
			Timestamp:   createdAt,
			Payload:     json.RawMessage(payload),
			Metadata:    metadataMap,
		}

		eventList = append(eventList, event)
	}

	if err := iter.Close(); err != nil {
		return nil, NewEventStoreError(ErrCodeStorage, "failed to retrieve events", err)
	}

	return eventList, nil
}

func (c *CassandraEventStore) getAggregateType(aggregateID string) string {
	return "default"
}

