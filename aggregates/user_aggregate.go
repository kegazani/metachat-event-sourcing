package aggregates

import (
	"errors"
	"time"

	"github.com/metachat/common/event-sourcing/events"
)

// UserAggregate represents the user aggregate
type UserAggregate struct {
	*BaseAggregate
	username    string
	email       string
	firstName   string
	lastName    string
	dateOfBirth string
	archetype   *events.Archetype
	modalities  []events.UserModality
}

// NewUserAggregate creates a new user aggregate
func NewUserAggregate(id string) *UserAggregate {
	return &UserAggregate{
		BaseAggregate: NewBaseAggregate(id),
		modalities:    make([]events.UserModality, 0),
	}
}

// CreateUser creates a new user
func (u *UserAggregate) CreateUser(username, email, firstName, lastName, dateOfBirth string) error {
	if u.username != "" {
		return errors.New("user already exists")
	}

	event, err := events.NewEvent(
		events.UserRegisteredEvent,
		u.GetID(),
		u.GetVersion()+1,
		events.UserRegisteredPayload{
			Username:    username,
			Email:       email,
			FirstName:   firstName,
			LastName:    lastName,
			DateOfBirth: dateOfBirth,
		},
		nil,
	)
	if err != nil {
		return err
	}

	u.AddUncommittedEvent(event)
	return nil
}

// UpdateProfile updates the user profile
func (u *UserAggregate) UpdateProfile(firstName, lastName, dateOfBirth, avatar, bio string) error {
	if u.username == "" {
		return errors.New("user does not exist")
	}

	event, err := events.NewEvent(
		events.UserProfileUpdatedEvent,
		u.GetID(),
		u.GetVersion()+1,
		events.UserProfileUpdatedPayload{
			FirstName:   firstName,
			LastName:    lastName,
			DateOfBirth: dateOfBirth,
			Avatar:      avatar,
			Bio:         bio,
		},
		nil,
	)
	if err != nil {
		return err
	}

	u.AddUncommittedEvent(event)
	return nil
}

// AssignArchetype assigns an archetype to the user
func (u *UserAggregate) AssignArchetype(archetypeID, archetypeName string, confidence float64, description string) error {
	if u.username == "" {
		return errors.New("user does not exist")
	}

	event, err := events.NewEvent(
		events.UserArchetypeAssignedEvent,
		u.GetID(),
		u.GetVersion()+1,
		events.UserArchetypeAssignedPayload{
			ArchetypeID:   archetypeID,
			ArchetypeName: archetypeName,
			Confidence:    confidence,
			Description:   description,
		},
		nil,
	)
	if err != nil {
		return err
	}

	u.AddUncommittedEvent(event)
	return nil
}

// UpdateArchetype updates the user archetype
func (u *UserAggregate) UpdateArchetype(archetypeID, archetypeName string, confidence float64, description string) error {
	if u.username == "" {
		return errors.New("user does not exist")
	}

	event, err := events.NewEvent(
		events.UserArchetypeUpdatedEvent,
		u.GetID(),
		u.GetVersion()+1,
		events.UserArchetypeUpdatedPayload{
			ArchetypeID:   archetypeID,
			ArchetypeName: archetypeName,
			Confidence:    confidence,
			Description:   description,
		},
		nil,
	)
	if err != nil {
		return err
	}

	u.AddUncommittedEvent(event)
	return nil
}

// UpdateModalities updates the user modalities
func (u *UserAggregate) UpdateModalities(modalities []events.UserModality) error {
	if u.username == "" {
		return errors.New("user does not exist")
	}

	event, err := events.NewEvent(
		events.UserModalitiesUpdatedEvent,
		u.GetID(),
		u.GetVersion()+1,
		events.UserModalitiesUpdatedPayload{
			Modalities: modalities,
		},
		nil,
	)
	if err != nil {
		return err
	}

	u.AddUncommittedEvent(event)
	return nil
}

// ApplyEvent applies an event to the aggregate
func (u *UserAggregate) ApplyEvent(event *events.Event) error {
	switch event.Type {
	case events.UserRegisteredEvent:
		return u.applyUserRegistered(event)
	case events.UserProfileUpdatedEvent:
		return u.applyUserProfileUpdated(event)
	case events.UserArchetypeAssignedEvent:
		return u.applyUserArchetypeAssigned(event)
	case events.UserArchetypeUpdatedEvent:
		return u.applyUserArchetypeUpdated(event)
	case events.UserModalitiesUpdatedEvent:
		return u.applyUserModalitiesUpdated(event)
	default:
		return errors.New("unknown event type")
	}
}

// applyUserRegistered applies the UserRegistered event
func (u *UserAggregate) applyUserRegistered(event *events.Event) error {
	var payload events.UserRegisteredPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return err
	}

	u.username = payload.Username
	u.email = payload.Email
	u.firstName = payload.FirstName
	u.lastName = payload.LastName
	u.dateOfBirth = payload.DateOfBirth
	u.IncrementVersion()
	return nil
}

// applyUserProfileUpdated applies the UserProfileUpdated event
func (u *UserAggregate) applyUserProfileUpdated(event *events.Event) error {
	var payload events.UserProfileUpdatedPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return err
	}

	if payload.FirstName != "" {
		u.firstName = payload.FirstName
	}
	if payload.LastName != "" {
		u.lastName = payload.LastName
	}
	if payload.DateOfBirth != "" {
		u.dateOfBirth = payload.DateOfBirth
	}
	u.IncrementVersion()
	return nil
}

// applyUserArchetypeAssigned applies the UserArchetypeAssigned event
func (u *UserAggregate) applyUserArchetypeAssigned(event *events.Event) error {
	var payload events.UserArchetypeAssignedPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return err
	}

	u.archetype = &events.Archetype{
		ID:          payload.ArchetypeID,
		Name:        payload.ArchetypeName,
		Description: payload.Description,
		Score:       payload.Confidence,
	}
	u.IncrementVersion()
	return nil
}

// applyUserArchetypeUpdated applies the UserArchetypeUpdated event
func (u *UserAggregate) applyUserArchetypeUpdated(event *events.Event) error {
	var payload events.UserArchetypeUpdatedPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return err
	}

	u.archetype = &events.Archetype{
		ID:          payload.ArchetypeID,
		Name:        payload.ArchetypeName,
		Description: payload.Description,
		Score:       payload.Confidence,
	}
	u.IncrementVersion()
	return nil
}

// applyUserModalitiesUpdated applies the UserModalitiesUpdated event
func (u *UserAggregate) applyUserModalitiesUpdated(event *events.Event) error {
	var payload events.UserModalitiesUpdatedPayload
	if err := event.UnmarshalPayload(&payload); err != nil {
		return err
	}

	u.modalities = payload.Modalities
	u.IncrementVersion()
	return nil
}

// GetUsername returns the username
func (u *UserAggregate) GetUsername() string {
	return u.username
}

// GetEmail returns the email
func (u *UserAggregate) GetEmail() string {
	return u.email
}

// GetFullName returns the full name
func (u *UserAggregate) GetFullName() string {
	return u.firstName + " " + u.lastName
}

// GetArchetype returns the archetype
func (u *UserAggregate) GetArchetype() *events.Archetype {
	return u.archetype
}

// GetModalities returns the modalities
func (u *UserAggregate) GetModalities() []events.UserModality {
	return u.modalities
}

// GetFirstName returns the first name
func (u *UserAggregate) GetFirstName() string {
	return u.firstName
}

// GetLastName returns the last name
func (u *UserAggregate) GetLastName() string {
	return u.lastName
}

// GetDateOfBirth returns the date of birth
func (u *UserAggregate) GetDateOfBirth() string {
	return u.dateOfBirth
}

// GetAvatar returns the avatar (empty for now as it's not stored in the aggregate)
func (u *UserAggregate) GetAvatar() string {
	return ""
}

// GetBio returns the bio (empty for now as it's not stored in the aggregate)
func (u *UserAggregate) GetBio() string {
	return ""
}

// GetCreatedAt returns the created at timestamp (empty for now as it's not stored in the aggregate)
func (u *UserAggregate) GetCreatedAt() time.Time {
	return time.Time{}
}

// GetUpdatedAt returns the updated at timestamp (empty for now as it's not stored in the aggregate)
func (u *UserAggregate) GetUpdatedAt() time.Time {
	return time.Time{}
}
