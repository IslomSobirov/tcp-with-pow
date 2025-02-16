FROM golang:1.22.5 AS builder

WORKDIR /build

COPY . .

RUN go mod download

RUN GO111MODULE=on CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main ./cmd/client

FROM scratch

COPY --from=builder /build/main /
COPY --from=builder /build/config/config.json /config/config.json

ENTRYPOINT ["/main"]
