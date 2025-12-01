package aggregates

import (
	"errors"

	"github.com/kegazani/metachat-event-sourcing/events"
)

// Aggregate represents a domain aggregate
type Aggregate interface {
	GetID() string
	GetVersion() int
	IncrementVersion()
	ApplyEvent(event *events.Event) error
	GetUncommittedEvents() []*events.Event
	ClearUncommittedEvents()
}

// BaseAggregate provides common functionality for aggregates
type BaseAggregate struct {
	id                string
	version           int
	uncommittedEvents []*events.Event
}

// NewBaseAggregate creates a new base aggregate
func NewBaseAggregate(id string) *BaseAggregate {
	return &BaseAggregate{
		id:                id,
		version:           0,
		uncommittedEvents: make([]*events.Event, 0),
	}
}

// GetID returns the aggregate ID
func (a *BaseAggregate) GetID() string {
	return a.id
}

// GetVersion returns the aggregate version
func (a *BaseAggregate) GetVersion() int {
	return a.version
}

// IncrementVersion increments the aggregate version
func (a *BaseAggregate) IncrementVersion() {
	a.version++
}

// AddUncommittedEvent adds an event to the uncommitted events list
func (a *BaseAggregate) AddUncommittedEvent(event *events.Event) {
	a.uncommittedEvents = append(a.uncommittedEvents, event)
}

// GetUncommittedEvents returns the uncommitted events
func (a *BaseAggregate) GetUncommittedEvents() []*events.Event {
	return a.uncommittedEvents
}

// ClearUncommittedEvents clears the uncommitted events
func (a *BaseAggregate) ClearUncommittedEvents() {
	a.uncommittedEvents = make([]*events.Event, 0)
}

// LoadFromHistory loads the aggregate state from a history of events
func (a *BaseAggregate) LoadFromHistory(eventList []*events.Event) error {
	for _, event := range eventList {
		if event.AggregateID != a.id {
			return errors.New("event aggregate ID does not match")
		}

		// Base aggregate doesn't apply events directly
		// This should be implemented by specific aggregates

		a.version = event.Version
	}

	return nil
}
