package usecase

import (
	"context"
	"database/sql"
	"errors"
	"librenote/app/model"
	"librenote/app/response"
	"librenote/infrastructure/config"
	"librenote/infrastructure/middlewares"
	"net/http"
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

type userUsecase struct {
	repo           model.UserRepository
	contextTimeout time.Duration
}

func NewUserUsecase(repo model.UserRepository, timeout time.Duration) model.UserUsecase {
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
		return response.ErrConflict
	}

	// generate password salted hash
	hash, err := bcrypt.GenerateFromPassword([]byte(m.Hash), bcrypt.DefaultCost)
	if err != nil {
		return
	}

	m.Hash = string(hash)

	// store
	err = u.repo.CreateUser(ctx, m)

	return
}

func (u *userUsecase) Login(c context.Context, email, password string) (token string, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.repo.GetUserByEmail(ctx, email)
	if err != nil {
		return "", response.WrapError(errors.New("email/password is incorrect"), http.StatusUnauthorized)
	}

	// check password
	err = bcrypt.CompareHashAndPassword([]byte(user.Hash), []byte(password))
	if err != nil {
		return "", response.WrapError(errors.New("email/password is incorrect"), http.StatusUnauthorized)
	}

	// check user state
	if user.IsActive == 0 || user.IsTrashed == 1 {
		return "", response.WrapError(errors.New("user not exist or inactive"), http.StatusUnauthorized)
	}

	token, err = createToken(user.ID)
	if err != nil {
		return "", err
	}

	return token, nil
}

func createToken(userID int32) (string, error) {
	jwtCfg := config.Get().Jwt

	claims := &middlewares.JwtCustomClaims{
		UserID: userID,
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(jwtCfg.ExpireTime).Unix(),
		},
	}

	unsignedToken := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	token, err := unsignedToken.SignedString([]byte(jwtCfg.SecretKey))

	if err != nil {
		return "", err
	}

	return token, err
}

func (u *userUsecase) GetUserDetails(c context.Context, id int32) (details *model.UserDetails, err error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, response.ErrNotFound
		}

		return nil, err
	}

	// check user state
	if user.IsActive == 0 || user.IsTrashed == 1 {
		return nil, response.WrapError(errors.New("user not exist or inactive"), http.StatusUnauthorized)
	}

	details = &model.UserDetails{
		FullName:        user.FullName,
		Email:           user.Email,
		ListViewEnabled: user.ListViewEnabled,
		DarkModeEnabled: user.DarkModeEnabled,
	}

	return details, nil
}

func (u *userUsecase) GetUser(c context.Context, id int32) (*model.User, error) {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	user, err := u.repo.GetUser(ctx, id)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return &user, response.ErrNotFound
		}

		return &user, err
	}

	// check user state
	if user.IsActive == 0 || user.IsTrashed == 1 {
		return &user, response.WrapError(errors.New("user not exist or inactive"), http.StatusUnauthorized)
	}

	return &user, nil
}

func (u *userUsecase) Update(c context.Context, m *model.User, p model.Password) error {
	ctx, cancel := context.WithTimeout(c, u.contextTimeout)
	defer cancel()

	if p.IsChanged {
		// check old password is correct
		err := bcrypt.CompareHashAndPassword([]byte(m.Hash), []byte(p.OldPassword))
		if err != nil {
			return response.WrapError(errors.New("old password doesn't match"), http.StatusBadRequest)
		}

		// generate password salted hash
		hash, err := bcrypt.GenerateFromPassword([]byte(p.NewPassword), bcrypt.DefaultCost)
		if err != nil {
			return err
		}

		m.Hash = string(hash)
	}

	// update
	return u.repo.UpdateUser(ctx, m)
}
