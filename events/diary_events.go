package events

// DiaryEntryCreatedPayload represents the payload for DiaryEntryCreated event
type DiaryEntryCreatedPayload struct {
	UserID      string `json:"user_id"`
	Title       string `json:"title"`
	Content     string `json:"content"`
	TokenCount  int    `json:"token_count"`
	SessionID   string `json:"session_id"`
	Tags        []string `json:"tags,omitempty"`
}

// DiaryEntryUpdatedPayload represents the payload for DiaryEntryUpdated event
type DiaryEntryUpdatedPayload struct {
	Title       string `json:"title,omitempty"`
	Content     string `json:"content,omitempty"`
	TokenCount  int    `json:"token_count,omitempty"`
	Tags        []string `json:"tags,omitempty"`
}

// DiaryEntryDeletedPayload represents the payload for DiaryEntryDeleted event
type DiaryEntryDeletedPayload struct {
	Reason string `json:"reason,omitempty"`
}

// DiarySessionStartedPayload represents the payload for DiarySessionStarted event
type DiarySessionStartedPayload struct {
	UserID    string `json:"user_id"`
	StartTime string `json:"start_time"`
	Source    string `json:"source"` // "web", "mobile", etc.
}

// DiarySessionEndedPayload represents the payload for DiarySessionEnded event
type DiarySessionEndedPayload struct {
	SessionID   string `json:"session_id"`
	EndTime     string `json:"end_time"`
	EntryCount  int    `json:"entry_count"`
	TokenCount  int    `json:"token_count"`
}