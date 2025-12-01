package events

import (
	"encoding/json"
	"time"

	"github.com/google/uuid"
)

// EventType represents the type of event
type EventType string

const (
	// User events
	UserRegisteredEvent        EventType = "UserRegistered"
	UserProfileUpdatedEvent    EventType = "UserProfileUpdated"
	UserArchetypeAssignedEvent EventType = "UserArchetypeAssigned"
	UserArchetypeUpdatedEvent  EventType = "UserArchetypeUpdated"
	UserModalitiesUpdatedEvent EventType = "UserModalitiesUpdated"

	// Diary events
	DiaryEntryCreatedEvent   EventType = "DiaryEntryCreated"
	DiaryEntryUpdatedEvent   EventType = "DiaryEntryUpdated"
	DiaryEntryDeletedEvent   EventType = "DiaryEntryDeleted"
	DiarySessionStartedEvent EventType = "DiarySessionStarted"
	DiarySessionEndedEvent   EventType = "DiarySessionEnded"

	// Mood analysis events
	MoodAnalyzedEvent EventType = "MoodAnalyzed"

	// Aggregation events
	DailyMoodAggregatedEvent   EventType = "DailyMoodAggregated"
	WeeklyMoodAggregatedEvent  EventType = "WeeklyMoodAggregated"
	MonthlyMoodAggregatedEvent EventType = "MonthlyMoodAggregated"

	// Archetype events
	ArchetypeCalculationTriggeredEvent EventType = "ArchetypeCalculationTriggered"
	ArchetypeAssignedEvent             EventType = "ArchetypeAssigned"
	ArchetypeUpdatedEvent              EventType = "ArchetypeUpdated"

	// User portrait events
	UserPortraitUpdatedEvent EventType = "UserPortraitUpdated"
)

// Event represents a domain event
type Event struct {
	ID          string                 `json:"id"`
	Type        EventType              `json:"type"`
	AggregateID string                 `json:"aggregate_id"`
	Version     int                    `json:"version"`
	Timestamp   time.Time              `json:"timestamp"`
	Payload     json.RawMessage        `json:"payload"`
	Metadata    map[string]interface{} `json:"metadata"`
}

// NewEvent creates a new event
func NewEvent(eventType EventType, aggregateID string, version int, payload interface{}, metadata map[string]interface{}) (*Event, error) {
	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return nil, err
	}

	if metadata == nil {
		metadata = make(map[string]interface{})
	}

	return &Event{
		ID:          uuid.New().String(),
		Type:        eventType,
		AggregateID: aggregateID,
		Version:     version,
		Timestamp:   time.Now(),
		Payload:     payloadBytes,
		Metadata:    metadata,
	}, nil
}

// UnmarshalPayload unmarshals the event payload to the provided type
func (e *Event) UnmarshalPayload(v interface{}) error {
	return json.Unmarshal(e.Payload, v)
}

// EventMetadata represents metadata for an event
type EventMetadata struct {
	CorrelationID string                 `json:"correlation_id"`
	CausationID   string                 `json:"causation_id"`
	UserID        string                 `json:"user_id,omitempty"`
	Extra         map[string]interface{} `json:"extra,omitempty"`
}

// NewEventMetadata creates new event metadata
func NewEventMetadata(correlationID, causationID, userID string) *EventMetadata {
	return &EventMetadata{
		CorrelationID: correlationID,
		CausationID:   causationID,
		UserID:        userID,
		Extra:         make(map[string]interface{}),
	}
}

// ToMap converts metadata to map
func (m *EventMetadata) ToMap() map[string]interface{} {
	result := map[string]interface{}{
		"correlation_id": m.CorrelationID,
		"causation_id":   m.CausationID,
	}

	if m.UserID != "" {
		result["user_id"] = m.UserID
	}

	for k, v := range m.Extra {
		result[k] = v
	}

	return result
}
