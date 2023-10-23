package coockie

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"github.com/golang-jwt/jwt/v4"
	"go.uber.org/zap"
	"net/http"
	handlers "practicumserver/internal/handlers/allhandlers"
	"time"
)

const SECRET_KEY = "key"
const TOKEN_EXP = time.Hour * 6
const MAX_LEN_USER_ID = 16

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func NewUserID(ctx context.Context, hndlrs *handlers.Handlers) (string, error) {
	userIDbytes := make([]byte, MAX_LEN_USER_ID)
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
		fmt.Println(existsUserID)
		if existsUserID {
			return userID, nil
		}
	}
}

func BuildJWTString(ctx context.Context, hndlrs *handlers.Handlers, userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})
	tokenString, err := token.SignedString([]byte(SECRET_KEY))
	if err != nil {
		return "", err
	}
	return tokenString, err
}

func GetUserID(tokenString string) (string, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SECRET_KEY), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("Token is not valid")
	}
	return claims.UserID, nil
}

func MiddlewareCheckCoockie(log *zap.Logger, hndlrs *handlers.Handlers) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string
			cookie, err := r.Cookie("Authorization")
			if err != nil {
				fmt.Println(err)
				userID, err = NewUserID(r.Context(), hndlrs)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Error("Error:", zap.Error(err))
					return
				}
				tokenString, err := BuildJWTString(r.Context(), hndlrs, userID)
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
				userID, err = GetUserID(cookie.Value)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Error("Error:", zap.Error(err))
					return
				}
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}
