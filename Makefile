
run-okx:
	go run -race cmd/okx/main.go

generate:
	go generate ./...

test:
	go test ./... -v

local-lint: ## run linter on host machine
	@echo "+ $@"
	golangci-lint run
.PHONY: local-lint

run-local-piston:
	go run local_run/piston_server/main.go

local-migrate-up:
	sql-migrate up -config=migrations/quant/dbconfig.yaml -env="local"

build-docker-dev: ## build docker image for development
 ifdef NOT_INSIDE_DEV_CONTAINER
	@echo "+ $@"
	docker build --tag torque:dev -f docker/dev.dockerfile .
	go mod download
 endif
.PHONY: build-docker-dev

ifndef INSIDE_DEV_CONTAINER
  NOT_INSIDE_DEV_CONTAINER = 1
endif

ifdef NOT_INSIDE_DEV_CONTAINER
  # $(1) - command to run inside container
  # $(2) - optional `-it` for interactive terminals
  RUN_IN_DEV_CONTAINER = docker run --rm $(2)                                  \
                         --name torque-dev                               \
                         -v "`pwd`":/app:delegated                             \
                         -v "`go env GOMODCACHE`":/go/pkg/mod:delegated        \
                         -v "`go env GOCACHE`":/root/.cache/go-build:delegated \
                         -e GOLANGCI_LINT_CACHE=/app/.cache/golangci-lint      \
                         -w /app                                               \
                         torque:dev                                      \
                         $(1)
else
  RUN_IN_DEV_CONTAINER = $(1)
endif

gen: build-docker-dev ## run code generation in docker
	@echo "+ $@"
	$(call RUN_IN_DEV_CONTAINER, make _gen)
.PHONY: gen

_gen:
	go generate ./...
.PHONY: _gen

lint: build-docker-dev ## run linter in docker
	@echo "+ $@"
	$(call RUN_IN_DEV_CONTAINER, golangci-lint run)
.PHONY: lint


testenv-start: ## start test environment in docker compose
	@echo "+ $@"
	./docker/testenv/env.sh start
.PHONY: testenv-start

testenv-restart: ## restart test environment in docker compose
	@echo "+ $@"
	./docker/testenv/env.sh restart
.PHONY: testenv-restart

testenv-stop: ## stop test environment in docker compose
	@echo "+ $@"
	./docker/testenv/env.sh stop
.PHONY: testenv-stop