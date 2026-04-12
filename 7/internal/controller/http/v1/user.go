package v1

import (
	"net/http"

	"practice-7/internal/entity"
	"practice-7/internal/usecase"
	"practice-7/pkg/logger"
	"practice-7/utils"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

type userRoutes struct {
	uc usecase.UserInterface
	l  logger.Interface
}

func newUserRoutes(handler *gin.RouterGroup, uc usecase.UserInterface, l logger.Interface) {
	r := &userRoutes{uc, l}

	h := handler.Group("/users")
	{
		h.POST("/", r.RegisterUser)
		h.POST("/login", r.LoginUser)

		protected := h.Group("/")
		protected.Use(utils.JWTAuthMiddleware())
		{
			protected.GET("/me", r.GetMe)

			admin := protected.Group("/")
			admin.Use(utils.RoleMiddleware("admin"))
			{
				admin.PATCH("/promote/:id", r.PromoteUser)
			}
		}
	}
}

func (r *userRoutes) RegisterUser(c *gin.Context) {
	var dto entity.CreateUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	hashedPassword, err := utils.HashPassword(dto.Password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error hashing password"})
		return
	}

	user := &entity.User{
		Username: dto.Username,
		Email:    dto.Email,
		Password: hashedPassword,
		Role:     "user",
	}

	createdUser, sessionID, err := r.uc.RegisterUser(user)
	if err != nil {
		c.JSON(http.StatusConflict, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message":    "User registered successfully.",
		"session_id": sessionID,
		"user":       toUserResponse(createdUser),
	})
}

func (r *userRoutes) LoginUser(c *gin.Context) {
	var dto entity.LoginUserDTO
	if err := c.ShouldBindJSON(&dto); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := r.uc.LoginUser(&dto)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (r *userRoutes) GetMe(c *gin.Context) {
	userIDStr := c.GetString("userID")
	userID, err := uuid.Parse(userIDStr)
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid user ID in token")
		return
	}

	user, err := r.uc.GetUserByID(userID)
	if err != nil {
		errorResponse(c, http.StatusNotFound, "user not found")
		return
	}

	c.JSON(http.StatusOK, gin.H{"user": toUserResponse(user)})
}

func (r *userRoutes) PromoteUser(c *gin.Context) {
	id, err := uuid.Parse(c.Param("id"))
	if err != nil {
		errorResponse(c, http.StatusBadRequest, "invalid user ID")
		return
	}

	user, err := r.uc.PromoteUser(id)
	if err != nil {
		errorResponse(c, http.StatusNotFound, err.Error())
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "user promoted to admin",
		"user":    toUserResponse(user),
	})
}

func toUserResponse(u *entity.User) entity.UserResponse {
	return entity.UserResponse{
		ID:       u.ID.String(),
		Username: u.Username,
		Email:    u.Email,
		Role:     u.Role,
	}
}
