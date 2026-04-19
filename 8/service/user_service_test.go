package service

import (
	"errors"
	"testing"

	"practice-8/repository"

	"github.com/stretchr/testify/assert"
	"go.uber.org/mock/gomock"
)

func TestGetUserByID(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().GetUserByID(1).Return(user, nil)

	result, err := userService.GetUserByID(1)
	assert.NoError(t, err)
	assert.Equal(t, user, result)
}

func TestCreateUser(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	userService := NewUserService(mockRepo)

	user := &repository.User{ID: 1, Name: "Bakytzhan Agai"}
	mockRepo.EXPECT().CreateUser(user).Return(nil)

	err := userService.CreateUser(user)
	assert.NoError(t, err)
}

func TestRegisterUser_UserAlreadyExists(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	existing := &repository.User{ID: 2, Name: "Alice", Email: "alice@example.com"}
	mockRepo.EXPECT().GetByEmail("alice@example.com").Return(existing, nil)

	err := svc.RegisterUser(&repository.User{Name: "Bob"}, "alice@example.com")
	assert.ErrorContains(t, err, "user with this email already exists")
}

func TestRegisterUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	newUser := &repository.User{ID: 3, Name: "Charlie", Email: "charlie@example.com"}
	mockRepo.EXPECT().GetByEmail("charlie@example.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(nil)

	err := svc.RegisterUser(newUser, "charlie@example.com")
	assert.NoError(t, err)
}

func TestRegisterUser_CreateUserError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	newUser := &repository.User{ID: 4, Name: "Dave", Email: "dave@example.com"}
	dbErr := errors.New("db write failed")
	mockRepo.EXPECT().GetByEmail("dave@example.com").Return(nil, nil)
	mockRepo.EXPECT().CreateUser(newUser).Return(dbErr)

	err := svc.RegisterUser(newUser, "dave@example.com")
	assert.ErrorIs(t, err, dbErr)
}

func TestUpdateUserName_EmptyName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.UpdateUserName(2, "")
	assert.ErrorContains(t, err, "name cannot be empty")
}

func TestUpdateUserName_UserNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	repoErr := errors.New("user not found")
	mockRepo.EXPECT().GetUserByID(99).Return(nil, repoErr)

	err := svc.UpdateUserName(99, "NewName")
	assert.ErrorIs(t, err, repoErr)
}

func TestUpdateUserName_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	mockRepo.EXPECT().GetUserByID(2).Return(&repository.User{ID: 2, Name: "OldName"}, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *repository.User) error {
		assert.Equal(t, "NewName", u.Name, "name must be updated before calling UpdateUser")
		return nil
	})

	err := svc.UpdateUserName(2, "NewName")
	assert.NoError(t, err)
}

func TestUpdateUserName_UpdateFails(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	dbErr := errors.New("db write error")
	mockRepo.EXPECT().GetUserByID(2).Return(&repository.User{ID: 2, Name: "OldName"}, nil)
	mockRepo.EXPECT().UpdateUser(gomock.Any()).DoAndReturn(func(u *repository.User) error {
		assert.Equal(t, "NewName", u.Name, "name must be updated before calling UpdateUser")
		return dbErr
	})

	err := svc.UpdateUserName(2, "NewName")
	assert.ErrorIs(t, err, dbErr)
}

func TestDeleteUser_AdminNotAllowed(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	err := svc.DeleteUser(1)
	assert.ErrorContains(t, err, "it is not allowed to delete admin user")
}

func TestDeleteUser_Success(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	deleted := false
	mockRepo.EXPECT().DeleteUser(5).DoAndReturn(func(id int) error {
		deleted = true
		return nil
	})

	err := svc.DeleteUser(5)
	assert.NoError(t, err)
	assert.True(t, deleted, "DeleteUser should have been called on the repository")
}

func TestDeleteUser_RepositoryError(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	mockRepo := repository.NewMockUserRepository(ctrl)
	svc := NewUserService(mockRepo)

	repoErr := errors.New("db delete failed")
	mockRepo.EXPECT().DeleteUser(5).Return(repoErr)

	err := svc.DeleteUser(5)
	assert.ErrorIs(t, err, repoErr)
}
