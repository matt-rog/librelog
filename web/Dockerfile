FROM node:22-alpine AS frontend
WORKDIR /app
COPY frontend/package*.json ./
RUN npm ci
COPY frontend/ ./
RUN npm run build

FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
COPY --from=frontend /app/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -o /web .

FROM gcr.io/distroless/static-debian12
COPY --from=build /web /web
CMD ["/web"]
