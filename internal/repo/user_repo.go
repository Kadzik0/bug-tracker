package repo

import (
	"context"

	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/kadzik0/bug-tracker/internal/model"
)

type UserRepo struct {
	pool *pgxpool.Pool
}

func NewUserRepo(pool *pgxpool.Pool) *UserRepo {
	return &UserRepo{pool: pool}
}

func (r *UserRepo) Create(ctx context.Context, user *model.User) error {
	var sqlQuery string
	sqlQuery = `INSERT INTO users (id, email, name, created_at) VALUES ($1, $2, $3, $4)`

	_, err := r.pool.Exec(ctx, sqlQuery, user.ID, user.Email, user.Name, user.CreatedAt)
	return err
}

func (r *UserRepo) List(ctx context.Context) ([]*model.User, error) {
	var sqlQuery string
	sqlQuery = `SELECT id, email, name, created_at FROM users`

	result, err := r.pool.Query(ctx, sqlQuery)
	if err != nil {
		return nil, err
	}
	defer result.Close()

	var users []*model.User
	for result.Next() {
		var user model.User
		if err := result.Scan(&user.ID, &user.Email, &user.Name, &user.CreatedAt); err != nil {
			return nil, err
		}
		users = append(users, &user)
	}
	if err := result.Err(); err != nil {
		return nil, err
	}
	return users, nil
}
