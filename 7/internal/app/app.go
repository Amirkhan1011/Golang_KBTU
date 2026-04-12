package app

import (
	"fmt"

	"practice-7/config"
	v1 "practice-7/internal/controller/http/v1"
	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/internal/usecase/repo"
	"practice-7/pkg/logger"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func Run(cfg *config.Config) error {
	utils.SetJWTSecret(cfg.JWTSecret)

	l := logger.New()

	userRepo := repo.NewUserRepo()

	adminHash, err := utils.HashPassword("admin123")
	if err != nil {
		return fmt.Errorf("seed admin: %w", err)
	}
	userRepo.SeedAdmin(&entity.User{
		ID:       uuid.New(),
		Username: "admin",
		Email:    "admin@example.com",
		Password: adminHash,
		Role:     "admin",
	})
	l.Info("Seeded admin  → username: admin  | password: admin123")

	userUC := usecase.NewUserUseCase(userRepo)

	gin.SetMode(gin.ReleaseMode)
	router := gin.New()
	router.Use(gin.Logger(), gin.Recovery())

	v1.NewRouter(router, userUC, l)

	l.Info("Server listening on :%s", cfg.Port)
	return router.Run(":" + cfg.Port)
}
