CURDIR=$(shell pwd)
BINDIR=${CURDIR}/bin
GOVER=$(shell go version | perl -nle '/(go\d\S+)/; print $$1;')
SMARTIMPORTS=${BINDIR}/smartimports_${GOVER}
LINTVER=v1.51.1
LINTBIN=${BINDIR}/lint_${GOVER}_${LINTVER}
PACKAGE=${CURDIR}/cmd/app


all: format build test lint

build: bindir
	go build -o ${BINDIR}/app ${PACKAGE}

run: config
	go run ${PACKAGE}

.PHONY: .bin-deps
.bin-deps:
	$(info Installing binary dependencies...)

	GOBIN=$(LOCAL_BIN) go install github.com/pressly/goose/v3/cmd/goose@latest

MIGRATE:=$(LOCAL_BIN)/goose
MIGRATE_DATABASE_DRIVER=postgres
MIGRATE_DATABASE_HOST=localhost
MIGRATE_DATABASE_PORT=7432
MIGRATE_DATABASE_URL="postgresql://postgres:passwd@$(MIGRATE_DATABASE_HOST):$(MIGRATE_DATABASE_PORT)/postgres?sslmode=disable"

.PHONY: migrate
migrate:
	$(MIGRATE) -dir migrations $(MIGRATE_DATABASE_DRIVER) $(MIGRATE_DATABASE_URL) up

.PHONY: reset_migrations
reset_migrations:
	$(MIGRATE) -dir migrations $(MIGRATE_DATABASE_DRIVER) $(MIGRATE_DATABASE_URL) reset

.PHONY: go_generate
go_generate:
	go generate ./...
