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
	RoleDirector   = "director"
	RoleAccountant = "accountant"
)

func NewUserService(usersRepository usersRepository, timeout time.Duration) *UserService {
	return &UserService{usersRepository: usersRepository, timeout: timeout}
}

type UserService struct {
	usersRepository usersRepository
	timeout         time.Duration
}

func (us *UserService) FindAllActiveUsers(ctx context.Context) ([]database.User, error) {
	ctx, cancel := context.WithTimeout(ctx, us.timeout)
	defer cancel()

	return us.usersRepository.FindAllActiveUsers(ctx)
}

func (us *UserService) AddUser(ctx context.Context, user database.User) error {
	ctx, cancel := context.WithTimeout(ctx, us.timeout)
	defer cancel()

	err := us.validateUser(user)
	if err != nil {
		return err
	}

	_, err = us.usersRepository.FindUserByID(ctx, user.ID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return err
	}

	user.Enabled = true

	if errors.Is(err, database.ErrNotFound) {
		err = us.usersRepository.CreateUser(ctx, user)
		if err != nil {
			return err
		}
	} else {
		err := us.usersRepository.UpdateUser(ctx, user)
		if err != nil {
			return err
		}
	}
	return nil
}

func (us *UserService) DisableUser(ctx context.Context, userID int64) error {
	ctx, cancel := context.WithTimeout(ctx, us.timeout)
	defer cancel()

	if userID == 0 {
		return errors.New("invalid user id parameter")
	}

	u, err := us.usersRepository.FindUserByID(ctx, userID)
	if err != nil {
		return err
	}

	u.Enabled = false
	err = us.usersRepository.UpdateUser(ctx, u)
	if err != nil {
		return err
	}

	return nil
}

func (us *UserService) UserHasRole(ctx context.Context, userID int64, roles ...string) (bool, error) {
	ctx, cancel := context.WithTimeout(ctx, us.timeout)
	defer cancel()

	if userID == 0 {
		return false, errors.New("invalid id user parameter")
	}

	u, err := us.usersRepository.FindUserByID(ctx, userID)
	if err != nil && !errors.Is(err, database.ErrNotFound) {
		return false, err
	}

	if !errors.Is(err, database.ErrNotFound) || u.Enabled {
		for _, role := range roles {
			if u.Role == role {
				return true, nil
			}
		}

	}
	return false, nil
}

func (us *UserService) validateUser(user database.User) error {
	if user.ID == 0 || user.Name == "" || (user.Role != RoleDirector && user.Role != RoleAccountant) {
		return errors.New("invalid user parameters")
	}
	return nil
}

func (us *UserService) ConvertToHTML(users []database.User) string {
	buf := bytes.NewBufferString("")
	table := tablewriter.NewWriter(buf)
	table.SetHeader([]string{"ID", "NAME", "ROLE", "CREATED", "UPDATED"})
	for _, u := range users {
		table.Append([]string{strconv.FormatInt(u.ID, 10), u.Name, u.Role, u.CreatedAt.Format(time.Layout), u.UpdatedAt.Format(time.Layout)})
	}
	table.Render()
	return "<pre>\n" + buf.String() + "\n</pre>"
}
