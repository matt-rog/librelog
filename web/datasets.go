// AI-assisted code
package main

import (
	"encoding/json"
	"net/http"

	"github.com/gocql/gocql"
)

func handleListDatasets(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	datasets, err := dbListDatasets(session, userID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to list datasets")
		return
	}
	if datasets == nil {
		datasets = []Dataset{}
	}
	writeJSON(w, http.StatusOK, datasets)
}

func handleCreateDataset(w http.ResponseWriter, r *http.Request) {
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
	if err := dbCreateDataset(session, userID, logID, req.Name, req.Description); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to create dataset")
		return
	}

	writeJSON(w, http.StatusCreated, Dataset{
		LogID:       logID,
		UserID:      userID.String(),
		Name:        req.Name,
		Description: req.Description,
	})
}

func handleGetDataset(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logID := r.PathValue("id")

	dataset, err := dbGetDataset(session, userID, logID)
	if err != nil {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}

	writeJSON(w, http.StatusOK, dataset)
}

func handleUpdateDataset(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logID := r.PathValue("id")

	existing, err := dbGetDataset(session, userID, logID)
	if err != nil {
		writeError(w, http.StatusNotFound, "dataset not found")
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

	if err := dbUpdateDataset(session, userID, logID, name, description); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to update dataset")
		return
	}

	writeJSON(w, http.StatusOK, Dataset{
		LogID:       logID,
		UserID:      userID.String(),
		Name:        name,
		Description: description,
	})
}

func handleDeleteDataset(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logID := r.PathValue("id")

	if _, err := dbGetDataset(session, userID, logID); err != nil {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}

	if err := dbDeleteDataset(session, userID, logID); err != nil {
		writeError(w, http.StatusInternalServerError, "failed to delete dataset")
		return
	}

	writeJSON(w, http.StatusOK, map[string]string{"status": "ok"})
}
