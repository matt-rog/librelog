FROM cassandra:latest

COPY init.cql /docker-entrypoint-initdb.d/init.cql
CMD ["sh", "-c", "exec docker-entrypoint.sh cassandra -f"]
