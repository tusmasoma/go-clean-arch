package handler

import (
	"net/http"

	"github.com/labstack/echo/v4"
	"github.com/tusmasoma/go-clean-arch/usecase"
	"github.com/tusmasoma/go-tech-dojo/pkg/log"
)

type UserHandler interface {
	GetUser(c echo.Context) error
	CreateUser(c echo.Context) error
	UpdateUser(c echo.Context) error
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

func (uh *userHandler) GetUser(c echo.Context) error {
	ctx := c.Request().Context()

	user, err := uh.uuc.GetUser(ctx)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	response := GetUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}
	return c.JSON(http.StatusOK, response)
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) CreateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var requestBody CreateUserRequest
	if err := c.Bind(&requestBody); err != nil {
		log.Error("Failed to decode request body", log.Ferror(err))
		return c.NoContent(http.StatusBadRequest)
	}
	if !uh.isValidCreateUserRequest(&requestBody) {
		return c.NoContent(http.StatusBadRequest)
	}

	token, err := uh.uuc.CreateUserAndToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}

	c.Response().Header().Set("Authorization", "Bearer "+token)
	return c.NoContent(http.StatusOK)
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

func (uh *userHandler) UpdateUser(c echo.Context) error {
	ctx := c.Request().Context()

	var requestBody UpdateUserRequest
	if err := c.Bind(&requestBody); err != nil {
		log.Error("Failed to decode request body", log.Ferror(err))
		return c.NoContent(http.StatusBadRequest)
	}
	if !uh.isValidUpdateUserRequest(&requestBody) {
		return c.NoContent(http.StatusBadRequest)
	}

	if err := uh.uuc.UpdateUser(ctx, requestBody.Name); err != nil {
		return c.NoContent(http.StatusInternalServerError)
	}
	return c.NoContent(http.StatusOK)
}

func (uh *userHandler) isValidUpdateUserRequest(requestBody *UpdateUserRequest) bool {
	if requestBody.Name == "" {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}
