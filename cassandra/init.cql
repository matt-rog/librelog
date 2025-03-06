CREATE KEYSPACE IF NOT EXISTS librelog WITH replication = 
    {'class': 'SimpleStrategy', 'replication_factor': '1'};

USE librelog;

CREATE TABLE IF NOT EXISTS logs (
    user_id UUID,
    log_set TEXT,
    timestamp TIMESTAMP,
    data TEXT,
    PRIMARY KEY ((user_id, log_set), timestamp)
) WITH CLUSTERING ORDER BY (timestamp DESC);

CREATE INDEX IF NOT EXISTS log_set_index ON logs(data);

DESC TABLE logs;
