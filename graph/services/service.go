package services

import (
	"context"
	"mygql/graph/model"

	"github.com/volatiletech/sqlboiler/v4/boil"
)

type Services interface {
	UserService
}

type services struct {
	*userService
}

type UserService interface {
	GetUserByName(ctx context.Context, name string) (*model.User, error)
}

func New(exec boil.ContextExecutor) Services {
	return &services{
		userService: &userService{exec: exec},
	}
}
