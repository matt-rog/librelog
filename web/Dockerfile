FROM node:22-alpine AS frontend
WORKDIR /project/web/frontend
COPY web/frontend/package*.json ./
RUN npm ci
COPY web/frontend/ ./
COPY docs/ /project/docs/
RUN npm run build

FROM golang:1.24 AS build
WORKDIR /src
COPY web/go.mod web/go.sum ./
RUN go mod download
COPY web/*.go ./
COPY --from=frontend /project/web/frontend/dist ./frontend/dist
RUN CGO_ENABLED=0 go build -o /web .

FROM gcr.io/distroless/static-debian12
COPY --from=build /web /web
CMD ["/web"]
