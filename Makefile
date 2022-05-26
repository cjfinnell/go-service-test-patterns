################################################################################
# Commands
################################################################################

.PHONY: clean
clean:
	-@docker compose down --remove-orphans -t 0
	-@docker compose rm -f 
	-@rm -rf .bin/ coverage.txt results.xml
	-@go clean -testcache

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

.PHONY: mocks
mocks: .bin/mockery
	@go generate .

.PHONY: test.unit
test.unit: .bin/gotestsum
	@.bin/gotestsum --format testname -- -cover -coverprofile coverage.txt -short

.PHONY: test.integration
test.integration: .bin/gotestsum
	@.bin/gotestsum --format testname --junitfile results.xml -- -cover -coverprofile coverage.txt

.PHONY: test
test:
	@docker compose up --force-recreate --exit-code-from int-test --build int-test
	@docker cp int-test:/test/coverage.txt .
	@docker cp int-test:/test/results.xml .

.PHONY: fuzz
fuzz:
	@docker compose up --force-recreate -d redis
	-go test -fuzz FuzzIntegration 
	@docker compose down

################################################################################
# Tools
################################################################################

.bin/golangci-lint: $(wildcard vendor/github.com/golangci/*/*.go)
	@echo "building linter..."
	@cd vendor/github.com/golangci/golangci-lint/cmd/golangci-lint && go build -o $(shell git rev-parse --show-toplevel)/.bin/golangci-lint .

.bin/mockery: $(wildcard vendor/github.com/vektra/mockery/*/*.go) redis.go
	@echo "building mock generator..."
	@cd vendor/github.com/vektra/mockery/v2 && go build -o $(shell git rev-parse --show-toplevel)/.bin/mockery .

.bin/gotestsum: $(wildcard vendor/gotest.tools/*/*.go)
	@echo "building test runner..."
	@cd vendor/gotest.tools/gotestsum && go build -o $(shell git rev-parse --show-toplevel)/.bin/gotestsum .
