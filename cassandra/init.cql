CREATE KEYSPACE IF NOT EXISTS librelog WITH replication =
    {'class': 'SimpleStrategy', 'replication_factor': '1'};

USE librelog;

CREATE TABLE IF NOT EXISTS logs (
    user_id UUID,
    log_id TEXT,
    recv_time TIMESTAMP,
    data TEXT,
    PRIMARY KEY ((user_id, log_id), recv_time)
) WITH CLUSTERING ORDER BY (recv_time DESC);

CREATE TABLE IF NOT EXISTS logs_meta (
    user_id UUID,
    log_id TEXT,
    name TEXT,
    description TEXT,
    data TEXT,
    PRIMARY KEY ((user_id), log_id)
);

CREATE TABLE IF NOT EXISTS users (
    user_id UUID,
    account_number_hash TEXT,
    password_hash TEXT,
    name TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS users_by_account (
    account_number_hash TEXT,
    user_id UUID,
    PRIMARY KEY (account_number_hash)
);

CREATE TABLE IF NOT EXISTS tokens (
    token_hash TEXT,
    user_id UUID,
    name TEXT,
    prefix TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (token_hash)
) WITH default_time_to_live = 2592000;

CREATE TABLE IF NOT EXISTS tokens_by_user (
    user_id UUID,
    token_hash TEXT,
    name TEXT,
    prefix TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY ((user_id), token_hash)
);
