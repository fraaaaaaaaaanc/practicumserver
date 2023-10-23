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
	"practicumserver/internal/handlers/allhandlers"
	storage "practicumserver/internal/storage/pg"
	"time"
)

const LEN_ID = 16
const TOKEN_EXP = time.Hour * 5
const KEY = "key"

type Claims struct {
	jwt.RegisteredClaims
	UserID string
}

func NewUserId(hndlr *handlers.Handlers) (string, error) {
	userIDBytes := make([]byte, LEN_ID)
	for {
		_, err := rand.Read(userIDBytes)
		if err != nil {
			return "", err
		}
		userID := base64.StdEncoding.EncodeToString(userIDBytes)
		if DBStorage, ok := hndlr.Storage.(*storage.DBStorage); ok {
			var exists bool
			row := DBStorage.DB.QueryRowContext(context.Background(),
				"SELECT EXISTS (SELECT 1 FROM links WHERE userid = $1)",
				userID)
			if err = row.Scan(&exists); err != nil {
				return "", err
			}
			if !exists {
				return userID, nil
			}
		}
	}
	return "", errors.New("nil userid")
}

func BuildJWTString(userID string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, Claims{
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(time.Now().Add(TOKEN_EXP)),
		},
		UserID: userID,
	})
	return token.SignedString([]byte(KEY))
}

func GetUserID(tokenString string) (string, error) {
	claims := &Claims{}
	token, err := jwt.ParseWithClaims(tokenString, claims, func(token *jwt.Token) (interface{}, error) {
		return []byte(KEY), nil
	})
	if err != nil {
		return "", err
	}
	if !token.Valid {
		return "", errors.New("Token is not valid")
	}
	return claims.UserID, nil
}

func MiddlewareCheckCookies(log *zap.Logger, hndlr *handlers.Handlers) func(h http.Handler) http.Handler {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			var userID string
			token, err := r.Cookie("Authorization")
			if err != nil {
				userID, err = NewUserId(hndlr)
				fmt.Println(userID)
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
				cookie := http.Cookie{
					Name:     "Authorization",
					Value:    tokenString,
					Path:     "/",
					MaxAge:   18000,
					HttpOnly: true,
				}
				http.SetCookie(w, &cookie)
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			userID, err = GetUserID(token.Value)
			fmt.Println(userID)
			if err != nil {
				w.WriteHeader(http.StatusInternalServerError)
				log.Error("Error:", zap.Error(err))
				return
			}
			ctx := context.WithValue(r.Context(), "userID", userID)
			r = r.WithContext(ctx)

			h.ServeHTTP(w, r)
		})
	}
}
