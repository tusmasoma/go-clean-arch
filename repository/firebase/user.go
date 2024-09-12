package firebase

import (
	"context"
	"fmt"

	"github.com/tusmasoma/go-clean-arch/entity"
	"github.com/tusmasoma/go-clean-arch/repository"
)

// DOCS: https://cloud.google.com/identity-platform/docs/use-rest-api?hl=ja

type userRepository struct {
	client *Client
}

func NewUserRepository(client *Client) repository.UserRepository {
	return &userRepository{
		client: client,
	}
}

func (ur *userRepository) Get(ctx context.Context, id string) (*entity.User, error) {
	user, err := ur.client.cli.GetUser(ctx, id)
	if err != nil {
		return nil, err
	}
	return &entity.User{
		ID:    user.UID,
		Name:  user.DisplayName,
		Email: user.Email,
	}, nil
}

type createUserRequest struct {
	Email             string `json:"email"`
	Password          string `json:"password"`
	ReturnSecureToken bool   `json:"returnSecureToken"`
}

type createUserResponse struct {
	IDToken      string `json:"idToken"`
	Email        string `json:"email"`
	RefreshToken string `json:"refreshToken"`
	ExpiresIn    string `json:"expiresIn"`
	LocalID      string `json:"localId"`
}

func (ur *userRepository) Create(ctx context.Context, user entity.User) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:signInWithPassword?key=%s", ur.client.apiKey)
	reqBody := createUserRequest{
		Email:             user.Email,
		Password:          user.Password,
		ReturnSecureToken: true,
	}

	var respBody createUserResponse
	if err := ur.client.postCall(ctx, url, reqBody, &respBody); err != nil {
		return err
	}
	return nil
}

type updateUserRequest struct {
	IDToken     string `json:"idToken"`
	DisplayName string `json:"displayName"`
	Email       string `json:"email"`
	Password    string `json:"password"`
}

type updateUserResponse struct {
	// 一部抜粋
	LocalID string `json:"localId"`
}

func (ur *userRepository) Update(ctx context.Context, user entity.User) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:update?key=%s", ur.client.apiKey)
	reqBody := updateUserRequest{
		IDToken:     user.ID, // TODO: change to IDToken
		DisplayName: user.Name,
		Email:       user.Email,
		Password:    user.Password,
	}

	var respBody updateUserResponse
	if err := ur.client.postCall(ctx, url, reqBody, &respBody); err != nil {
		return err
	}
	return nil
}

type deleteUserRequest struct {
	IDToken string `json:"idToken"`
}

type deleteUserResponse struct{}

func (ur *userRepository) Delete(ctx context.Context, idToken string) error {
	url := fmt.Sprintf("https://identitytoolkit.googleapis.com/v1/accounts:delete?key=%s", ur.client.apiKey)
	reqBody := deleteUserRequest{
		IDToken: idToken,
	}

	var respBody deleteUserResponse
	if err := ur.client.postCall(ctx, url, reqBody, &respBody); err != nil {
		return err
	}
	return nil
}

func (ur *userRepository) LockUserByEmail(_ context.Context, _ string) (bool, error) {
	// TODO: Firebase Auth does not provide a way to lock user by email
	return false, nil
}
