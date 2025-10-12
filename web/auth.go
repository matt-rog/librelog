// AI-assisted code
package main

import (
	"context"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/gocql/gocql"
	"golang.org/x/crypto/bcrypt"
)

type contextKey string

const userIDKey contextKey = "user_id"

func getUserID(r *http.Request) gocql.UUID {
	return r.Context().Value(userIDKey).(gocql.UUID)
}

func requireAuth(next http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		auth := r.Header.Get("Authorization")
		if !strings.HasPrefix(auth, "Bearer ") {
			writeError(w, http.StatusUnauthorized, "missing token")
			return
		}
		token := strings.TrimPrefix(auth, "Bearer ")
		tokenHash := hashSHA256(token)

		userID, err := dbGetUserByToken(session, tokenHash)
		if err != nil {
			writeError(w, http.StatusUnauthorized, "invalid token")
			return
		}

		ctx := context.WithValue(r.Context(), userIDKey, userID)
		next(w, r.WithContext(ctx))
	}
}

func hashSHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func generateAccountNumber() (string, error) {
	b := make([]byte, 5)
	if _, err := rand.Read(b); err != nil {
		return "", err
	}
	num := 0
	for _, v := range b {
		num = num*256 + int(v)
	}
	return fmt.Sprintf("%010d", num%10000000000), nil
}

func handleSignup(w http.ResponseWriter, r *http.Request) {
	var req struct {
		Password string `json:"password"`
		Name     string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Password == "" {
		writeError(w, http.StatusBadRequest, "password required")
		return
	}

	accountNumber, err := generateAccountNumber()
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate account number")
		return
	}

	accountHash := hashSHA256(accountNumber)
	passwordHash, err := bcrypt.GenerateFromPassword([]byte(req.Password), bcrypt.DefaultCost)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to hash password")
		return
	}

	userID := gocql.TimeUUID()
	if err := dbCreateUser(session, userID, accountHash, string(passwordHash), req.Name); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create user")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"account_number": accountNumber,
	})
}

func handleLogin(w http.ResponseWriter, r *http.Request) {
	var req struct {
		AccountNumber string `json:"account_number"`
		Password      string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.AccountNumber == "" || req.Password == "" {
		writeError(w, http.StatusBadRequest, "account_number and password required")
		return
	}

	accountHash := hashSHA256(req.AccountNumber)
	userID, err := dbGetUserIDByAccount(session, accountHash)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	user, err := dbGetUser(session, userID)
	if err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(req.Password)); err != nil {
		writeError(w, http.StatusUnauthorized, "invalid credentials")
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	token := hex.EncodeToString(tokenBytes)
	tokenHash := hashSHA256(token)

	if err := dbCreateToken(session, tokenHash, userID, "", ""); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{
		"token": token,
	})
}

func handleLogout(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	token := strings.TrimPrefix(auth, "Bearer ")
	tokenHash := hashSHA256(token)

	userID := getUserID(r)
	if err := dbDeleteToken(session, tokenHash, userID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
