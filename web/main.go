// AI-assisted code
package main

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"

	"github.com/gocql/gocql"
)

var session *gocql.Session

func writeJSON(w http.ResponseWriter, status int, v interface{}) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(status)
	json.NewEncoder(w).Encode(v)
}

func writeError(w http.ResponseWriter, status int, msg string) {
	writeJSON(w, status, map[string]string{"error": msg})
}

func main() {
	urls := strings.Split(os.Getenv("CASSANDRA_CLUSTER"), " ")
	cluster := gocql.NewCluster(urls...)
	cluster.Keyspace = "librelog"
	cluster.Consistency = gocql.Quorum

	var err error
	session, err = cluster.CreateSession()
	if err != nil {
		log.Fatal("Failed to connect to Cassandra:", err)
	}
	defer session.Close()

	mux := http.NewServeMux()

	mux.HandleFunc("POST /api/signup", handleSignup)
	mux.HandleFunc("POST /api/login", handleLogin)
	mux.HandleFunc("POST /api/logout", requireAuth(handleLogout))

	mux.HandleFunc("GET /api/datasets", requireAuth(handleListDatasets))
	mux.HandleFunc("POST /api/datasets", requireAuth(handleCreateDataset))
	mux.HandleFunc("GET /api/datasets/{id}", requireAuth(handleGetDataset))
	mux.HandleFunc("PUT /api/datasets/{id}", requireAuth(handleUpdateDataset))
	mux.HandleFunc("DELETE /api/datasets/{id}", requireAuth(handleDeleteDataset))

	mux.HandleFunc("GET /api/datasets/{id}/logs", requireAuth(handleQueryLogs))

	log.Println("web api listening on :8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
