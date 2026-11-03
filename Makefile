.PHONY: build run test docker-up docker-down logs clean deb lint lint-fix lint-nginx lint-docker lint-frontend lint-all

GIXY ?= $(shell command -v gixy 2>/dev/null || echo $(HOME)/.local/bin/gixy)

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

lint:
	golangci-lint run ./...

lint-fix:
	golangci-lint run ./... --fix

# gixy resolves include paths absolutely (/etc/nginx/...), so stage a tmp tree
lint-nginx:
	docker compose build nginx
	bash web/nginx/entrypoint_test.sh
	docker compose run --rm --no-deps -e DOMAIN=nginx-test.local nginx nginx -t
	@tmp=$$(mktemp -d); \
	mkdir -p "$$tmp/snippets" "$$tmp/conf.d"; \
	sed "s|/etc/nginx/|$$tmp/|g" web/nginx/nginx.conf > "$$tmp/nginx.conf"; \
	cp web/nginx/snippets/security-headers.conf "$$tmp/snippets/"; \
	touch "$$tmp/mime.types" "$$tmp/conf.d/block_ai_bots.conf" "$$tmp/conf.d/block_search_bots.conf"; \
	$(GIXY) -l "$$tmp/nginx.conf" || true; \
	$(GIXY) -lll "$$tmp/nginx.conf"; ret=$$?; rm -rf "$$tmp"; exit $$ret

lint-docker:
	docker run --rm -v $(PWD):/work -w /work hadolint/hadolint hadolint --no-fail Dockerfile web/nginx/Dockerfile
	docker run --rm -v $(PWD):/work -w /work ghcr.io/aquasecurity/trivy:latest config . --severity HIGH,CRITICAL --quiet --exit-code 1 --skip-dirs packaging --ignorefile .trivyignore.yaml
	docker run --rm -v $(PWD):/work -w /work ghcr.io/aquasecurity/trivy:latest fs --scanners secret --quiet --exit-code 1 --skip-dirs packaging .
	docker compose config -q

lint-frontend:
	cd web && npm install --no-fund --no-audit && npx eslint static/js/

lint-all: lint lint-nginx lint-docker lint-frontend
