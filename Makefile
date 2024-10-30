cf = -f deploy/docker/compose.yaml
af = -f deploy/docker/compose-api-test.yaml
COVER_FILE ?= coverage.out
REPORT_FILE ?= coverage.html

# Run in docker 

init: ## Initialize environement variables
	@if [ ! -f ./deploy/docker/.env ]; then \
		cp ./deploy/docker/.env.sample ./deploy/docker/.env; \
		echo "Adjust configuration in ./deploy/docker/.env"; \
	fi;
build: ## Build docker containers
	docker compose $(cf) build
up: init ## Start docker containers
	docker compose $(cf) up -d --remove-orphans
down: ## Stop docker containers
	docker compose $(cf) down
rebuild: ## Rebuild and start docker containers
	@make down
	@make build
	@make up
restart: ## Restart docker containers
	docker compose $(cf) restart

# Hurl API testing in docker

apitestbuild: ## Build containers for API testing
	docker compose $(af) build
apitestup: ## Start containers for API testing
	docker compose $(af) up -d --remove-orphans
apitestdown: ## Stop containers for API testing
	docker compose $(af) down
apitestrun: ## Run Hurl testing scripts in docker container and in mutual network
	docker run --rm -v ./test/:/test --net md-ship-public ghcr.io/orange-opensource/hurl:latest --test --color --variables-file=/test/api/docker-vars /test/api/customer.hurl
apitest: ## Build and start docker services and run API testing on them
	@make apitestbuild
	@make apitestup
	@make apitestrun
	@make apitestdown

## Local development

generate: ## Generate go source files
	go generate ./...
tools: ## Install needed tools, e.g. linter
	@echo Installing tools from req-tools.txt
	@grep '@' deploy/local/req-tools.txt | xargs -tI % go install %
lint: tools ## Static check of the sources
	golangci-lint run --fix
format: ## Format source code
	go fmt ./...
clean: ## Clean the project from built files
	rm -f $(COVER_FILE) 
	rm -f $(REPORT_FILE) 
	rm -f shipping_service 
	rm -f shipping_service.exe
help: ## Print this help
	@grep -E '^[a-zA-Z0-9_-]+:.*## .*$$' $(MAKEFILE_LIST) | sort | awk 'BEGIN {FS = ":.*## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

## Local testing

test: ## Run unit tests
	go test -count=1 ./...
test-with-coverage: ## Run short unit tests with coverage
	go test -count=1 -short ./... -coverprofile=$(COVER_FILE)
$(COVER_FILE):
	$(MAKE) test-with-coverage
show-total-coverage: $(COVER_FILE) ## Show percentage of coverage
	go tool cover -func=$(COVER_FILE) | grep ^total
check-coverage-threshold: $(COVER_FILE) ## Verify if coverage percentage is enough
	go tool cover -func $(COVER_FILE) \
	| grep "total:" | awk '{print ((int($$3) > 80) != 1) }'
report: $(COVER_FILE) ## HTML report for test coverage
	go tool cover -html=$(COVER_FILE) -o $(REPORT_FILE)
	rm -f $(COVER_FILE)
hurl: ## Run hurl API testing on localhost installation
	hurl --variables-file=.\test\api\local-vars .\test\api\customer.hurl

.PHONY: \
	apitest \
	apitestbuild \
	apitestdown \
	apitestrun \
	apitestup \
	build \
	check-coverage-threshold \
	clean \
	coverage \
	down \
	format \
	generate \
	help \
	hurl \
	lint \
	report \
	show-total-coverage \
	test \
	test-with-coverage \
	tools \
	up \