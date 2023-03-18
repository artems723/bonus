package handler

import (
	"bonus/internal/model"
	"bonus/internal/repository"
	"context"
	"encoding/json"
	"errors"
	"github.com/golang-jwt/jwt/v5"
	"log"
	"net/http"
	"time"
)

func (h *Handler) RegisterHandler(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	// Read JSON and store to user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	// Check errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Create user
	err = h.userService.Create(r.Context(), user)
	if err != nil && !errors.Is(err, repository.ErrLoginIsTaken) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, repository.ErrLoginIsTaken) {
		http.Error(w, err.Error(), http.StatusConflict)
		return
	}
	log.Printf("User was registered: %v\n", user.Login)

	token, err := Login(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

func (h *Handler) LoginHandler(w http.ResponseWriter, r *http.Request) {
	var user *model.User
	// Read JSON and store to user struct
	err := json.NewDecoder(r.Body).Decode(&user)
	// Check errors
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}
	// Check user was registered and password ok
	ok, err := h.userService.CheckUser(r.Context(), user)
	if err != nil && !errors.Is(err, repository.ErrUserNotFound) {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	if errors.Is(err, repository.ErrUserNotFound) {
		http.Error(w, err.Error(), http.StatusUnauthorized)
		return
	}
	if !ok {
		http.Error(w, ErrWrongPassword.Error(), http.StatusUnauthorized)
		return
	}
	log.Printf("User was authenticated: %v\n", user.Login)

	token, err := Login(user)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}
	w.Header().Add("Authorization", token)
	w.WriteHeader(http.StatusOK)
}

var jwtSecretKey = []byte("J6Zao6mOK#j^bUJVcG")

func Login(user *model.User) (string, error) {
	payload := jwt.MapClaims{
		"sub": user.Login,
		"exp": time.Now().Add(time.Hour * 72).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, payload)
	return token.SignedString(jwtSecretKey)
}

func ParseToken(accessToken string) (string, error) {
	token, err := jwt.Parse(accessToken, func(token *jwt.Token) (interface{}, error) {
		_, ok := token.Method.(*jwt.SigningMethodHMAC)
		if !ok {
			return "", ErrWrongAuthToken
		}
		return jwtSecretKey, nil
	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		return claims["sub"].(string), nil
	} else {
		return "", ErrWrongAuthToken
	}
}

// Authenticator middleware
func Authenticator(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		accessToken := r.Header.Get("Authorization")
		if accessToken == "" {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, "token is empty", http.StatusUnauthorized)
			return
		}

		login, err := ParseToken(accessToken)
		if err != nil {
			w.Header().Set("WWW-Authenticate", `Basic realm="restricted", charset="UTF-8"`)
			http.Error(w, err.Error(), http.StatusUnauthorized)
			return
		}

		// set username to context
		ctx := r.Context()
		ctx = context.WithValue(ctx, "login", login)
		r = r.WithContext(ctx)
		log.Printf("user %s was authenticated\n", login)
		// authenticated, pass it through
		next.ServeHTTP(w, r)
	})
}

var ErrWrongPassword = errors.New("wrong password")
var ErrWrongAuthToken = errors.New("wrong auth token")
