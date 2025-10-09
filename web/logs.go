// AI-assisted code
package main

import (
	"net/http"
	"strconv"
	"time"
)

func handleQueryLogs(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logID := r.PathValue("id")

	if _, err := dbGetDataset(session, userID, logID); err != nil {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}

	limit := 100
	if l := r.URL.Query().Get("limit"); l != "" {
		parsed, err := strconv.Atoi(l)
		if err != nil || parsed < 1 || parsed > 1000 {
			writeError(w, http.StatusBadRequest, "limit must be 1-1000")
			return
		}
		limit = parsed
	}

	var before, after *time.Time
	if b := r.URL.Query().Get("before"); b != "" {
		t, err := time.Parse(time.RFC3339, b)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid before timestamp")
			return
		}
		before = &t
	}
	if a := r.URL.Query().Get("after"); a != "" {
		t, err := time.Parse(time.RFC3339, a)
		if err != nil {
			writeError(w, http.StatusBadRequest, "invalid after timestamp")
			return
		}
		after = &t
	}

	entries, err := dbQueryLogs(session, userID, logID, limit, before, after)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query logs")
		return
	}
	if entries == nil {
		entries = []LogEntry{}
	}
	writeJSON(w, http.StatusOK, entries)
}
