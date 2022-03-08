package usecase

import (
	"librenote/app/model"
	"librenote/app/user/repository"
)

type UserUsecase interface {
	Registration(user *model.User) error
}

type userUsecase struct {
	repo repository.UserRepository
}

func NewUserUsecase(repo repository.UserRepository) UserUsecase {
	return &userUsecase{
		repo: repo,
	}
}

func (u *userUsecase) Registration(user *model.User) error {
	// check user already exists

	// generate password hash & salt

	// store

	return nil
}
