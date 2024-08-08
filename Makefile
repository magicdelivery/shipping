cf = -f deploy/docker/compose.yaml
af = -f deploy/docker/compose-api-test.yaml

build:
	docker compose $(cf) build
up:
	docker compose $(cf) up -d
down:
	docker compose $(cf) down
rebuild:
	@make down
	@make build
	@make up
generate:
	go generate ./...
test:
	go test -count=1 ./...
hurl:
	hurl --variables-file=.\test\api\local-vars .\test\api\customer.hurl
api-test:
	docker compose $(af) build
	docker compose $(af) -p mdtest up -d
	docker run --rm -v .\test\:/test --net md-ship-public ghcr.io/orange-opensource/hurl:latest --test --color --variables-file=/test/api/docker-vars /test/api/customer.hurl
	docker compose $(af) -p mdtest down

.PHONY: \
	hurl \
	build \
	up \
	down \
	generate \
	test \
	api-test