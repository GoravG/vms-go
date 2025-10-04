package token

import (
	"sync"
)

var (
	currentToken string
	mutex        sync.RWMutex
)

func SetToken(token string) {
	mutex.Lock()
	defer mutex.Unlock()
	currentToken = token
}

func GetToken() string {
	mutex.RLock()
	defer mutex.RUnlock()
	return currentToken
}
