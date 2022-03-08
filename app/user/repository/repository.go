package repository

import (
	"context"
	"librenote/app/model"
)

type UserRepository interface {
	CreateUser(tx context.Context, user *model.User) error
	GetUser(tx context.Context, id int32) (model.User, error)
	GetUserByEmail(tx context.Context, email string) (model.User, error)
	UpdateUser(tx context.Context, user *model.User) error
}
