// AI-assisted code
package main

import (
	"time"

	"github.com/gocql/gocql"
)

type Dataset struct {
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

func dbCreateToken(session *gocql.Session, tokenHash string, userID gocql.UUID) error {
	return session.Query(
		`INSERT INTO tokens (token_hash, user_id, created_at) VALUES (?, ?, ?)`,
		tokenHash, userID, time.Now(),
	).Exec()
}

func dbGetUserByToken(session *gocql.Session, tokenHash string) (gocql.UUID, error) {
	var userID gocql.UUID
	err := session.Query(
		`SELECT user_id FROM tokens WHERE token_hash = ?`, tokenHash,
	).Scan(&userID)
	return userID, err
}

func dbDeleteToken(session *gocql.Session, tokenHash string) error {
	return session.Query(`DELETE FROM tokens WHERE token_hash = ?`, tokenHash).Exec()
}

func dbListDatasets(session *gocql.Session, userID gocql.UUID) ([]Dataset, error) {
	iter := session.Query(
		`SELECT log_id, name, description FROM logs_meta WHERE user_id = ?`, userID,
	).Iter()

	var datasets []Dataset
	var d Dataset
	for iter.Scan(&d.LogID, &d.Name, &d.Description) {
		d.UserID = userID.String()
		datasets = append(datasets, d)
	}
	if err := iter.Close(); err != nil {
		return nil, err
	}
	return datasets, nil
}

func dbGetDataset(session *gocql.Session, userID gocql.UUID, logID string) (Dataset, error) {
	var d Dataset
	err := session.Query(
		`SELECT log_id, name, description, data FROM logs_meta WHERE user_id = ? AND log_id = ?`,
		userID, logID,
	).Scan(&d.LogID, &d.Name, &d.Description, &d.Data)
	d.UserID = userID.String()
	return d, err
}

func dbCreateDataset(session *gocql.Session, userID gocql.UUID, logID, name, description string) error {
	return session.Query(
		`INSERT INTO logs_meta (user_id, log_id, name, description) VALUES (?, ?, ?, ?)`,
		userID, logID, name, description,
	).Exec()
}

func dbUpdateDataset(session *gocql.Session, userID gocql.UUID, logID, name, description string) error {
	return session.Query(
		`UPDATE logs_meta SET name = ?, description = ? WHERE user_id = ? AND log_id = ?`,
		name, description, userID, logID,
	).Exec()
}

func dbDeleteDataset(session *gocql.Session, userID gocql.UUID, logID string) error {
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
