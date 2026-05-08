package repo

import (
	"context"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kadzik0/bug-tracker/internal/model"
)

type CommentRepo struct {
	pool *pgxpool.Pool
}

func NewCommentRepo(pool *pgxpool.Pool) *CommentRepo {
	return &CommentRepo{pool: pool}
}

func (r *CommentRepo) Create(ctx context.Context, comment *model.IncidentComment) error {
	var sqlQueryComment string
	sqlQueryComment = `INSERT INTO incident_comments (id, incident_id, author_id, body, created_at) VALUES ($1, $2, $3, $4, $5)`

	_, err := r.pool.Exec(ctx, sqlQueryComment, comment.ID, comment.IncidentID, comment.AuthorID, comment.Body, comment.CreatedAt)

	return err
}

func (r *CommentRepo) ListComments(ctx context.Context, incidentID uuid.UUID) ([]*model.IncidentComment, error) {
	var sqlQuery string
	sqlQuery = `SELECT id, incident_id, author_id, body, created_at FROM incident_comments WHERE incident_id = $1 ORDER BY created_at ASC`

	rows, err := r.pool.Query(ctx, sqlQuery, incidentID)
	if err != nil {
		return nil, err
	}
	defer rows.Close()

	var comments []*model.IncidentComment
	for rows.Next() {
		var comment model.IncidentComment
		if err := rows.Scan(&comment.ID, &comment.IncidentID, &comment.AuthorID, &comment.Body, &comment.CreatedAt); err != nil {
			return nil, err
		}
		comments = append(comments, &comment)
	}

	return comments, nil
}
