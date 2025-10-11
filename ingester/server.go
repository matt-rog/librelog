// AI-assisted code
package main

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
)

type LogObject struct {
	LogSet string          `json:"log_set"`
	Data   json.RawMessage `json:"data"`
}

var upgrader = websocket.Upgrader{}

var session *gocql.Session

func hashSHA256(s string) string {
	h := sha256.Sum256([]byte(s))
	return hex.EncodeToString(h[:])
}

func authenticateToken(token string) (gocql.UUID, error) {
	tokenHash := hashSHA256(token)
	var userID gocql.UUID
	err := session.Query(
		`SELECT user_id FROM tokens WHERE token_hash = ?`, tokenHash,
	).Scan(&userID)
	return userID, err
}

func ingestWS(w http.ResponseWriter, r *http.Request) {
	token := r.URL.Query().Get("token")
	if token == "" {
		http.Error(w, "missing token", http.StatusUnauthorized)
		return
	}
	userID, err := authenticateToken(token)
	if err != nil {
		http.Error(w, "invalid token", http.StatusUnauthorized)
		return
	}

	c, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		log.Print("upgrade:", err)
		return
	}
	defer c.Close()

	for {
		mt, message, err := c.ReadMessage()
		if err != nil {
			log.Println("read:", err)
			break
		}

		var lo LogObject
		if err := json.Unmarshal(message, &lo); err != nil {
			log.Println("invalid json:", err)
			c.WriteMessage(mt, []byte(`{"error":"invalid json"}`))
			continue
		}

		err = session.Query(
			`INSERT INTO logs (user_id, log_id, recv_time, data) VALUES (?,?,?,?)`,
			userID, lo.LogSet, time.Now(), lo.Data,
		).Exec()
		if err != nil {
			log.Println("insert error:", err)
			c.WriteMessage(mt, []byte(`{"error":"insert error"}`))
			continue
		}

		c.WriteMessage(mt, []byte(`{"status":"ok"}`))
	}
}

func ingestREST(w http.ResponseWriter, r *http.Request) {
	auth := r.Header.Get("Authorization")
	if !strings.HasPrefix(auth, "Bearer ") {
		http.Error(w, `{"error":"missing token"}`, http.StatusUnauthorized)
		return
	}
	token := strings.TrimPrefix(auth, "Bearer ")
	userID, err := authenticateToken(token)
	if err != nil {
		http.Error(w, `{"error":"invalid token"}`, http.StatusUnauthorized)
		return
	}

	var lo LogObject
	if err := json.NewDecoder(r.Body).Decode(&lo); err != nil {
		http.Error(w, `{"error":"invalid json"}`, http.StatusBadRequest)
		return
	}

	err = session.Query(
		`INSERT INTO logs (user_id, log_id, recv_time, data) VALUES (?,?,?,?)`,
		userID, lo.LogSet, time.Now(), lo.Data,
	).Exec()
	if err != nil {
		http.Error(w, `{"error":"insert error"}`, http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write([]byte(`{"status":"ok"}`))
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
	mux.HandleFunc("GET /ingest", ingestWS)
	mux.HandleFunc("POST /ingest", ingestREST)

	log.Println("ingester listening on :9000")
	log.Fatal(http.ListenAndServe(":9000", mux))
}
