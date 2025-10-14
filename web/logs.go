// AI-assisted code
package main

import (
	"encoding/csv"
	"encoding/json"
	"fmt"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/gocql/gocql"
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

func collectAllLogs(userID gocql.UUID, logID string) ([]LogEntry, error) {
	var all []LogEntry
	const batch = 1000
	var before *time.Time
	for {
		entries, err := dbQueryLogs(session, userID, logID, batch, before, nil)
		if err != nil {
			return nil, err
		}
		all = append(all, entries...)
		if len(entries) < batch {
			break
		}
		before = &entries[len(entries)-1].RecvTime
	}
	return all, nil
}

func handleExportLogs(w http.ResponseWriter, r *http.Request) {
	userID := getUserID(r)
	logID := r.PathValue("id")

	ds, err := dbGetDataset(session, userID, logID)
	if err != nil {
		writeError(w, http.StatusNotFound, "dataset not found")
		return
	}

	entries, err := collectAllLogs(userID, logID)
	if err != nil {
		writeError(w, http.StatusInternalServerError, "failed to query logs")
		return
	}
	if entries == nil {
		entries = []LogEntry{}
	}

	format := r.URL.Query().Get("format")
	if format == "csv" {
		exportCSV(w, ds.Name, entries)
	} else {
		exportJSON(w, ds.Name, entries)
	}
}

type exportEntry struct {
	RecvTime time.Time       `json:"recv_time"`
	Data     json.RawMessage `json:"data"`
}

func exportJSON(w http.ResponseWriter, name string, entries []LogEntry) {
	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.json"`, name))

	out := make([]exportEntry, len(entries))
	for i, e := range entries {
		out[i].RecvTime = e.RecvTime
		if json.Valid([]byte(e.Data)) {
			out[i].Data = json.RawMessage(e.Data)
		} else {
			b, _ := json.Marshal(e.Data)
			out[i].Data = b
		}
	}
	json.NewEncoder(w).Encode(out)
}

func exportCSV(w http.ResponseWriter, name string, entries []LogEntry) {
	w.Header().Set("Content-Type", "text/csv")
	w.Header().Set("Content-Disposition", fmt.Sprintf(`attachment; filename="%s.csv"`, name))

	// collect all keys across entries
	keySet := map[string]bool{}
	parsed := make([]map[string]interface{}, len(entries))
	for i, e := range entries {
		var m map[string]interface{}
		if err := json.Unmarshal([]byte(e.Data), &m); err != nil {
			m = map[string]interface{}{"data": e.Data}
		}
		parsed[i] = m
		for k := range m {
			keySet[k] = true
		}
	}

	cols := make([]string, 0, len(keySet))
	for k := range keySet {
		cols = append(cols, k)
	}
	sort.Strings(cols)

	cw := csv.NewWriter(w)
	header := append([]string{"time"}, cols...)
	cw.Write(header)

	for i, e := range entries {
		row := []string{e.RecvTime.Format(time.RFC3339)}
		for _, col := range cols {
			v, ok := parsed[i][col]
			if !ok {
				row = append(row, "")
				continue
			}
			switch val := v.(type) {
			case string:
				row = append(row, val)
			case float64:
				row = append(row, strconv.FormatFloat(val, 'f', -1, 64))
			case bool:
				row = append(row, strconv.FormatBool(val))
			default:
				b, _ := json.Marshal(val)
				row = append(row, string(b))
			}
		}
		cw.Write(row)
	}
	cw.Flush()
}
