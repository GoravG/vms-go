package security

import (
	"crypto/hmac"
	"crypto/sha256"
	"encoding/hex"
	"vms_go/internal/config"
)

var (
	HMACSecretKey = getHMACKey()
)

func getHMACKey() []byte {
	key := config.FromEnv().HMAC_SECRET_KEY
	return []byte(key)
}
func GenerateHMAC(data []byte) string {
	mac := hmac.New(sha256.New, HMACSecretKey)
	mac.Write(data)
	return hex.EncodeToString(mac.Sum(nil))
}
