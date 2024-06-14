package database

import (
	"time"
)

type User struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Role      string    `db:"role"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Enabled   bool      `db:"enabled"`
}

type GoodCategory struct {
	ID          int64     `db:"id"`
	Name        string    `db:"name"`
	Description string    `db:"role"`
	Unit        string    `db:"unit"`
	Cost        float64   `db:"cost"`
	CreatedBy   int64     `db:"created_by"`
	CreatedAt   time.Time `db:"created_at"`
	UpdatedAt   time.Time `db:"updated_at"`
	Enabled     bool      `db:"enabled"`
}

type Operation struct {
	ID          int64     `db:"id"`
	CategoryID  int64     `db:"category_id"`
	Value       float64   `db:"value"`
	CurrBalance float64   `db:"current_balance"`
	CreatedBy   int64     `db:"created_by"`
	CreatedAt   time.Time `db:"created_at"`
}

type CurrBalanceRow struct {
	ID          int64      `db:"id"`
	Name        string     `db:"name"`
	Unit        string     `db:"unit"`
	Cost        float64    `db:"cost"`
	CurrBalance *float64   `db:"current_balance"`
	TotalCost   *float64   `db:"total_cost"`
	LastOpTime  *time.Time `db:"last_op_time"`
	LastOpBy    *int64     `db:"last_op_by"`
	LastOpValue *float64   `db:"last_op_val"`
}
