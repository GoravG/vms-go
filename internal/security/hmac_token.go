package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"encoding/json"
	"errors"
	"strings"
	"vms_go/internal/config"
)

var (
	HMACSecretKey = getHMACKey()
)

type Claims struct {
	Email  string `json:"email"`
	UserID int    `json:"user_id"`
}

func getHMACKey() []byte {
	key := config.FromEnv().HMAC_SECRET_KEY
	return []byte(key)
}

// GenerateToken creates a token with claims, similar to JWT
func GenerateToken(claims Claims) (string, error) {
	// Convert claims to JSON
	payload, err := json.Marshal(claims)
	if err != nil {
		return "", err
	}

	// Base64 encode the payload
	encodedPayload := base64.RawURLEncoding.EncodeToString(payload)

	// Generate HMAC of the encoded payload
	mac := hmac.New(sha256.New, HMACSecretKey)
	mac.Write([]byte(encodedPayload))
	signature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	// Combine payload and signature with a dot separator
	return encodedPayload + "." + signature, nil
}

// VerifyAndExtractClaims verifies the token and returns the claims
func VerifyAndExtractClaims(token string) (*Claims, error) {
	// Split token into payload and signature
	parts := strings.Split(token, ".")
	if len(parts) != 2 {
		return nil, errors.New("invalid token format")
	}

	encodedPayload, receivedSignature := parts[0], parts[1]

	// Verify signature
	mac := hmac.New(sha256.New, HMACSecretKey)
	mac.Write([]byte(encodedPayload))
	expectedSignature := base64.RawURLEncoding.EncodeToString(mac.Sum(nil))

	if !hmac.Equal([]byte(receivedSignature), []byte(expectedSignature)) {
		return nil, errors.New("invalid signature")
	}

	// Decode payload
	payloadBytes, err := base64.RawURLEncoding.DecodeString(encodedPayload)
	if err != nil {
		return nil, err
	}

	// Parse claims
	var claims Claims
	if err := json.Unmarshal(payloadBytes, &claims); err != nil {
		return nil, err
	}

	return &claims, nil
}

// VerifyHMAC checks if the provided data matches the given HMAC
func VerifyHMAC(data []byte, providedHMAC string) bool {
	mac := hmac.New(sha256.New, HMACSecretKey)
	mac.Write(data)
	expectedMAC := hex.EncodeToString(mac.Sum(nil))
	return hmac.Equal([]byte(providedHMAC), []byte(expectedMAC))
}
