package service

import (
	"WarehouseTgBot/internal/database"
	"bytes"
	"context"
	"errors"
	"github.com/olekukonko/tablewriter"
	"strconv"
	"time"
)

func NewOperationService(operationRepository operationRepository, timeout time.Duration) *OperationService {
	return &OperationService{operationRepository: operationRepository, timeout: timeout}
}

type OperationService struct {
	operationRepository operationRepository
	timeout             time.Duration
}

func (os *OperationService) AddOperation(ctx context.Context, op database.Operation) error {
	ctx, cancel := context.WithTimeout(ctx, os.timeout)
	defer cancel()

	err := os.validateCategory(op)
	if err != nil {
		return err
	}

	err = os.operationRepository.CreateOperation(ctx, op)
	if err != nil {
		return err
	}
	return nil
}

func (os *OperationService) RemoveLastOperation(ctx context.Context, categoryID int64) (database.Operation, error) {
	ctx, cancel := context.WithTimeout(ctx, os.timeout)
	defer cancel()

	lo, err := os.operationRepository.RemoveLastOperationForCategory(ctx, categoryID)
	if err != nil {
		return lo, err
	}
	return lo, nil
}

func (os *OperationService) ShowCurrentBalance(ctx context.Context) ([]database.CurrBalanceRow, error) {
	ctx, cancel := context.WithTimeout(ctx, os.timeout)
	defer cancel()

	cats, err := os.operationRepository.ShowCurrentBalancePerCategory(ctx)
	if err != nil {
		return nil, err
	}
	return cats, nil
}

func (os *OperationService) validateCategory(op database.Operation) error {
	if op.CategoryID == 0 || op.Value == 0 {
		return errors.New("invalid operation parameters")
	}
	return nil
}

func (os *OperationService) ConvertToHTML(cbs []database.CurrBalanceRow) string {
	buf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"ID", "NAME", "UNIT", "COST", "CURRENT_BALANCE", "CURRENT_TOTAL_COST",
		"LAST_OPERATION_TIME", "OPERATION_PERSON", "OPERATION_VAL"})

	for _, c := range cbs {

		var balance = ""
		if c.CurrBalance != nil {
			balance = strconv.FormatFloat(*c.CurrBalance, 'f', 2, 64)
		}
		var totalCost = ""
		if c.TotalCost != nil {
			totalCost = strconv.FormatFloat(*c.TotalCost, 'f', 2, 64)
		}
		var lastOperationTime = ""
		if c.LastOpTime != nil {
			lastOperationTime = c.LastOpTime.Format(time.Layout)
		}

		var lastOpBy = ""
		if c.LastOpBy != nil {
			lastOpBy = strconv.FormatInt(*c.LastOpBy, 10)
		}
		var lastOpValue = ""
		if c.LastOpValue != nil {
			lastOpValue = strconv.FormatFloat(*c.LastOpValue, 'f', 2, 64)
		}

		table.Append([]string{strconv.FormatInt(c.ID, 10), c.Name, c.Unit,
			strconv.FormatFloat(c.Cost, 'f', 2, 64),
			balance,
			totalCost,
			lastOperationTime,
			lastOpBy,
			lastOpValue,
		})
	}
	table.Render()
	return "<pre>\n" + buf.String() + "\n</pre>"
}
