// Package cookie provides functions and middleware for handling user authentication and session management using cookies.
//
// This package includes functionality for generating and validating JWT tokens, creating unique user IDs, and checking
// user authentication through cookies. It offers a complete solution for user authentication and managing session
// data within web applications.
package cookie

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
	handlers "practicumserver/internal/handlers/allhandlers"
	"practicumserver/internal/models"
	"time"
)

// SecretKey is the secret key used for signing and verifying JWT tokens.
const SecretKey = "key"

// TokenExp represents the token expiration time.
const TokenExp = time.Hour * 6

// MaxLenUserID defines the maximum length of the generated user ID.
const MaxLenUserID = 16

// Claims represents custom JWT claims that include the UserID.
type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

// NewUserID generates a new user ID by creating random bytes and checking its uniqueness.
func NewUserID(ctx context.Context, hndlrs *handlers.Handlers) (string, error) {
	userIDbytes := make([]byte, MaxLenUserID)
	for {
		_, err := rand.Read(userIDbytes)
		if err != nil {
			return "", err
		}
		userID := base64.StdEncoding.EncodeToString(userIDbytes)
		existsUserID, err := hndlrs.Storage.CheckUserID(ctx, userID)
		if err != nil {
			return "", err
		}
		if existsUserID {
			return userID, nil
		}
	}
}

// BuildJWTString creates a JWT token with the provided user ID and signs it with the SecretKey.
func BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TokenExp)),
		},
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(SecretKey))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

// GetUserID retrieves the user ID from a JWT token.
func GetUserID(tokenString string) (string, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("token is not valid")
	}
	return claims.UserID, nil
}

// MiddlewareCheckCookie is a middleware function that checks and manages user authentication via cookies.
func MiddlewareCheckCookie(log *zap.Logger, hndlrs *handlers.Handlers) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				// If no Authorization cookie is found, generate a new user ID and create a JWT token.
				userID, err = NewUserID(r.Context(), hndlrs)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Error("Error:", zap.Error(err))
					return
				}
				tokenString, err := BuildJWTString(userID)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Error("Error:", zap.Error(err))
					return
				}
				newCookie := &http.Cookie{
					Name:     "Authorization",
					Value:    tokenString,
					Path:     "/",
					MaxAge:   7200,
					HttpOnly: true,
				}
				http.SetCookie(w, newCookie)
			} else {
				// If an Authorization cookie is found, validate it and extract the user ID.
				userID, err = GetUserID(cookie.Value)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Error("Error:", zap.Error(err))
					return
				}
			}
			// Set the user ID in the request context and continue processing.
			ctx := context.WithValue(r.Context(), models.UserIDKey, userID)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}
