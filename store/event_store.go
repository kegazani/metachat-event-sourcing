package store

import (
	"context"
	"errors"

	"github.com/kegazani/metachat-event-sourcing/events"
)

// EventStore defines the interface for event storage
type EventStore interface {
	// SaveEvents saves a batch of events to the store
	SaveEvents(ctx context.Context, events []*events.Event) error

	// GetEventsByAggregateID retrieves all events for a specific aggregate
	GetEventsByAggregateID(ctx context.Context, aggregateID string) ([]*events.Event, error)

	// GetEventsByType retrieves all events of a specific type
	GetEventsByType(ctx context.Context, eventType events.EventType) ([]*events.Event, error)

	// GetEventsByAggregateIDAndVersion retrieves events for an aggregate up to a specific version
	GetEventsByAggregateIDAndVersion(ctx context.Context, aggregateID string, version int) ([]*events.Event, error)

	// GetEventsByTimeRange retrieves events within a time range
	GetEventsByTimeRange(ctx context.Context, startTime, endTime string) ([]*events.Event, error)
}

// EventStoreError represents an error from the event store
type EventStoreError struct {
	Code    string
	Message string
	Err     error
}

// Error implements the error interface
func (e *EventStoreError) Error() string {
	if e.Err != nil {
		return e.Message + ": " + e.Err.Error()
	}
	return e.Message
}

// Unwrap returns the underlying error
func (e *EventStoreError) Unwrap() error {
	return e.Err
}

// NewEventStoreError creates a new event store error
func NewEventStoreError(code, message string, err error) *EventStoreError {
	return &EventStoreError{
		Code:    code,
		Message: message,
		Err:     err,
	}
}

// Predefined error codes
const (
	ErrCodeConnectionFailed = "CONNECTION_FAILED"
	ErrCodeEventNotFound    = "EVENT_NOT_FOUND"
	ErrCodeVersionConflict  = "VERSION_CONFLICT"
	ErrCodeSerialization    = "SERIALIZATION_ERROR"
	ErrCodeStorage          = "STORAGE_ERROR"
)

// Predefined errors
var (
	ErrConnectionFailed = NewEventStoreError(ErrCodeConnectionFailed, "failed to connect to event store", nil)
	ErrEventNotFound    = NewEventStoreError(ErrCodeEventNotFound, "event not found", nil)
	ErrVersionConflict  = NewEventStoreError(ErrCodeVersionConflict, "version conflict", nil)
	ErrSerialization    = NewEventStoreError(ErrCodeSerialization, "serialization error", nil)
	ErrStorage          = NewEventStoreError(ErrCodeStorage, "storage error", nil)
)

// IsEventStoreError checks if an error is an EventStoreError
func IsEventStoreError(err error) bool {
	var e *EventStoreError
	return errors.As(err, &e)
}

// GetEventStoreErrorCode returns the error code from an EventStoreError
func GetEventStoreErrorCode(err error) string {
	var e *EventStoreError
	if errors.As(err, &e) {
		return e.Code
	}
	return ""
}
