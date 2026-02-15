FROM golang:1.25-alpine AS builder

WORKDIR /app
COPY go.mod go.sum ./

RUN go mod download

COPY cmd ./cmd
COPY internal ./internal

RUN go build -o server ./cmd/server

FROM alpine:latest

RUN apk --no-cache add ca-certificates && update-ca-certificates

RUN addgroup -S appgroup && adduser -S appuser -G appgroup

WORKDIR /app

COPY --from=builder /app/server ./server

RUN chown -R appuser:appgroup /app

USER appuser

EXPOSE 9000

HEALTHCHECK --interval=30s --timeout=5s --start-period=10s --retries=3 \
	CMD wget -qO- http://localhost:9000/api/status/up >/dev/null 2>&1 || exit 1

CMD ["./server"]

