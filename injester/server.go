package main

import (
	"encoding/json"
	"flag"
	"log"
	"net/http"
	"os"
	"strings"
	"time"

	"github.com/gocql/gocql"
	"github.com/gorilla/websocket"
)

type LogObject struct {
    UserID string          `json:"user_id"`
    LogSet string          `json:"log_set"`
    Data   json.RawMessage `json:"data"`
}

var addr = flag.String("addr", "localhost:9000", "http service address")

var upgrader = websocket.Upgrader{}

var session *gocql.Session

func injest(w http.ResponseWriter, r *http.Request) {
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
			continue
		}
	
		var lo LogObject
		err = json.Unmarshal([]byte(message), &lo)
		if err != nil {
			log.Println("invalid json:", err)
			c.WriteMessage(mt, []byte("invalid json"))
			continue
		}

		currentTime := time.Now().UnixMilli()

		err = session.Query(`INSERT INTO logs (user_id, log_set, timestamp, data) VALUES (?,?,?,?)`, lo.UserID, lo.LogSet, currentTime, lo.Data).Exec();
		if err != nil {
        	log.Println("insert error:", err)
			c.WriteMessage(mt, []byte("insert error"))
			continue
    	}

		log.Printf("recv: %s", message)
		err = c.WriteMessage(mt, message)
		if err != nil {
			log.Println("write:", err)
			continue
		}
	}
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


	flag.Parse()
	log.SetFlags(0)
	http.HandleFunc("/injest", injest)
	log.Fatal(http.ListenAndServe(*addr, nil))
}
