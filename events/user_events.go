package events

// UserRegisteredPayload represents the payload for UserRegistered event
type UserRegisteredPayload struct {
	Username    string `json:"username"`
	Email       string `json:"email"`
	FirstName   string `json:"first_name"`
	LastName    string `json:"last_name"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
}

// UserProfileUpdatedPayload represents the payload for UserProfileUpdated event
type UserProfileUpdatedPayload struct {
	FirstName   string `json:"first_name,omitempty"`
	LastName    string `json:"last_name,omitempty"`
	DateOfBirth string `json:"date_of_birth,omitempty"`
	Avatar      string `json:"avatar,omitempty"`
	Bio         string `json:"bio,omitempty"`
}

// UserArchetypeAssignedPayload represents the payload for UserArchetypeAssigned event
type UserArchetypeAssignedPayload struct {
	ArchetypeID    string  `json:"archetype_id"`
	ArchetypeName  string  `json:"archetype_name"`
	Confidence     float64 `json:"confidence"`
	Description    string  `json:"description"`
}

// UserArchetypeUpdatedPayload represents the payload for UserArchetypeUpdated event
type UserArchetypeUpdatedPayload struct {
	ArchetypeID    string  `json:"archetype_id"`
	ArchetypeName  string  `json:"archetype_name"`
	Confidence     float64 `json:"confidence"`
	Description    string  `json:"description"`
}

// UserModalitiesUpdatedPayload represents the payload for UserModalitiesUpdated event
type UserModalitiesUpdatedPayload struct {
	Modalities []UserModality `json:"modalities"`
}

// UserModality represents a user modality
type UserModality struct {
	ID          string  `json:"id"`
	Name        string  `json:"name"`
	Type        string  `json:"type"`
	Enabled     bool    `json:"enabled"`
	Weight      float64 `json:"weight"`
	Config      map[string]interface{} `json:"config,omitempty"`
}