################################################################################
# Commands
################################################################################

.PHONY: clean
clean:
	-@docker compose down --remove-orphans -t 0
	-@docker compose rm -f 
	-@rm -rf .bin/

.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: run
run:
	-@docker compose up --force-recreate --build app

.PHONY: lint
lint: .bin/golangci-lint
	@.bin/golangci-lint run

################################################################################
# Tools
################################################################################

.bin/golangci-lint: $(wildcard vendor/github.com/golangci/*/*.go)
	@echo "building linter..."
	@cd vendor/github.com/golangci/golangci-lint/cmd/golangci-lint && go build -o $(shell git rev-parse --show-toplevel)/.bin/golangci-lint .
