.PHONY: build run test docker-up docker-down logs clean deb

build:
	CGO_ENABLED=0 go build -tags netgo -o bin/mcp-is ./cmd/server

run:
	go run ./cmd/server

test:
	go test ./...

docker-up:
	docker compose up -d --build

docker-down:
	docker compose down

logs:
	docker compose logs -f

clean:
	rm -rf bin/
	rm -rf packaging/mcp-is_*
	docker compose down -v

deb:
	./packaging/build-deb.sh
