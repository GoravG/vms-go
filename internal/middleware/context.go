package middleware

import (
	"context"
)

type ctxKey string

var userKey ctxKey = "username"

func WithUser(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, userKey, username)
}

func GetUser(ctx context.Context) (string, bool) {
	u, ok := ctx.Value(userKey).(string)
	return u, ok
}
