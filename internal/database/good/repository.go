package good

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

func (r *Repository) CreateCategory(ctx context.Context, gc database.GoodCategory) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
		INSERT INTO goods (name, description, unit, cost, created_by)
		VALUES ($1, $2, $3, $4, $5)
	`
	if _, err := r.db.Exec(ctx, query, gc.Name, gc.Description, gc.Unit, gc.Cost, gc.CreatedBy); err != nil {
		var writerErr *pgconn.PgError
		if errors.As(err, &writerErr) && writerErr.Code == "23505" {
			return database.ErrConflict
		}
		return fmt.Errorf("postgres Exec: %w", err)
	}
	return nil
}

func (r *Repository) UpdateCategory(ctx context.Context, gc database.GoodCategory) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `UPDATE goods SET name=$2, description=$3, enabled=$4, updated_at=$5, unit=$6, cost=$7  WHERE id=$1`
	if _, err := r.db.Exec(ctx, query, gc.ID, gc.Name, gc.Description, gc.Enabled, time.Now(), gc.Unit, gc.Cost); err != nil {
		return fmt.Errorf("postgres Exec: %w", err)
	}
	return nil
}

func (r *Repository) FindCategoryByID(ctx context.Context, gcID int64) (database.GoodCategory, error) {
	var gc database.GoodCategory

	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	if err := r.db.QueryRow(ctx,
		`SELECT * FROM goods WHERE id=$1`, gcID).
		Scan(&gc.ID, &gc.Name, &gc.Description, &gc.Unit, &gc.Cost, &gc.CreatedBy, &gc.CreatedAt, &gc.UpdatedAt, &gc.Enabled); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return gc, database.ErrNotFound
		}
		return gc, fmt.Errorf("postgres QueryRow Decode: %w", err)
	}
	return gc, nil
}

func (r *Repository) FindAllActiveCategories(ctx context.Context) ([]database.GoodCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var users []database.GoodCategory

	query := `SELECT * FROM goods WHERE enabled = TRUE`

	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var gc database.GoodCategory
		err := rows.Scan(&gc.ID, &gc.Name, &gc.Description, &gc.Unit, &gc.Cost, &gc.CreatedBy, &gc.CreatedAt, &gc.UpdatedAt, &gc.Enabled)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		users = append(users, gc)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return users, nil
}
