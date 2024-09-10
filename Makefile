.PHONY: all

dev:
	make build && make start

.PHONY: api
api:
	docker compose -f docker-compose.yml --profile api --profile cache build && \
	docker compose -f docker-compose.yml --profile api --profile cache up

.PHONY: gin
gin:
	docker compose -f docker-compose.yml --profile go build && \
	docker compose -f docker-compose.yml --profile go up

instrument:
	docker compose -f docker-compose.yml --profile all --profile instrument build && \
	docker compose -f docker-compose.yml --profile all --profile instrument up

build:
	docker compose -f docker-compose.yml --profile all build && docker image prune -f

start:
	docker compose -f docker-compose.yml --profile all up

stop:
	docker compose -f docker-compose.yml --profile all down

pre-commit:
	pre-commit run --all-files
