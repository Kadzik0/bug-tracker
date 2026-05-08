package repo

import (
	"context"
	"fmt"
	"strings"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kadzik0/bug-tracker/internal/model"
)

type IncidentRepo struct {
	pool *pgxpool.Pool
}

type ListIncidentsOptions struct {
	Status      *string
	Priority    *string
	Environment *string
	AssigneeID  *string
	Q           *string
}

func NewIncidentRepo(pool *pgxpool.Pool) *IncidentRepo {
	return &IncidentRepo{pool: pool}
}

func (r *IncidentRepo) Create(ctx context.Context, incident *model.Incident) error {
	var sqlQuery string
	sqlQuery = `INSERT INTO incidents (id, title, description, environment, priority, status, reporter_id, assignee_id, created_at, updated_at) VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9, $10)`

	_, err := r.pool.Exec(ctx, sqlQuery, incident.ID, incident.Title, incident.Description, incident.Environment, incident.Priority, incident.Status, incident.ReporterID, incident.AssigneeID, incident.CreatedAt, incident.UpdatedAt)
	return err
}

func (r *IncidentRepo) List(ctx context.Context, opts *ListIncidentsOptions) ([]*model.Incident, error) {
	var conditions []string
	var args []any
	i := 1

	if opts.Status != nil {
		conditions = append(conditions, fmt.Sprintf("status = $%d", i))
		args = append(args, *opts.Status)
		i++
	}
	if opts.Priority != nil {
		conditions = append(conditions, fmt.Sprintf("priority = $%d", i))
		args = append(args, *opts.Priority)
		i++
	}
	if opts.Environment != nil {
		conditions = append(conditions, fmt.Sprintf("environment = $%d", i))
		args = append(args, *opts.Environment)
		i++
	}
	if opts.AssigneeID != nil {
		conditions = append(conditions, fmt.Sprintf("assignee_id = $%d", i))
		args = append(args, *opts.AssigneeID)
		i++
	}
	if opts.Q != nil {
		conditions = append(conditions, fmt.Sprintf("title ILIKE $%d", i))
		args = append(args, "%"+*opts.Q+"%")
		i++
	}

	query := `SELECT id, title, description, environment, priority, status, reporter_id, assignee_id, created_at, updated_at FROM incidents`
	if len(conditions) > 0 {
		query += " WHERE " + strings.Join(conditions, " AND ")
	}

	result, err := r.pool.Query(ctx, query, args...)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var incidents []*model.Incident
	for result.Next() {
		var incident model.Incident
		if err := result.Scan(&incident.ID, &incident.Title, &incident.Description, &incident.Environment, &incident.Priority, &incident.Status, &incident.ReporterID, &incident.AssigneeID, &incident.CreatedAt, &incident.UpdatedAt); err != nil {
			return nil, err
		}
		incidents = append(incidents, &incident)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	return incidents, nil
}

func (r *IncidentRepo) GetByID(ctx context.Context, id uuid.UUID) (*model.Incident, error) {
	var sqlQuery string
	sqlQuery = `SELECT id, title, description, environment, priority, status, reporter_id, assignee_id, created_at, updated_at FROM incidents WHERE id = $1`

	var incident model.Incident
	if err := r.pool.QueryRow(ctx, sqlQuery, id).Scan(&incident.ID, &incident.Title, &incident.Description, &incident.Environment, &incident.Priority, &incident.Status, &incident.ReporterID, &incident.AssigneeID, &incident.CreatedAt, &incident.UpdatedAt); err != nil {
		return nil, err
	}
	return &incident, nil
}

func (r *IncidentRepo) Update(ctx context.Context, incident *model.Incident) error {
	var sqlQuery string
	sqlQuery = `UPDATE incidents SET title = $1, description = $2, environment = $3, priority = $4, status = $5, reporter_id = $6, assignee_id = $7, updated_at = $8 WHERE id = $9`

	_, err := r.pool.Exec(ctx, sqlQuery, incident.Title, incident.Description, incident.Environment, incident.Priority, incident.Status, incident.ReporterID, &incident.AssigneeID, incident.UpdatedAt, incident.ID)
	return err
}
