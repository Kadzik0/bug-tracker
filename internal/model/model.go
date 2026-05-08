package model

import (
	"time"

	"github.com/google/uuid"
)

type (
	User struct {
		ID        uuid.UUID `json:"id"`
		Email     string    `json:"email"`
		Name      string    `json:"name"`
		CreatedAt time.Time `json:"created_at"`
	}

	Incident struct {
		ID          uuid.UUID   `json:"id"`
		Title       string      `json:"title"`
		Description string      `json:"description"`
		Environment Environment `json:"environment"`
		Priority    Priority    `json:"priority"`
		Status      Status      `json:"status"`
		ReporterID  uuid.UUID   `json:"reporter_id"`
		AssigneeID  *uuid.UUID  `json:"assignee_id"`
		CreatedAt   time.Time   `json:"created_at"`
		UpdatedAt   time.Time   `json:"updated_at"`
	}

	Environment string

	Priority string

	Status string

	IncidentComment struct {
		ID         uuid.UUID `json:"id"`
		IncidentID uuid.UUID `json:"incident_id"`
		AuthorID   uuid.UUID `json:"author_id"`
		Body       string    `json:"body"`
		CreatedAt  time.Time `json:"created_at"`
	}

	IncidentEvent struct {
		ID          uuid.UUID `json:"id"`
		IncidentID  uuid.UUID `json:"incident_id"`
		ActorID     uuid.UUID `json:"actor_id"`
		EventType   EventType `json:"type"`
		Description string    `json:"description"`
		CreatedAt   time.Time `json:"created_at"`
	}

	EventType string
)

const (
	EnvironmentProd  Environment = "prod"
	EnvironmentStage Environment = "stage"
	EnvironmentDev   Environment = "dev"

	PriorityP1 Priority = "P1"
	PriorityP2 Priority = "P2"
	PriorityP3 Priority = "P3"
	PriorityP4 Priority = "P4"

	StatusOpen       Status = "OPEN"
	StatusInProgress Status = "IN_PROGRESS"
	StatusResolved   Status = "RESOLVED"
	StatusClosed     Status = "CLOSED"

	EventTypeCreated         EventType = "CREATED"
	EventTypeStatusChanged   EventType = "STATUS_CHANGED"
	EventTypeAssigned        EventType = "ASSIGNED"
	EventTypePriorityChanged EventType = "PRIORITY_CHANGED"
	EventTypeCommentAdded    EventType = "COMMENT_ADDED"
)
