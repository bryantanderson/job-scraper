
dev:
	make build && make start

gin:
	docker compose -f docker-compose.yml --profile go build && \
	docker compose -f docker-compose.yml --profile go up

fastapi:
	docker compose -f docker-compose.yml --profile python build && \
	docker compose -f docker-compose.yml --profile python up

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
