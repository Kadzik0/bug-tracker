package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kadzik0/bug-tracker/internal/model"
)

type EventRepo struct {
	pool *pgxpool.Pool
}

func NewEventRepo(pool *pgxpool.Pool) *EventRepo {
	return &EventRepo{pool: pool}
}

func (r *EventRepo) Create(ctx context.Context, event *model.IncidentEvent) error {
	var sqlQueryEvent string
	sqlQueryEvent = `INSERT INTO incident_events (id, incident_id, actor_id, type, description, created_at) VALUES ($1, $2, $3, $4, $5, $6)`

	_, err := r.pool.Exec(ctx, sqlQueryEvent, event.ID, event.IncidentID, event.ActorID, event.EventType, event.Description, event.CreatedAt)
	return err
}

func (r *EventRepo) ListEvents(ctx context.Context, incidentID uuid.UUID) ([]*model.IncidentEvent, error) {
	var sqlQuery string
	sqlQuery = `SELECT id, incident_id, actor_id, type, description, created_at FROM incident_events WHERE incident_id = $1 ORDER BY created_at ASC`

	rows, err := r.pool.Query(ctx, sqlQuery, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var events []*model.IncidentEvent
	for rows.Next() {
		var event model.IncidentEvent
		if err := rows.Scan(&event.ID, &event.IncidentID, &event.ActorID, &event.EventType, &event.Description, &event.CreatedAt); err != nil {
			return nil, err
		}
		events = append(events, &event)
	}

	return events, nil
}
