#!make
include .env
LINTER=golangci-lint

deploy: lint docker-compose.up migrate.up
	@echo "----- deploy -----"

DB_CONNECTION="host=$(DB_HOST) port=$(DB_PORT) user=$(POSTGRES_USER) password=$(POSTGRES_PASSWORD) dbname=$(POSTGRES_DB) sslmode=$(DB_SSL_MODE)"
MIGRATIONS_FOLDER="migrations"
TAG_NAME=dyleme/schedudler

.PHONY: docker-compose.up
docker-compose.up: 
	@echo "----- deploy by docker -----"
	@docker compose up -d


.PHONY: docker-compose.down
docker-compose.down:
	docker compose down

.PHONY: redeploy
redeploy:
	docker compose pull service
	docker compose down service
	docker compose up -d

.PHONY: migrate.up
migrate.up:
	@echo "----- running migrations up -----"
	@cd $(MIGRATIONS_FOLDER);\
	goose postgres ${DB_CONNECTION} up


.PHONY: migrate.down
migrate.down:
	@cd $(MIGRATIONS_FOLDER);\
	goose postgres ${DB_CONNECTION} down


.PHONY: migrate.create
migrate.create:
	@cd $(MIGRATIONS_FOLDER);\
	goose create $(name) sql

.PHONY: gen
gen: gen.go

.PHONY: gen.go
gen.go:
	@echo "----------- Generate go files ----------------"
	go generate ./...

.PHONY: lint
lint:
	@echo "----------- Lint project ----------------"
	golangci-lint run
	
	
.PHONY: lint-fix
lint-fix:
	@echo "----------- Lint project ----------------"
	golangci-lint run --fix	


.PHONY: test.build
test.build:
	@echo "----------- Test project build ----------------"
	@go build -o test_main ./cmd/timetable/
	@rm test_main

.PHONY: test
test: test.build
	@echo "----------- Test project ----------------"
	@go test ./...

.PHONY: docker.build
docker.build:
	docker build -t $(TAG_NAME) .

.PHONY: docker.push
docker.push: docker.build
	docker push $(TAG_NAME)
	
.PHONY: install
install: install.generators install.mocks install.linter install.tools
	
.PHONY: install.generators
install.generators:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@v1.26.0 
	go install github.com/oapi-codegen/oapi-codegen/v2/cmd/oapi-codegen@v2.3.0
	
.PHONY: install.tools
install.tools:
	go install github.com/pressly/goose/v3/cmd/goose@v3.21.1
	go install github.com/daixiang0/gci@v0.13.4

.PHONY: install.linter
install.linter:
	curl -sSfL https://raw.githubusercontent.com/golangci/golangci-lint/master/install.sh | sh -s -- -b $$(go env GOPATH)/bin v1.61.0
	
.PHONY: install.mocks
install.mocks:
	go install go.uber.org/mock/mockgen@v0.4.0


