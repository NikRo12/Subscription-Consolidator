# ─── Этап 1: сборка ───────────────────────────────────────────
FROM golang:1.25-alpine AS builder

RUN apk add --no-cache git
ENV GOTOOLCHAIN=auto

WORKDIR /build

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN go mod tidy

RUN GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux go build -o /server ./cmd/server
RUN GOTOOLCHAIN=auto CGO_ENABLED=0 GOOS=linux go build -o /parser ./cmd/parser

FROM alpine:latest AS server

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /server ./server

COPY --from=builder /build/certs ./certs

EXPOSE 8080
CMD ["./server"]

FROM alpine:latest AS parser

RUN apk add --no-cache ca-certificates

WORKDIR /app

COPY --from=builder /parser ./parser
COPY --from=builder /build/certs ./certs

CMD ["./parser"]