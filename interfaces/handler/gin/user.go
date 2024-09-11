package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/usecase"
)

type UserHandler interface {
	GetUser(c *gin.Context)
	CreateUser(c *gin.Context)
	UpdateUser(c *gin.Context)
}

type userHandler struct {
	uuc usecase.UserUseCase
}

func NewUserHandler(uuc usecase.UserUseCase) UserHandler {
	return &userHandler{
		uuc: uuc,
	}
}

type GetUserResponse struct {
	ID    string `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
}

func (uh *userHandler) GetUser(c *gin.Context) {
	ctx := c.Request.Context()

	user, err := uh.uuc.GetUser(ctx)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	response := GetUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	c.JSON(http.StatusOK, response)
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) CreateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var requestBody CreateUserRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Error("Failed to decode request body", log.Ferror(err))
		c.Status(http.StatusBadRequest)
		return
	}
	if !uh.isValidCreateUserRequest(&requestBody) {
		c.Status(http.StatusBadRequest)
		return
	}

	token, err := uh.uuc.CreateUserAndToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}

	c.Header("Authorization", "Bearer "+token)
	c.Status(http.StatusOK)
}

func (uh *userHandler) isValidCreateUserRequest(requestBody *CreateUserRequest) bool {
	if requestBody.Email == "" || requestBody.Password == "" {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}

type UpdateUserRequest struct {
	Name string `json:"name"`
	// Email 	string `json:"email"`
	// Password 	string `json:"password"`
}

func (uh *userHandler) UpdateUser(c *gin.Context) {
	ctx := c.Request.Context()

	var requestBody UpdateUserRequest
	if err := c.ShouldBindJSON(&requestBody); err != nil {
		log.Error("Failed to decode request body", log.Ferror(err))
		c.Status(http.StatusBadRequest)
		return
	}
	if !uh.isValidUpdateUserRequest(&requestBody) {
		c.Status(http.StatusBadRequest)
		return
	}

	if err := uh.uuc.UpdateUser(ctx, requestBody.Name); err != nil {
		c.Status(http.StatusInternalServerError)
		return
	}
	c.Status(http.StatusOK)
}

func (uh *userHandler) isValidUpdateUserRequest(requestBody *UpdateUserRequest) bool {
	if requestBody.Name == "" {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}
