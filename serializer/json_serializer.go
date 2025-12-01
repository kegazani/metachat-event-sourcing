package serializer

import (
	"encoding/json"

	"github.com/kegazani/metachat-event-sourcing/events"
)

// Serializer defines the interface for event serialization
type Serializer interface {
	// Serialize serializes an event to bytes
	Serialize(event *events.Event) ([]byte, error)

	// Deserialize deserializes bytes to an event
	Deserialize(data []byte) (*events.Event, error)
}

// JSONSerializer is a JSON implementation of Serializer
type JSONSerializer struct{}

// NewJSONSerializer creates a new JSON serializer
func NewJSONSerializer() *JSONSerializer {
	return &JSONSerializer{}
}

// Serialize serializes an event to JSON bytes
func (j *JSONSerializer) Serialize(event *events.Event) ([]byte, error) {
	return json.Marshal(event)
}

// Deserialize deserializes JSON bytes to an event
func (j *JSONSerializer) Deserialize(data []byte) (*events.Event, error) {
	var event events.Event
	if err := json.Unmarshal(data, &event); err != nil {
		return nil, err
	}
	return &event, nil
}
