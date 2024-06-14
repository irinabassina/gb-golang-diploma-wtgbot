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

const (
	UnitPcs = "pcs"
	UnitKG  = "kg"
)

func NewCategoryService(categoryRepository categoryRepository, timeout time.Duration) *CategoryService {
	return &CategoryService{categoryRepository: categoryRepository, timeout: timeout}
}

type CategoryService struct {
	categoryRepository categoryRepository
	timeout            time.Duration
}

func (cs *CategoryService) FindAllActiveCategories(ctx context.Context) ([]database.GoodCategory, error) {
	ctx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()

	return cs.categoryRepository.FindAllActiveCategories(ctx)
}

func (cs *CategoryService) AddCategory(ctx context.Context, gc database.GoodCategory) error {
	ctx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()

	err := cs.validateCategory(gc)
	if err != nil {
		return err
	}

	err = cs.categoryRepository.CreateCategory(ctx, gc)
	if err != nil {
		return err
	}
	return nil
}

func (cs *CategoryService) UpdateCategory(ctx context.Context, gc database.GoodCategory) error {
	ctx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()

	err := cs.validateCategory(gc)
	if err != nil {
		return err
	} else if gc.ID == 0 {
		return errors.New("category ID cannot be empty")
	}

	_, err = cs.categoryRepository.FindCategoryByID(ctx, gc.ID)
	if err != nil {
		return err
	}

	err = cs.categoryRepository.UpdateCategory(ctx, gc)
	if err != nil {
		return err

	}
	return nil
}

func (cs *CategoryService) DisableCategory(ctx context.Context, gcID int64) error {
	ctx, cancel := context.WithTimeout(ctx, cs.timeout)
	defer cancel()

	if gcID == 0 {
		return errors.New("invalid category id parameter")
	}

	gc, err := cs.categoryRepository.FindCategoryByID(ctx, gcID)
	if err != nil {
		return err
	}

	gc.Enabled = false
	err = cs.categoryRepository.UpdateCategory(ctx, gc)
	if err != nil {
		return err
	}

	return nil
}

func (cs *CategoryService) validateCategory(gc database.GoodCategory) error {
	if gc.Cost <= 0 || gc.Name == "" || gc.Description == "" || (gc.Unit != UnitPcs && gc.Unit != UnitKG) {
		return errors.New("invalid good category parameters")
	}
	return nil
}

func (cs *CategoryService) ConvertToHTML(cats []database.GoodCategory) string {
	buf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"ID", "NAME", "DESC", "UNIT", "COST", "CREATED_BY", "CREATED", "UPDATED"})
	for _, c := range cats {
		table.Append([]string{strconv.FormatInt(c.ID, 10), c.Name, c.Description, c.Unit,
			strconv.FormatFloat(c.Cost, 'f', 2, 64), strconv.FormatInt(c.CreatedBy, 10),
			c.CreatedAt.Format(time.Layout), c.UpdatedAt.Format(time.Layout)})
	}
	table.Render()
	return "<pre>\n" + buf.String() + "\n</pre>"
}
