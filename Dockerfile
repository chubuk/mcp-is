FROM golang:1.26-alpine AS builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -tags netgo -ldflags="-s -w" -o mcp-is ./cmd/server

FROM alpine:3.22
RUN apk --no-cache upgrade && \
    apk --no-cache add ca-certificates tzdata docker-cli bash && \
    adduser -D -g '' -u 1001 -s /sbin/nologin appuser
USER appuser
WORKDIR /home/appuser/
COPY --from=builder /app/mcp-is .
EXPOSE 8080
ENTRYPOINT ["./mcp-is"]
