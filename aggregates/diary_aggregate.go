package aggregates

import (
	"errors"
	"time"

	"github.com/kegazani/metachat-event-sourcing/events"
)

// DiaryAggregate represents the diary entry aggregate
type DiaryAggregate struct {
	*BaseAggregate
	userID     string
	title      string
	content    string
	tokenCount int
	sessionID  string
	tags       []string
	deleted    bool
}

// NewDiaryAggregate creates a new diary aggregate
func NewDiaryAggregate(id string) *DiaryAggregate {
	return &DiaryAggregate{
		BaseAggregate: NewBaseAggregate(id),
		tags:          make([]string, 0),
	}
}

// CreateEntry creates a new diary entry
func (d *DiaryAggregate) CreateEntry(userID, title, content string, tokenCount int, sessionID string, tags []string) error {
	if d.title != "" {
		return errors.New("diary entry already exists")
	}

	event, err := events.NewEvent(
		events.DiaryEntryCreatedEvent,
		d.GetID(),
		d.GetVersion()+1,
		events.DiaryEntryCreatedPayload{
			UserID:     userID,
			Title:      title,
			Content:    content,
			TokenCount: tokenCount,
			SessionID:  sessionID,
			Tags:       tags,
		},
		nil,
	)
	if err != nil {
		return err
	}

	d.AddUncommittedEvent(event)
	return nil
}

// UpdateEntry updates a diary entry
func (d *DiaryAggregate) UpdateEntry(title, content string, tokenCount int, tags []string) error {
	if d.deleted {
		return errors.New("diary entry has been deleted")
	}

	if d.title == "" {
		return errors.New("diary entry does not exist")
	}

	event, err := events.NewEvent(
		events.DiaryEntryUpdatedEvent,
		d.GetID(),
		d.GetVersion()+1,
		events.DiaryEntryUpdatedPayload{
			Title:      title,
			Content:    content,
			TokenCount: tokenCount,
			Tags:       tags,
		},
		nil,
	)
	if err != nil {
		return err
	}

	d.AddUncommittedEvent(event)
	return nil
}

// DeleteEntry deletes a diary entry
func (d *DiaryAggregate) DeleteEntry(reason string) error {
	if d.deleted {
		return errors.New("diary entry has already been deleted")
	}

	if d.title == "" {
		return errors.New("diary entry does not exist")
	}

	event, err := events.NewEvent(
		events.DiaryEntryDeletedEvent,
		d.GetID(),
		d.GetVersion()+1,
		events.DiaryEntryDeletedPayload{
			Reason: reason,
		},
		nil,
	)
	if err != nil {
		return err
	}

	d.AddUncommittedEvent(event)
	return nil
}

// ApplyEvent applies an event to the aggregate
func (d *DiaryAggregate) ApplyEvent(event *events.Event) error {
	switch event.Type {
	case events.DiaryEntryCreatedEvent:
		return d.applyDiaryEntryCreated(event)
	case events.DiaryEntryUpdatedEvent:
		return d.applyDiaryEntryUpdated(event)
	case events.DiaryEntryDeletedEvent:
		return d.applyDiaryEntryDeleted(event)
	default:
		return errors.New("unknown event type")
	}
}

// applyDiaryEntryCreated applies the DiaryEntryCreated event
func (d *DiaryAggregate) applyDiaryEntryCreated(event *events.Event) error {
	var payload events.DiaryEntryCreatedPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return err
	}

	d.userID = payload.UserID
	d.title = payload.Title
	d.content = payload.Content
	d.tokenCount = payload.TokenCount
	d.sessionID = payload.SessionID
	d.tags = payload.Tags
	d.deleted = false
	d.IncrementVersion()
	return nil
}

// applyDiaryEntryUpdated applies the DiaryEntryUpdated event
func (d *DiaryAggregate) applyDiaryEntryUpdated(event *events.Event) error {
	var payload events.DiaryEntryUpdatedPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return err
	}

	if payload.Title != "" {
		d.title = payload.Title
	}
	if payload.Content != "" {
		d.content = payload.Content
	}
	if payload.TokenCount > 0 {
		d.tokenCount = payload.TokenCount
	}
	if payload.Tags != nil {
		d.tags = payload.Tags
	}
	d.IncrementVersion()
	return nil
}

// applyDiaryEntryDeleted applies the DiaryEntryDeleted event
func (d *DiaryAggregate) applyDiaryEntryDeleted(event *events.Event) error {
	var payload events.DiaryEntryDeletedPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return err
	}

	d.deleted = true
	d.IncrementVersion()
	return nil
}

// GetUserID returns the user ID
func (d *DiaryAggregate) GetUserID() string {
	return d.userID
}

// GetTitle returns the title
func (d *DiaryAggregate) GetTitle() string {
	return d.title
}

// GetContent returns the content
func (d *DiaryAggregate) GetContent() string {
	return d.content
}

// GetTokenCount returns the token count
func (d *DiaryAggregate) GetTokenCount() int {
	return d.tokenCount
}

// GetSessionID returns the session ID
func (d *DiaryAggregate) GetSessionID() string {
	return d.sessionID
}

// GetTags returns the tags
func (d *DiaryAggregate) GetTags() []string {
	return d.tags
}

// IsDeleted returns whether the entry is deleted
func (d *DiaryAggregate) IsDeleted() bool {
	return d.deleted
}

// GetCreatedAt returns the creation time of the aggregate
func (d *DiaryAggregate) GetCreatedAt() time.Time {
	// For now, return a zero time since we don't track creation time in the aggregate
	// In a real implementation, this would be stored in the aggregate
	return time.Time{}
}

// GetUpdatedAt returns the last update time of the aggregate
func (d *DiaryAggregate) GetUpdatedAt() time.Time {
	// For now, return a zero time since we don't track update time in the aggregate
	// In a real implementation, this would be stored in the aggregate
	return time.Time{}
}

// DiarySession represents a diary session
type DiarySession struct {
	ID          string
	UserID      string
	Title       string
	Description string
	StartedAt   time.Time
	EndedAt     time.Time
	Status      string
	EntryCount  int
}

// DiaryAnalytics represents diary analytics
type DiaryAnalytics struct {
	TotalEntries             int64
	TotalSessions            int64
	ActiveUsers              int64
	AverageEntriesPerSession float64
	AverageSessionsPerUser   float64
	LastUpdated              time.Time
}
