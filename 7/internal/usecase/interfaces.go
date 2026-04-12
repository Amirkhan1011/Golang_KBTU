package usecase

import (
	"practice-7/internal/entity"

	"github.com/google/uuid"
)

type UserInterface interface {
	RegisterUser(user *entity.User) (*entity.User, string, error)
	LoginUser(input *entity.LoginUserDTO) (string, error)
	GetUserByID(id uuid.UUID) (*entity.User, error)
	PromoteUser(id uuid.UUID) (*entity.User, error)
}
