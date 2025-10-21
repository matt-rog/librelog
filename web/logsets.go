// AI-assisted code
package main

import (
	"encoding/json"
	"net/http"

	"github.com/gocql/gocql"
)

func handleListLogsets(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logsets, err := dbListLogsets(session, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list logsets")
		return
	}
	if logsets == nil {
		logsets = []Logset{}
	}
	writeJSON(w, http.StatusOK, logsets)
}

func handleCreateLogset(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	var req struct {
		Name        string `json:"name"`
		Description string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}
	if req.Name == "" {
		writeError(w, http.StatusBadRequest, "name required")
		return
	}

	logID := gocql.TimeUUID().String()
	if err := dbCreateLogset(session, userID, logID, req.Name, req.Description); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create logset")
		return
	}

	writeJSON(w, http.StatusCreated, Logset{
		LogID:       logID,
		UserID:      userID.String(),
		Name:        req.Name,
		Description: req.Description,
	})
}

func handleGetLogset(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logID := r.PathValue("id")

	logset, err := dbGetLogset(session, userID, logID)
	if err != nil {
		writeError(w, http.StatusNotFound, "logset not found")
		return
	}

	writeJSON(w, http.StatusOK, logset)
}

func handleUpdateLogset(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logID := r.PathValue("id")

	existing, err := dbGetLogset(session, userID, logID)
	if err != nil {
		writeError(w, http.StatusNotFound, "logset not found")
		return
	}

	var req struct {
		Name        *string `json:"name"`
		Description *string `json:"description"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeError(w, http.StatusBadRequest, "invalid json")
		return
	}

	name := existing.Name
	description := existing.Description
	if req.Name != nil {
		name = *req.Name
	}
	if req.Description != nil {
		description = *req.Description
	}

	if err := dbUpdateLogset(session, userID, logID, name, description); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update logset")
		return
	}

	writeJSON(w, http.StatusOK, Logset{
		LogID:       logID,
		UserID:      userID.String(),
		Name:        name,
		Description: description,
	})
}

func handleDeleteLogset(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logID := r.PathValue("id")

	if _, err := dbGetLogset(session, userID, logID); err != nil {
		writeError(w, http.StatusNotFound, "logset not found")
		return
	}

	if err := dbDeleteLogset(session, userID, logID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete logset")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
