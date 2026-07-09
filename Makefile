.PHONY: build up down test test-backend test-frontend coverage coverage-backend coverage-frontend clean

IMAGE ?= calculator:latest

## Build the production, single-container Docker image.
build:
	docker build -t $(IMAGE) .

## Start the local dev stack (Vite dev server + Go server, hot-reloading).
up:
	docker compose up --build

## Stop and remove the local dev stack.
down:
	docker compose down

## Run backend and frontend test suites.
test: test-backend test-frontend

test-backend:
	cd apps/backend && go test ./...

test-frontend:
	cd apps/frontend && npm install --no-audit --no-fund && npm test

## Generate coverage reports for backend and frontend.
coverage: coverage-backend coverage-frontend

coverage-backend:
	cd apps/backend && go test ./... -coverprofile=coverage.out && go tool cover -func=coverage.out

coverage-frontend:
	cd apps/frontend && npm install --no-audit --no-fund && npm run coverage

## Remove build artifacts and containers/volumes created by `up`.
clean:
	docker compose down -v --remove-orphans
	rm -rf apps/frontend/node_modules apps/frontend/dist apps/frontend/coverage
	rm -f apps/backend/coverage.out
	rm -rf apps/backend/internal/web/dist/*
	echo '<!doctype html><title>calculator</title>' > apps/backend/internal/web/dist/index.html
