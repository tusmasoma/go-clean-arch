package handler

import (
	"encoding/json"
	"io"
	"net/http"

	"github.com/tusmasoma/go-tech-dojo/pkg/log"

	"github.com/tusmasoma/go-clean-arch/usecase"
)

type UserHandler interface {
	GetUser(w http.ResponseWriter, r *http.Request)
	CreateUser(w http.ResponseWriter, r *http.Request)
	UpdateUser(w http.ResponseWriter, r *http.Request)
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

func (uh *userHandler) GetUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	user, err := uh.uuc.GetUser(ctx)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	if err = json.NewEncoder(w).Encode(GetUserResponse{
		ID:    user.ID,
		Name:  user.Name,
		Email: user.Email,
	}); err != nil {
		http.Error(w, "Failed to encode user to JSON", http.StatusInternalServerError)
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
}

type CreateUserRequest struct {
	Email    string `json:"email"`
	Password string `json:"password"`
}

func (uh *userHandler) CreateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody CreateUserRequest
	defer r.Body.Close()
	if !uh.isValidCreateUserRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	token, err := uh.uuc.CreateUserAndToken(ctx, requestBody.Email, requestBody.Password)
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}

	w.Header().Set("Authorization", "Bearer "+token)
	w.WriteHeader(http.StatusOK)
}

func (uh *userHandler) isValidCreateUserRequest(body io.ReadCloser, requestBody *CreateUserRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body: %v", err)
		return false
	}
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

func (uh *userHandler) UpdateUser(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	var requestBody UpdateUserRequest
	defer r.Body.Close()
	if !uh.isValidUpdateUserRequest(r.Body, &requestBody) {
		w.WriteHeader(http.StatusBadRequest)
		return
	}

	if err := uh.uuc.UpdateUser(ctx, requestBody.Name); err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		return
	}
	w.WriteHeader(http.StatusOK)
}

func (uh *userHandler) isValidUpdateUserRequest(body io.ReadCloser, requestBody *UpdateUserRequest) bool {
	if err := json.NewDecoder(body).Decode(requestBody); err != nil {
		log.Error("Failed to decode request body: %v", err)
		return false
	}
	if requestBody.Name == "" {
		log.Warn("Invalid request body: %v", requestBody)
		return false
	}
	return true
}
