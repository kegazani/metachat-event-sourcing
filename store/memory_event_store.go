package store

import (
	"context"
	"sort"
	"sync"
	"time"

	"github.com/kegazani/metachat-event-sourcing/events"
)

// MemoryEventStore is an in-memory implementation of EventStore
// This is mainly for testing and development purposes
type MemoryEventStore struct {
	mu     sync.RWMutex
	events []*events.Event
	index  map[string][]int // aggregateID -> event indices
}

// NewMemoryEventStore creates a new in-memory event store
func NewMemoryEventStore() *MemoryEventStore {
	return &MemoryEventStore{
		events: make([]*events.Event, 0),
		index:  make(map[string][]int),
	}
}

// SaveEvents saves a batch of events to the store
func (m *MemoryEventStore) SaveEvents(ctx context.Context, eventList []*events.Event) error {
	m.mu.Lock()
	defer m.mu.Unlock()

	for _, event := range eventList {
		// Check for version conflicts
		for _, existingEvent := range m.events {
			if existingEvent.AggregateID == event.AggregateID && existingEvent.Version == event.Version {
				return ErrVersionConflict
			}
		}

		// Add event to the store
		m.events = append(m.events, event)

		// Update index
		m.index[event.AggregateID] = append(m.index[event.AggregateID], len(m.events)-1)
	}

	return nil
}

// GetEventsByAggregateID retrieves all events for a specific aggregate
func (m *MemoryEventStore) GetEventsByAggregateID(ctx context.Context, aggregateID string) ([]*events.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	indices, ok := m.index[aggregateID]
	if !ok {
		return []*events.Event{}, nil
	}

	result := make([]*events.Event, 0, len(indices))
	for _, idx := range indices {
		result = append(result, m.events[idx])
	}

	// Sort events by version
	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil
}

// GetEventsByType retrieves all events of a specific type
func (m *MemoryEventStore) GetEventsByType(ctx context.Context, eventType events.EventType) ([]*events.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	result := make([]*events.Event, 0)
	for _, event := range m.events {
		if event.Type == eventType {
			result = append(result, event)
		}
	}

	return result, nil
}

// GetEventsByAggregateIDAndVersion retrieves events for an aggregate up to a specific version
func (m *MemoryEventStore) GetEventsByAggregateIDAndVersion(ctx context.Context, aggregateID string, version int) ([]*events.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	indices, ok := m.index[aggregateID]
	if !ok {
		return []*events.Event{}, nil
	}

	result := make([]*events.Event, 0, len(indices))
	for _, idx := range indices {
		event := m.events[idx]
		if event.Version <= version {
			result = append(result, event)
		}
	}

	// Sort events by version
	sort.Slice(result, func(i, j int) bool {
		return result[i].Version < result[j].Version
	})

	return result, nil
}

// GetEventsByTimeRange retrieves events within a time range
func (m *MemoryEventStore) GetEventsByTimeRange(ctx context.Context, startTime, endTime string) ([]*events.Event, error) {
	m.mu.RLock()
	defer m.mu.RUnlock()

	start, err := time.Parse(time.RFC3339, startTime)
	if err != nil {
		return nil, err
	}

	end, err := time.Parse(time.RFC3339, endTime)
	if err != nil {
		return nil, err
	}

	result := make([]*events.Event, 0)
	for _, event := range m.events {
		if event.Timestamp.After(start) && event.Timestamp.Before(end) {
			result = append(result, event)
		}
	}

	return result, nil
}

// Clear clears all events from the store (mainly for testing)
func (m *MemoryEventStore) Clear() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.events = make([]*events.Event, 0)
	m.index = make(map[string][]int)
}

// Len returns the number of events in the store (mainly for testing)
func (m *MemoryEventStore) Len() int {
	m.mu.RLock()
	defer m.mu.RUnlock()

	return len(m.events)
}
