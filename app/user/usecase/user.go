package usecase

import (
	"context"
	"errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/bcrypt"
	"librenote/app/model"
	"librenote/app/user/repository"
	"time"
)

type UserUsecase interface {
	Registration(c context.Context, m *model.User) (err error)
	Login(c context.Context, email, password string) (user model.User, err error)
	GetByID(c context.Context, id int32) (user model.User, err error)
	Update(c context.Context, m *model.User, p Password) error
}

type userUsecase struct {
	repo           repository.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(repo repository.UserRepository, timeout time.Duration) UserUsecase {
	return &userUsecase{
		repo:           repo,
		contextTimeout: timeout,
	}
}

func (u *userUsecase) Registration(c context.Context, m *model.User) (err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	// check user already exists
	existedUser, _ := u.repo.GetUserByEmail(ctx, m.Email)
	if existedUser != (model.User{}) {
		return errors.New("user already exist")
	}

	// generate password salted hash
	hash, err := bcrypt.GenerateFromPassword([]byte(m.Hash), bcrypt.MinCost)
	if err != nil {
		return
	}
	m.Hash = string(hash)

	// store
	err = u.repo.CreateUser(ctx, m)
	return
}

func (u *userUsecase) Login(c context.Context, email, password string) (user model.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err = u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		logrus.Error(err)
		return user, errors.New("user/password is incorrect")

	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		logrus.Error(err)
		return user, errors.New("user/password is incorrect")
	}

	return user, nil
}

func (u *userUsecase) GetByID(c context.Context, id int32) (user model.User, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err = u.repo.GetUser(ctx, id)
	if err != nil {
		return
	}
	user.Hash = ""

	return
}

type Password struct {
	OldPassword string
	NewPassword string
	IsChanged   bool
}

func (u *userUsecase) Update(c context.Context, m *model.User, p Password) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if p.IsChanged {
		// check old password is correct
		err := bcrypt.CompareHashAndPassword([]byte(m.Hash), []byte(p.OldPassword))
		if err != nil {
			return errors.New("old password doesn't match")
		}
		// generate password salted hash
		hash, err := bcrypt.GenerateFromPassword([]byte(p.NewPassword), bcrypt.MinCost)
		if err != nil {
			return err
		}
		m.Hash = string(hash)
	}

	// update
	err := u.repo.UpdateUser(ctx, m)
	return err
}
