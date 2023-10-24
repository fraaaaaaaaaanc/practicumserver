package cookie

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
	"practicumserver/internal/models"
	"time"
)

const SecretKey = "key"
const TokenExp = time.Hour * 6
const MaxLenUserID = 16

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

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
		fmt.Println(existsUserID)
		if existsUserID {
			return userID, nil
		}
	}
}

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

func GetUserID(tokenString string) (string, error) {
	var claims Claims
	token, err := jwt.ParseWithClaims(tokenString, &claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(SecretKey), nil
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
				userID, err = GetUserID(cookie.Value)
				if err != nil {
					w.WriteHeader(http.StatusInternalServerError)
					log.Error("Error:", zap.Error(err))
					return
				}
			}
			ctx := context.WithValue(r.Context(), models.UserIDKey, userID)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}
