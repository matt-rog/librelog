// AI-assisted code
package main

import (
	"time"

	"github.com/gocql/gocql"
)

type Logset struct {
	LogID       string `json:"log_id"`
	UserID      string `json:"user_id"`
	Name        string `json:"name"`
	Description string `json:"description"`
	Data        string `json:"data,omitempty"`
}

type LogEntry struct {
	RecvTime time.Time `json:"recv_time"`
	Data     string    `json:"data"`
}

type User struct {
	UserID             gocql.UUID
	AccountNumberHash  string
	PasswordHash       string
	Name               string
	CreatedAt          time.Time
}

func dbCreateUser(session *gocql.Session, userID gocql.UUID, accountHash, passwordHash, name string) error {
	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query(
		`INSERT INTO users (user_id, account_number_hash, password_hash, name, created_at) VALUES (?, ?, ?, ?, ?)`,
		userID, accountHash, passwordHash, name, time.Now(),
	)
	batch.Query(
		`INSERT INTO users_by_account (account_number_hash, user_id) VALUES (?, ?)`,
		accountHash, userID,
	)
	return session.ExecuteBatch(batch)
}

func dbGetUserIDByAccount(session *gocql.Session, accountHash string) (gocql.UUID, error) {
	var userID gocql.UUID
	err := session.Query(
		`SELECT user_id FROM users_by_account WHERE account_number_hash = ?`, accountHash,
	).Scan(&userID)
	return userID, err
}

func dbGetUser(session *gocql.Session, userID gocql.UUID) (User, error) {
	var u User
	err := session.Query(
		`SELECT user_id, account_number_hash, password_hash, name, created_at FROM users WHERE user_id = ?`, userID,
	).Scan(&u.UserID, &u.AccountNumberHash, &u.PasswordHash, &u.Name, &u.CreatedAt)
	return u, err
}

type APIKey struct {
	TokenHash string    `json:"token_hash"`
	Name      string    `json:"name"`
	Prefix    string    `json:"prefix"`
	CreatedAt time.Time `json:"created_at"`
}

func dbCreateToken(session *gocql.Session, tokenHash string, userID gocql.UUID, name, prefix string) error {
	now := time.Now()
	if name != "" {
		batch := session.NewBatch(gocql.LoggedBatch)
		batch.Query(
			`INSERT INTO tokens (token_hash, user_id, name, prefix, created_at) VALUES (?, ?, ?, ?, ?) USING TTL 0`,
			tokenHash, userID, name, prefix, now,
		)
		batch.Query(
			`INSERT INTO tokens_by_user (user_id, token_hash, name, prefix, created_at) VALUES (?, ?, ?, ?, ?)`,
			userID, tokenHash, name, prefix, now,
		)
		return session.ExecuteBatch(batch)
	}
	return session.Query(
		`INSERT INTO tokens (token_hash, user_id, created_at) VALUES (?, ?, ?)`,
		tokenHash, userID, now,
	).Exec()
}

func dbGetUserByToken(session *gocql.Session, tokenHash string) (gocql.UUID, error) {
	var userID gocql.UUID
	err := session.Query(
		`SELECT user_id FROM tokens WHERE token_hash = ?`, tokenHash,
	).Scan(&userID)
	return userID, err
}

func dbDeleteToken(session *gocql.Session, tokenHash string, userID gocql.UUID) error {
	batch := session.NewBatch(gocql.LoggedBatch)
	batch.Query(`DELETE FROM tokens WHERE token_hash = ?`, tokenHash)
	batch.Query(`DELETE FROM tokens_by_user WHERE user_id = ? AND token_hash = ?`, userID, tokenHash)
	return session.ExecuteBatch(batch)
}

func dbListTokens(session *gocql.Session, userID gocql.UUID) ([]APIKey, error) {
	iter := session.Query(
		`SELECT token_hash, name, prefix, created_at FROM tokens_by_user WHERE user_id = ?`, userID,
	).Iter()

	var keys []APIKey
	var k APIKey
	for iter.Scan(&k.TokenHash, &k.Name, &k.Prefix, &k.CreatedAt) {
		keys = append(keys, k)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return keys, nil
}

func dbListLogsets(session *gocql.Session, userID gocql.UUID) ([]Logset, error) {
	iter := session.Query(
		`SELECT log_id, name, description FROM logs_meta WHERE user_id = ?`, userID,
	).Iter()

	var logsets []Logset
	var d Logset
	for iter.Scan(&d.LogID, &d.Name, &d.Description) {
		d.UserID = userID.String()
		logsets = append(logsets, d)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return logsets, nil
}

func dbGetLogset(session *gocql.Session, userID gocql.UUID, logID string) (Logset, error) {
	var d Logset
	err := session.Query(
		`SELECT log_id, name, description, data FROM logs_meta WHERE user_id = ? AND log_id = ?`,
		userID, logID,
	).Scan(&d.LogID, &d.Name, &d.Description, &d.Data)
	d.UserID = userID.String()
	return d, err
}

func dbCreateLogset(session *gocql.Session, userID gocql.UUID, logID, name, description string) error {
	return session.Query(
		`INSERT INTO logs_meta (user_id, log_id, name, description) VALUES (?, ?, ?, ?)`,
		userID, logID, name, description,
	).Exec()
}

func dbUpdateLogset(session *gocql.Session, userID gocql.UUID, logID, name, description string) error {
	return session.Query(
		`UPDATE logs_meta SET name = ?, description = ? WHERE user_id = ? AND log_id = ?`,
		name, description, userID, logID,
	).Exec()
}

func dbDeleteLogset(session *gocql.Session, userID gocql.UUID, logID string) error {
	return session.Query(
		`DELETE FROM logs_meta WHERE user_id = ? AND log_id = ?`, userID, logID,
	).Exec()
}

func dbQueryLogs(session *gocql.Session, userID gocql.UUID, logID string, limit int, before, after *time.Time) ([]LogEntry, error) {
	query := `SELECT recv_time, data FROM logs WHERE user_id = ? AND log_id = ?`
	args := []interface{}{userID, logID}

	if before != nil && after != nil {
		query += ` AND recv_time < ? AND recv_time > ?`
		args = append(args, *before, *after)
	} else if before != nil {
		query += ` AND recv_time < ?`
		args = append(args, *before)
	} else if after != nil {
		query += ` AND recv_time > ?`
		args = append(args, *after)
	}

	query += ` LIMIT ?`
	args = append(args, limit)

	iter := session.Query(query, args...).Iter()

	var entries []LogEntry
	var e LogEntry
	for iter.Scan(&e.RecvTime, &e.Data) {
		entries = append(entries, e)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return entries, nil
}
