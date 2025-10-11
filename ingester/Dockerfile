FROM golang:1.24 AS build
WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY *.go ./
RUN CGO_ENABLED=0 go build -o /ingester .

FROM gcr.io/distroless/static-debian12
COPY --from=build /ingester /ingester
CMD ["/ingester"]
