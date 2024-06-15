package service

import (
	"WarehouseTgBot/internal/database"
	"context"
	"time"
)

type usersRepository interface {
	CreateUser(ctx context.Context, u database.User) error
	UpdateUser(ctx context.Context, u database.User) error
	FindUserByID(ctx context.Context, userID int64) (database.User, error)
	FindAllActiveUsers(ctx context.Context) ([]database.User, error)
}

type categoryRepository interface {
	CreateCategory(ctx context.Context, gc database.GoodCategory) error
	UpdateCategory(ctx context.Context, gc database.GoodCategory) error
	FindCategoryByID(ctx context.Context, gcID int64) (database.GoodCategory, error)
	FindAllActiveCategories(ctx context.Context) ([]database.GoodCategory, error)
}

type operationRepository interface {
	CreateOperation(ctx context.Context, op database.Operation) error
	FindLastOperationForCategory(ctx context.Context, categoryID int64) (database.Operation, error)
	RemoveLastOperationForCategory(ctx context.Context, categoryID int64) (database.Operation, error)
	ShowCurrentBalancePerCategory(ctx context.Context) ([]database.CurrBalanceRow, error)
	GetOperationsHistory(ctx context.Context, startTime time.Time) (string, error)
}
