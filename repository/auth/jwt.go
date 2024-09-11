package auth

import (
	"crypto"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"strings"

	"github.com/google/uuid"

	"github.com/tusmasoma/go-clean-arch/repository"
)

var (
	rawPublicKey = []byte(`-----BEGIN PUBLIC KEY-----
MIIBIjANBgkqhkiG9w0BAQEFAAOCAQ8AMIIBCgKCAQEAweNk5DXY50UZ5Vvd/glB
I1H3jrr74Ns3hH+9Btq5UWTogMimJk7ouEjrBKABmYGXNUEs9tCWzFDZ6wJW0nmL
XhlLyRXr6GDwu3hcQA5HTAA41vio+IGqdEWn5cV2woD2BlFjuXvx5CHaOMaZybYR
KEf1b9r6lB15WcxEDsfH7aZQp6RKG4ufI9C3fZif7VsYpRPC8B6qWvAdN871kJ0A
mBs4ZYCXQVwV+vN5liq5YG1V+Ju00wjdozs1zGsVgRot1YLvfPPQQ+fZj0OWzTyx
bwIz8+VOTM/WXSdEpk3d6q6Vcf78LRD4SJlbG0Ru/l3Y9P/CAxml2k83yjV2g55O
ewIDAQAB
-----END PUBLIC KEY-----`)

	rawSecretKey = []byte(`-----BEGIN PRIVATE KEY-----
MIIEvgIBADANBgkqhkiG9w0BAQEFAASCBKgwggSkAgEAAoIBAQDB42TkNdjnRRnl
W93+CUEjUfeOuvvg2zeEf70G2rlRZOiAyKYmTui4SOsEoAGZgZc1QSz20JbMUNnr
AlbSeYteGUvJFevoYPC7eFxADkdMADjW+Kj4gap0RaflxXbCgPYGUWO5e/HkIdo4
xpnJthEoR/Vv2vqUHXlZzEQOx8ftplCnpEobi58j0Ld9mJ/tWxilE8LwHqpa8B03
zvWQnQCYGzhlgJdBXBX683mWKrlgbVX4m7TTCN2jOzXMaxWBGi3Vgu9889BD59mP
Q5bNPLFvAjPz5U5Mz9ZdJ0SmTd3qrpVx/vwtEPhImVsbRG7+Xdj0/8IDGaXaTzfK
NXaDnk57AgMBAAECggEANb8KPbSrjth56DmCHSSFmPvkyv0MS3WZOzKJvLeu/WAi
j2iPnjjrjAIym9Ka/umMd+e8RiLmWnbjIaFBXhDxUEFk37Yi8gTFVsJzmBIdM3Uw
TG8br9+J17djZm9Jj3teN7wiD83K7PlxW6G6Cc9djDP+VmZ2Zc6R0BGuoAZDZp3i
7eflreYzuCha+nv4nIRcDayegIeccMW7BGzOpyB2GGXxjsFtZHYl597ZyL8RwKMm
+eBqHaa0BHAYenEWTM2Q76ctnRj0jAqaaYVT9EqUU/05KVIZOBLjA1e6ex3OUY1T
sn01i4oyWSvzA3EhjRmt6DuEDT5qK+AabY9VCDUCSQKBgQDy1B6YrENOVtN5BKJv
+IfRFCNHVCCjQ1ob2XEXS6QniHeWao+793q9vILOSd1USGeeIOwSysgO3LnCBH3R
j7Qrus35uE3l0gv4KKpaedtt3DaICubTvJXq15uZYpSJUiZLMVsS3JyUjLEPG+In
SkhIRkGBAIBPGZiod8eQzV0WfwKBgQDMZ7JXPOvTEkpEI7kV1oPS38WSYQvICpr0
VPEjIqv5yJNF7MA9J9BCude5m0mDQPWd9MBU2i0ycUdeTEe452Gdh8JfEedY+nlG
pXbfuMnw3MvaQ3GSBfWIZiAOCLDQJRmVdhrM/nIi80TG09pwMvmNGsbRS9cGijtI
hFvaFDAiBQKBgCRXOnz+ytPeiqeB2g2H1EumB+GU5Y2JduLUF+i0mUyRT9Ri/j/T
ObtLiwf0ZftHGrq/kpT9ZBNVVTeEFJBYQU6KFmlY+895L/FjpJsFwaEfY8nYV9M4
VfdfbRn3duNWOATozgh0m7pfk9/+/EmFBGxMl2EHAizUV9RemK9DDLthAoGBALWJ
f0WlcJhkTRsZUv9HJoq5fMIVeJ4wdRB9BDC9UVmlPs9ChjWKT5eDcEmC1hZBMiMY
RVzW7H85RjZErwpUTUjYUtOWlg5bXixVNi9Z8df+cPonHg2fR0Ld2Kg+JbKm0IMC
gqj/bqUFw1aGvyEY1LPyTRODNLS1PhOYoe8cMOd1AoGBAONqJMqUIjNQqw7cYppu
5+jFR/oAdlz6R4qHM/X2rOp0pXK1U3Z3wVtdpc2DPyQh1xXyrh6m1xr8Vn+wzihQ
0dH8t0SocHZOazOB3mshL8O+16HoRHlmCbYVdaG7CUZPcT+vhiO0dozsj8jHj4Iv
9Goa9sCO+rKNzDKL7ZONt1/p
-----END PRIVATE KEY-----`)
)

type authRepository struct{}

func NewAuthRepository() repository.AuthRepository {
	return &authRepository{}
}

// type Payload struct {
// 	JTI    string `json:"jti"`
// 	UserID string `json:"userId"`
// }

const expectedTokenParts = 3

func loadPrivateKey(keyBytes []byte) (*rsa.PrivateKey, error) {
	// PEMエンコードされたデータからPEMブロックをデコード
	block, _ := pem.Decode(keyBytes)
	if block == nil || (block.Type != "RSA PRIVATE KEY" && block.Type != "PRIVATE KEY") {
		return nil, fmt.Errorf("failed to decode PEM block containing the key")
	}

	// PEMブロックからRSA秘密鍵をパース
	privInterface, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse private key: %w", err)
	}

	privKey, ok := privInterface.(*rsa.PrivateKey)
	if !ok {
		return nil, fmt.Errorf("not RSA private key")
	}

	return privKey, nil
}

func loadPublicKey(keyBytes []byte) (*rsa.PublicKey, error) {
	// PEMエンコードされたデータからPEMブロックをデコード
	block, _ := pem.Decode(keyBytes)
	if block == nil || block.Type != "PUBLIC KEY" {
		return nil, fmt.Errorf("failed to decode PEM block containing the key")
	}

	// PEMブロックからRSA公開鍵をパース
	pubInterface, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, fmt.Errorf("failed to parse public key: %w", err)
	}

	pubKey, ok := pubInterface.(*rsa.PublicKey)
	if !ok {
		return nil, fmt.Errorf("not RSA public key")
	}

	return pubKey, nil
}

// Base64Url Encode
func base64UrlEncode(b []byte) string {
	return base64.RawURLEncoding.EncodeToString(b)
}

// Base64Url Decode
func base64UrlDecode(s string) ([]byte, error) {
	return base64.RawURLEncoding.DecodeString(s)
}

// アクセストークン(JWT形式)の生成
func (ar *authRepository) GenerateToken(userID, email string) (string, string) {
	// ヘッダの作成
	header := map[string]string{
		"typ": "JWT",
		"alg": "RS256",
	}
	headerBytes, _ := json.Marshal(header)
	encodedHeader := base64UrlEncode(headerBytes)

	// ペイロードの作成
	jti := uuid.New().String()
	payload := map[string]string{
		"jti":    jti,
		"userId": userID,
		"Email":  email,
	}
	payloadBytes, _ := json.Marshal(payload)
	encodedPayload := base64UrlEncode(payloadBytes)

	// エンコードされたヘッダとペイロードを結合
	jwtWithoutSignature := fmt.Sprintf("%s.%s", encodedHeader, encodedPayload)

	// SHA-256ハッシュを計算
	hashed := sha256.Sum256([]byte(jwtWithoutSignature))

	// 署名作成
	privKey, err := loadPrivateKey(rawSecretKey)
	if err != nil {
		panic(err)
	}
	signature, err := rsa.SignPKCS1v15(rand.Reader, privKey, crypto.SHA256, hashed[:])
	if err != nil {
		panic(err)
	}
	encodedSignature := base64UrlEncode(signature)

	// JWTを完成
	jwt := fmt.Sprintf("%s.%s", jwtWithoutSignature, encodedSignature)

	return jwt, jti
}

func (ar *authRepository) ValidateAccessToken(jwt string) error {
	// アクセストークンの検証
	parts := strings.Split(jwt, ".")
	if len(parts) != expectedTokenParts {
		return fmt.Errorf("invalid token")
	}
	// エンコードされたヘッダとペイロードを結合
	jwtWithoutSignature := fmt.Sprintf("%s.%s", parts[0], parts[1])
	// SHA-256ハッシュを計算
	hashed := sha256.Sum256([]byte(jwtWithoutSignature))

	// 著名作成
	signature, err := base64UrlDecode(parts[2])
	if err != nil {
		return fmt.Errorf("decoding failed: %w", err)
	}

	// 検証
	pubKey, err := loadPublicKey(rawPublicKey)
	if err != nil {
		return err
	}

	err = rsa.VerifyPKCS1v15(pubKey, crypto.SHA256, hashed[:], signature)
	log.Print(err)
	if err != nil {
		return fmt.Errorf("signature verification failed: %w", err)
	}

	return nil
}

func (ar *authRepository) GetPayloadFromToken(jwt string) (map[string]string, error) {
	var emptyPayload map[string]string
	// アクセストークンの検証
	parts := strings.Split(jwt, ".")
	if len(parts) != expectedTokenParts {
		return emptyPayload, fmt.Errorf("invalid token")
	}
	// エンコードされたヘッダとペイロードを結合
	encodedPayload := parts[1]
	// Base64Urlデコード
	payloadBytes, err := base64UrlDecode(encodedPayload)
	if err != nil {
		return emptyPayload, fmt.Errorf("decoding failed: %w", err)
	}

	// JSONデコード
	var payload map[string]string
	err = json.Unmarshal(payloadBytes, &payload)
	if err != nil {
		return emptyPayload, fmt.Errorf("JSON unmarshalling failed")
	}

	return payload, nil
}
