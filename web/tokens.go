// AI-assisted code
package main

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"net/http"
)

func handleListTokens(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	keys, err := dbListTokens(session, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list tokens")
		return
	}
	if keys == nil {
		keys = []APIKey{}
	}
	writeJSON(w, http.StatusOK, keys)
}

func handleCreateToken(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)

	var req struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name required")
		return
	}

	tokenBytes := make([]byte, 32)
	if _, err := rand.Read(tokenBytes); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to generate token")
		return
	}
	token := hex.EncodeToString(tokenBytes)
	tokenHash := hashSHA256(token)
	prefix := token[:8]

	if err := dbCreateToken(session, tokenHash, userID, req.Name, prefix); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create token")
		return
	}

	writeJSON(w, http.StatusCreated, map[string]string{
		"token":  token,
		"name":   req.Name,
		"prefix": prefix,
	})
}

func handleDeleteToken(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	tokenHash := r.PathValue("hash")

	if err := dbDeleteToken(session, tokenHash, userID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete token")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
