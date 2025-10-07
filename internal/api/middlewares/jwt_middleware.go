package middlewares

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"qiyana_vtu/pkg/utils"
	"strings"

	"github.com/golang-jwt/jwt/v5"
)

func JWTMiddleware(next http.Handler) http.Handler {
	fmt.Println("----------------- JWT middleware ------------------")

	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		fmt.Println("+++++++++++ inside JWT middleware")

		cookie, err := r.Cookie("Bearer")
		if err != nil {
			http.Error(w, "Unauthorized: Missing Bearer token", http.StatusUnauthorized)
			return
		}

		token := strings.TrimPrefix(cookie.Value, "Bearer ")

		jwtSecret := os.Getenv("JWT_SECRET")

		parsedToken, err := jwt.Parse(token, func(token *jwt.Token) (any, error) {

			return []byte(jwtSecret), nil
		}, jwt.WithValidMethods([]string{jwt.SigningMethodHS256.Alg()}))

		if err != nil {
			if errors.Is(err, jwt.ErrTokenExpired) {
				http.Error(w, "token expired", http.StatusUnauthorized)
				return
			}
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		if parsedToken.Valid {
			log.Println("Valid JWT")
		} else {
			http.Error(w, "invalid login token", http.StatusUnauthorized)
			log.Println("invalid JWT:", token)
		}

		claims, ok := parsedToken.Claims.(jwt.MapClaims)
		if !ok {
			http.Error(w, "invalid login token", http.StatusUnauthorized)
			log.Println("invalid login token:", token)
			return
		}

		ctx := context.WithValue(r.Context(), utils.ContextKey("role"), claims["role"])
		ctx = context.WithValue(ctx, utils.ContextKey("expiresAt"), claims["exp"])
		ctx = context.WithValue(ctx, utils.ContextKey("username"), claims["user"])
		ctx = context.WithValue(ctx, utils.ContextKey("userId"), claims["uid"])

		next.ServeHTTP(w, r.WithContext(ctx))
		fmt.Println("sent response from JWT middleware")
	})
}
