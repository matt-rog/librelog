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
    log_id TEXT,
    user_id UUID,
    name TEXT,
    description TEXT,
    data TEXT,
    PRIMARY KEY (log_id, user_id)
);

CREATE TABLE IF NOT EXISTS users (
    user_id UUID,
    account_number_hash TEXT,
    name TEXT,
    created_at TIMESTAMP,
    PRIMARY KEY (user_id)
);

CREATE TABLE IF NOT EXISTS tokens (
    token_hash TEXT,
    user_id UUID,
    PRIMARY KEY (token_hash)
);


CREATE INDEX IF NOT EXISTS log_set_index ON logs(data);