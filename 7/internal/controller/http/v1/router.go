package v1

import (
	"time"

	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
)

func NewRouter(handler *gin.Engine, uc usecase.UserInterface, l logger.Interface) {
	rateLimiter := utils.NewRateLimiter(5, 10*time.Second)
	handler.Use(rateLimiter.Middleware())

	v1 := handler.Group("/v1")
	newUserRoutes(v1, uc, l)
}
