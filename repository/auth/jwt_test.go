package auth

import (
	"reflect"
	"testing"

	jwtgo "github.com/dgrijalva/jwt-go"
	"github.com/google/uuid"
)

func Test_JWTToken(t *testing.T) {
	userID := uuid.MustParse("f6db2530-cd9b-4ac1-8dc1-38c795e6eec2")
	email := "test@gmail.com"

	repo := NewAuthRepository()

	// GenerateToken test
	jwt, jti := repo.GenerateToken(userID.String(), email)

	// JWTのフォーマットが正しいことを確認
	token, err := jwtgo.Parse(jwt, func(token *jwtgo.Token) (interface{}, error) {
		// ここで公開キーを使って署名を検証する（公開キーは環境に依存する）
		return loadPublicKey(rawPublicKey)
	})
	if err != nil {
		t.Errorf("Failed to parse JWT: %s", err)
	}

	// クレームを検証
	claims, ok := token.Claims.(jwtgo.MapClaims)
	if !ok {
		t.Errorf("Failed to parse claims")
	}

	if claims["Email"] != email {
		t.Errorf("Expected email %s, got %s", email, claims["Email"])
	}

	if claims["jti"] != jti {
		t.Errorf("Expected JTI %s, got %s", jti, claims["jti"])
	}

	// ValidateAccessToken test
	err = repo.ValidateAccessToken(jwt)
	if err != nil {
		t.Errorf("Failed to ValidateAccessToken: %s", err)
	}

	// GetPayloadFromToken test
	payload, err := repo.GetPayloadFromToken(jwt)
	if err != nil {
		t.Errorf("Failed to GetPayloadFromToken: %s", err)
	}
	wantPayload := map[string]string{
		"jti":    jti,
		"userId": userID.String(),
		"Email":  email,
	}
	if !reflect.DeepEqual(payload, wantPayload) {
		t.Errorf("GetPayloadFromToken() \n got = %v,\n want = %v", payload, wantPayload)
	}
}
