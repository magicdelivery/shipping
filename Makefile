cf = -f deploy/docker/compose.yaml
af = -f deploy/docker/compose-api-test.yaml
COVER_FILE ?= coverage.out

build: ## Build docker containers
	docker compose $(cf) build
up: ## Start docker containers
	docker compose $(cf) up -d
down: ## Stop docker containers
	docker compose $(cf) down
rebuild: ## Rebuild and start docker containers
	@make down
	@make build
	@make up
api-test: ## Build and start docker services and run API testing on them
	docker compose $(af) build
	docker compose $(af) -p mdtest up -d
	docker run --rm -v .\test\:/test --net md-ship-public ghcr.io/orange-opensource/hurl:latest --test --color --variables-file=/test/api/docker-vars /test/api/customer.hurl
	docker compose $(af) -p mdtest down

## Local development

generate: ## Generate go source files
	go generate ./...
tools: ## Install needed tools, e.g. linter
	@echo Installing tools from req-tools.txt
	@grep '@' deploy/local/req-tools.txt | xargs -tI % go install %
lint: tools ## Static check of the sources
	golangci-lint run --fix
help: ## Print this help
	@grep -E '^[a-zA-Z_-]+:.*## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Local testing

test:
	go test -count=1 ./...
shorttest: ## Run unit tests
	go test -count=1 -short ./... -coverprofile=$(COVER_FILE)
	go tool cover -func=$(COVER_FILE) | grep ^total
$(COVER_FILE):
	$(MAKE) shorttest
cover: $(COVER_FILE) ## Output coverage in human readable form in html
	go tool cover -html=$(COVER_FILE)
	rm -f $(COVER_FILE)
hurl: ## Run hurl API testing on localhost installation
	hurl --variables-file=.\test\api\local-vars .\test\api\customer.hurl

.PHONY: \
	api-test \
	build \
	cover \
	down \
	generate \
	help \
	hurl \
	lint \
	shorttest \
	test \
	tools \
	up \