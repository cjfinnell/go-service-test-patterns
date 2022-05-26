.PHONY: clean
clean:
	-@docker compose down --remove-orphans -t 0
	-@docker compose rm -f 

.PHONY: deps
deps:
	go mod tidy
	go mod vendor

.PHONY: run
run:
	-@docker compose up --force-recreate --build app
