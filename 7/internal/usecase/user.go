package usecase

import (
	"fmt"

	"practice-7/internal/entity"
	"practice-7/internal/usecase/repo"
	"practice-7/utils"

	"github.com/google/uuid"
)

type UserUseCase struct {
	repo *repo.UserRepo
}

func NewUserUseCase(r *repo.UserRepo) *UserUseCase {
	return &UserUseCase{repo: r}
}

func (u *UserUseCase) RegisterUser(user *entity.User) (*entity.User, string, error) {
	created, err := u.repo.RegisterUser(user)
	if err != nil {
		return nil, "", fmt.Errorf("register user: %w", err)
	}
	sessionID := uuid.New().String()
	return created, sessionID, nil
}

func (u *UserUseCase) LoginUser(input *entity.LoginUserDTO) (string, error) {
	userFromRepo, err := u.repo.GetByUsername(input.Username)
	if err != nil {
		return "", fmt.Errorf("user not found: %w", err)
	}

	if !utils.CheckPassword(userFromRepo.Password, input.Password) {
		return "", fmt.Errorf("invalid credentials")
	}

	token, err := utils.GenerateJWT(userFromRepo.ID, userFromRepo.Role)
	if err != nil {
		return "", fmt.Errorf("generate JWT: %w", err)
	}

	return token, nil
}

func (u *UserUseCase) GetUserByID(id uuid.UUID) (*entity.User, error) {
	user, err := u.repo.GetByID(id)
	if err != nil {
		return nil, fmt.Errorf("get user: %w", err)
	}
	return user, nil
}

func (u *UserUseCase) PromoteUser(id uuid.UUID) (*entity.User, error) {
	user, err := u.repo.UpdateRole(id, "admin")
	if err != nil {
		return nil, fmt.Errorf("promote user: %w", err)
	}
	return user, nil
}
