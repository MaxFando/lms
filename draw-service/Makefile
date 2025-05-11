PROJECT := github.com/MaxFando/lms/draw-service
LOCAL_BIN := $(shell pwd)/bin

appName = structura
compose = docker-compose -f docker-compose-debug.yml -p $(appName)

DB_DSN = postgresql://postgres:postgres@localhost:5432/lms?sslmode=disable
DB_MIGRATION_DSN = postgresql://postgres:postgres@localhost:5432/lms?search_path=draw

PHONY: install-lint install-deps db-create-migration db-migrate db-rollback lint lint-fix test generate generate-proto lint-proto lock-proto

install-lint:
	GOBIN=$(LOCAL_BIN) go install github.com/golangci/golangci-lint/cmd/golangci-lint@latest

install-deps: install-lint
	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@v3.21.1
	GOBIN=$(LOCAL_BIN) go install github.com/bufbuild/buf/cmd/buf@v1.50.0
	GOBIN=$(LOCAL_BIN) go install github.com/pseudomuto/protoc-gen-doc/cmd/protoc-gen-doc@v1.5.1
	GOBIN=$(LOCAL_BIN) go install google.golang.org/protobuf/cmd/protoc-gen-go@v1.34.1
	GOBIN=$(LOCAL_BIN) go install -mod=mod google.golang.org/grpc/cmd/protoc-gen-go-grpc@v1.4.0
	GOBIN=$(LOCAL_BIN) go install github.com/grpc-ecosystem/grpc-gateway/v2/protoc-gen-grpc-gateway@v2.26.0

db-create-migration:
	@echo "Enter migration name:"
	@read name; \
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_MIGRATION_DSN) $(LOCAL_BIN)/goose -dir migrations create $$name sql

db-migrate:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_MIGRATION_DSN) $(LOCAL_BIN)/goose -dir migrations up

db-rollback:
	GOOSE_DRIVER=postgres GOOSE_DBSTRING=$(DB_MIGRATION_DSN) $(LOCAL_BIN)/goose -dir migrations down

lint: install-lint
	$(LOCAL_BIN)/golangci-lint run -c .golangci.yaml

lint-fix: install-lint
	$(LOCAL_BIN)/golangci-lint run -c .golangci.yaml --fix ./...

test:
	go test -v -race ./...

generate:
	@echo "Generating code..."
	go generate ./...
	@echo "Code generated!"

generate-proto: install-deps
	$(LOCAL_BIN)/buf generate --template ./api/grpc/buf.gen.yaml
