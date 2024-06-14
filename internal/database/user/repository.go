package user

import (
	"WarehouseTgBot/internal/database"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgx/v4"
	"time"

	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4/pgxpool"
)

func NewRepository(userDB *pgxpool.Pool, timeout time.Duration) *Repository {
	return &Repository{db: userDB, timeout: timeout}
}

type Repository struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func (r *Repository) CreateUser(ctx context.Context, u database.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
		INSERT INTO users (id, name, role)
		VALUES ($1, $2, $3)
	`
	if _, err := r.db.Exec(ctx, query, u.ID, u.Name, u.Role); err != nil {
		var writerErr *pgconn.PgError
		if errors.As(err, &writerErr) && writerErr.Code == "23505" {
			return database.ErrConflict
		}
		return fmt.Errorf("postgres Exec: %w", err)
	}

	return nil
}

func (r *Repository) UpdateUser(ctx context.Context, u database.User) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `UPDATE users SET name=$2, role=$3, enabled=$4, updated_at=$5  WHERE id=$1`
	if _, err := r.db.Exec(ctx, query, u.ID, u.Name, u.Role, u.Enabled, time.Now()); err != nil {
		return fmt.Errorf("postgres Exec: %w", err)
	}
	return nil
}

func (r *Repository) FindUserByID(ctx context.Context, userID int64) (database.User, error) {
	var u database.User

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err := r.db.QueryRow(ctx,
		`SELECT id, name, role, created_at, updated_at, enabled FROM users WHERE id=$1`, userID).
		Scan(&u.ID, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt, &u.Enabled); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return u, database.ErrNotFound
		}
		return u, fmt.Errorf("postgres QueryRow Decode: %w", err)
	}
	return u, nil
}

func (r *Repository) FindAllActiveUsers(ctx context.Context) ([]database.User, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var users []database.User

	query := `SELECT id, name, role, created_at, updated_at, enabled FROM users WHERE enabled = TRUE`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var u database.User
		err := rows.Scan(&u.ID, &u.Name, &u.Role, &u.CreatedAt, &u.UpdatedAt, &u.Enabled)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, u)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return users, nil
}
