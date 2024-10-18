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

test:
	go test ./...

run: config
	go run ${PACKAGE}

lint: install-lint
	${LINTBIN} run

precommit: format build test lint
	echo "OK"

bindir:
	mkdir -p ${BINDIR}

format: install-smartimports
	${SMARTIMPORTS} -exclude internal/mocks

install-lint: bindir
	test -f ${LINTBIN} || \
		(GOBIN=${BINDIR} go install github.com/golangci/golangci-lint/cmd/golangci-lint@${LINTVER} && \
		mv ${BINDIR}/golangci-lint ${LINTBIN})

LOCAL_BIN:=$(CURDIR)/bin

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

.PHONY: test_coverage
test_coverage:
	go test ./... -count=1 -race -cover -coverprofile cover.out
	go tool cover -func cover.out

.PHONY: test_coverage_html
test_coverage_html:
	go test ./... -count=1 -race -cover -coverprofile cover.out
	go tool cover -html=cover.out
