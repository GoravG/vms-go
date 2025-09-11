package middleware

// func Auth(next http.Handler, secret []byte) http.Handler {
// 	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
// 		authHeader := r.Header.Get("Authorization")
// 		if authHeader == "" {
// 			http.Error(w, "missing token", http.StatusUnauthorized)
// 			return
// 		}

// 		tokenStr := strings.TrimPrefix(authHeader, "Bearer ")
// 		username, err := token.Validate(tokenStr, secret)
// 		if err != nil {
// 			http.Error(w, err.Error(), http.StatusUnauthorized)
// 			return
// 		}

// 		// Store username in context (optional for later use)
// 		ctx := r.Context()
// 		ctx = WithUser(ctx, username) // we'll define this helper
// 		next.ServeHTTP(w, r.WithContext(ctx))
// 	})
// }
