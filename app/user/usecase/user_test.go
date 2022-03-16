package usecase_test

import (
	"context"
	"errors"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"librenote/app/model"
	"librenote/app/model/mocks"
	"librenote/app/user/usecase"
	"testing"
	"time"
)

func TestRegistration(t *testing.T) {
	mockUserRepo := new(mocks.UserRepository)

	mockUser := model.User{
		FullName:  "Mr. Test",
		Email:     "mrtest@example.com",
		Hash:      "super_password",
		IsActive:  1,
		UpdatedAt: time.Now().UTC(),
	}

	t.Run("success", func(t *testing.T) {
		tMockUser := mockUser
		tMockUser.ID = 0

		mockUserRepo.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("string")).Return(model.User{}, errors.New("not found")).Once()
		mockUserRepo.On("CreateUser", mock.Anything, mock.AnythingOfType("*model.User")).Return(nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)

		err := u.Registration(context.TODO(), &tMockUser)
		assert.NoError(t, err)
		assert.Equal(t, mockUser.FullName, tMockUser.FullName)
		mockUserRepo.AssertExpectations(t)
	})

	t.Run("existing-user", func(t *testing.T) {
		existingUser := mockUser

		mockUserRepo.On("GetUserByEmail", mock.Anything, mock.AnythingOfType("string")).Return(existingUser, nil).Once()

		u := usecase.NewUserUsecase(mockUserRepo, time.Second*2)
		err := u.Registration(context.TODO(), &existingUser)

		assert.Error(t, err)
		mockUserRepo.AssertExpectations(t)
	})
}
