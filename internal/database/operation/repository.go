package operation

import (
	"WarehouseTgBot/internal/database"
	"context"
	"errors"
	"fmt"
	"github.com/jackc/pgconn"
	"github.com/jackc/pgx/v4"
	"github.com/jackc/pgx/v4/pgxpool"
	"time"
)

func NewRepository(userDB *pgxpool.Pool, timeout time.Duration) *Repository {
	return &Repository{db: userDB, timeout: timeout}
}

type Repository struct {
	db      *pgxpool.Pool
	timeout time.Duration
}

func (r *Repository) CreateOperation(ctx context.Context, op database.Operation) error {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query := `
		INSERT INTO operations (category_id, value, current_balance, created_by)
		VALUES ($1, $2, $3, $4)
	`

	lo, err := r.FindLastOperationForCategory(ctx, op.CategoryID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return fmt.Errorf("error finding last operation: %w", err)
	}
	var balance float64 = 0
	if !errors.Is(err, database.ErrNotFound) {
		balance = lo.CurrBalance
	}
	balance = balance + op.Value
	if balance < 0 {
		return fmt.Errorf("resulted balance cannot be negative")
	}

	if _, err := r.db.Exec(ctx, query, op.CategoryID, op.Value, balance, op.CreatedBy); err != nil {
		var writerErr *pgconn.PgError
		if errors.As(err, &writerErr) && writerErr.Code == "23505" {
			return database.ErrConflict
		}
		return fmt.Errorf("postgres Exec: %w", err)
	}
	return nil
}

func (r *Repository) FindLastOperationForCategory(ctx context.Context, categoryID int64) (database.Operation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var lo database.Operation
	if err := r.db.QueryRow(ctx,
		`SELECT * FROM operations WHERE category_id=$1 ORDER BY created_at DESC LIMIT 1`, categoryID).
		Scan(&lo.ID, &lo.CategoryID, &lo.Value, &lo.CurrBalance, &lo.CreatedBy, &lo.CreatedAt); err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return lo, database.ErrNotFound
		}
		return lo, fmt.Errorf("postgres QueryRow Decode: %w", err)
	}
	return lo, nil
}

func (r *Repository) RemoveLastOperationForCategory(ctx context.Context, categoryID int64) (database.Operation, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var lo database.Operation

	lo, err := r.FindLastOperationForCategory(ctx, categoryID)
	if err != nil {
		return lo, fmt.Errorf("error deleting last operation: %w", err)
	}

	if _, err := r.db.Exec(ctx, `DELETE FROM operations WHERE id=$1`, lo.ID); err != nil {
		return lo, fmt.Errorf("postgres Exec: %w", err)
	}
	return lo, nil
}

func (r *Repository) ShowCurrentBalancePerCategory(ctx context.Context) ([]database.CurrBalanceRow, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	var balanceRows []database.CurrBalanceRow
	query := `
		SELECT goods.id,
		   goods.name,
		   goods.unit,
		   goods.cost,
		   lastOp.current_balance,
		   (lastOp.current_balance * goods.cost) as total_cost,
		   lastOp.created_at                     as last_op_time,
		   lastOp.created_by                     as last_op_by,
		   lastOp.value                          as last_op_val
		FROM goods
		LEFT JOIN
			(SELECT * FROM 
				(SELECT *, row_number() over (partition by category_id order by created_at DESC ) as rw
					FROM operations) as "tmp"
			where rw = 1) as lastOp
	 	ON goods.id = lastOp.category_id;
	`
	rows, err := r.db.Query(ctx, query)
	if err != nil {
		return nil, fmt.Errorf("postgres Query: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		var cbr database.CurrBalanceRow
		err := rows.Scan(&cbr.ID, &cbr.Name, &cbr.Unit, &cbr.Cost, &cbr.CurrBalance, &cbr.TotalCost, &cbr.LastOpTime, &cbr.LastOpBy, &cbr.LastOpValue)
		if err != nil {
			return nil, fmt.Errorf("failed to scan row: %w", err)
		}
		balanceRows = append(balanceRows, cbr)
	}

	if err := rows.Err(); err != nil {
		return nil, fmt.Errorf("error during rows iteration: %w", err)
	}

	return balanceRows, nil
}

func (r *Repository) GetOperationsHistory(ctx context.Context, startTime time.Time) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, r.timeout)
	defer cancel()

	query :=
		`SELECT (goods.id,
			   goods.name,
			   goods.unit,
			   goods.cost,
			   lastOp.current_balance,
			   lastOp.created_at,
			   lastOp.created_by,
			   lastOp.value) :: TEXT
		FROM goods
				 JOIN
			 (SELECT *
			  FROM operations
			  WHERE created_at > $1) as lastOp
			 ON goods.id = lastOp.category_id
		ORDER BY lastOp.created_at DESC`

	rows, err := r.db.Query(ctx, query, startTime)
	if err != nil {
		return "", fmt.Errorf("postgres Query: %w", err)
	}
	defer rows.Close()

	resultCSV := "category_id,category_name,category_unit,category_cost,current_balance,created_at,created_by,value"
	for rows.Next() {
		rowStr := ""
		err := rows.Scan(&rowStr)
		if err != nil {
			return "", fmt.Errorf("failed to scan row: %w", err)
		}
		rowStr = rowStr[:len(rowStr)-1]
		rowStr = rowStr[1:]
		resultCSV = resultCSV + "\n" + rowStr
	}

	if err := rows.Err(); err != nil {
		return "", fmt.Errorf("error during rows iteration: %w", err)
	}

	return resultCSV, nil
}
