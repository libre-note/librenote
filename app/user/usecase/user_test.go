package usecase_test

import (
	"context"
	"errors"
	"librenote/app/model"
	"librenote/app/model/mocks"
	"librenote/app/user/usecase"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"golang.org/x/crypto/bcrypt"
)

func TestRegistration(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)
	nowTime := time.Now().UTC().Format("2006-01-02 15:04:05")

	mockUser := model.User{
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      "super_password",
		IsActive:  1,
		CreatedAt: nowTime,
		UpdatedAt: nowTime,
	}

	t.Run("success", func(t *testing.T) {
		tMockUser := mockUser
		tMockUser.ID = 0

		mockUserRepo.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(model.User{}, errors.New("not found")).Once()
		mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)

		err := u.Registration(context.TODO(), &tMockUser)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.FullName, tMockUser.FullName)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("existing-user", func(t *testing.T) {
		existingUser := mockUser

		mockUserRepo.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(existingUser, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		err := u.Registration(context.TODO(), &existingUser)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestLoginSuccessAndWrongPassword(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	hash, _ := bcrypt.GenerateFromPassword([]byte("super_password"), bcrypt.MinCost)
	mockUser := model.User{
		ID:        1,
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      string(hash),
		IsActive:  1,
		IsTrashed: 0,
	}

	t.Run("success", func(t *testing.T) {
		existingUser := mockUser

		mockUserRepo.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(existingUser, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		token, err := u.Login(context.TODO(), "mrtest@example.com", "super_password")

		assert.NoError(t, err)
		assert.NotEqual(t, "", token)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("wrong-password", func(t *testing.T) {
		existingUser := mockUser

		mockUserRepo.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(existingUser, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		_, err := u.Login(context.TODO(), "mrtest@example.com", "super")

		assert.Error(t, err)
		assert.EqualError(t, err, "email/password is incorrect")
	})
}

func TestLoginWrongEmailAndInactive(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	hash, _ := bcrypt.GenerateFromPassword([]byte("super_password"), bcrypt.MinCost)
	mockUser := model.User{
		ID:        1,
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      string(hash),
		IsActive:  1,
		IsTrashed: 0,
	}

	t.Run("wrong-email", func(t *testing.T) {
		mockUserRepo.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(model.User{}, errors.New("not found")).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		_, err := u.Login(context.TODO(), "test@example.com", "super_password")

		assert.Error(t, err)
		assert.EqualError(t, err, "email/password is incorrect")
	})

	t.Run("inactive", func(t *testing.T) {
		existingUser := mockUser
		existingUser.IsActive = 0

		mockUserRepo.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("string")).
			Return(existingUser, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		_, err := u.Login(context.TODO(), "mrtest@example.com", "super_password")

		assert.Error(t, err)
		assert.EqualError(t, err, "user not exist or inactive")
	})
}

func TestGetUserDetails(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	hash, _ := bcrypt.GenerateFromPassword([]byte("super_password"), bcrypt.MinCost)
	mockUser := model.User{
		ID:              1,
		FullName:        "Mr. Test",
		Email:           "mrtest@example.com",
		Hash:            string(hash),
		IsActive:        1,
		IsTrashed:       0,
		ListViewEnabled: 1,
		DarkModeEnabled: 1,
	}

	t.Run("success", func(t *testing.T) {
		existingUser := mockUser

		mockUserRepo.On("GetUser", mock.Anything, mock.AnythingOfType("int32")).
			Return(existingUser, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		details, err := u.GetUserDetails(context.TODO(), 1)

		assert.NoError(t, err)
		assert.Equal(t, existingUser.Email, details.Email)
		assert.Equal(t, existingUser.ListViewEnabled, details.ListViewEnabled)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("not-found", func(t *testing.T) {
		mockUserRepo.On("GetUser", mock.Anything, mock.AnythingOfType("int32")).
			Return(model.User{}, errors.New("no row found")).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		_, err := u.GetUserDetails(context.TODO(), 2)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}

func TestUpdate(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	hash, _ := bcrypt.GenerateFromPassword([]byte("super_password"), bcrypt.MinCost)
	mockUser := model.User{
		ID:        1,
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      string(hash),
		IsActive:  1,
		IsTrashed: 0,
	}

	t.Run("change-password", func(t *testing.T) {
		existingUser := mockUser
		pass := model.Password{
			OldPassword: "super_password",
			NewPassword: "super_new_pass",
			IsChanged:   true,
		}

		mockUserRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		err := u.Update(context.TODO(), &existingUser, pass)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("wrong-password", func(t *testing.T) {
		existingUser := mockUser
		pass := model.Password{
			OldPassword: "superpassword",
			NewPassword: "super_new_pass",
			IsChanged:   true,
		}

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		err := u.Update(context.TODO(), &existingUser, pass)

		assert.Error(t, err)
		assert.EqualError(t, err, "old password doesn't match")
		mockUserRepo.AssertExpectations(t)
	})
}

func TestDelete(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	hash, _ := bcrypt.GenerateFromPassword([]byte("super_password"), bcrypt.MinCost)
	mockUser := model.User{
		ID:        1,
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      string(hash),
		IsActive:  1,
		IsTrashed: 0,
	}

	t.Run("delete", func(t *testing.T) {
		existingUser := mockUser
		existingUser.IsTrashed = 1
		pass := model.Password{IsChanged: false}

		mockUserRepo.On("UpdateUser", mock.Anything, mock.AnythingOfType("*model.User")).
			Return(nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		err := u.Update(context.TODO(), &existingUser, pass)

		assert.NoError(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}
