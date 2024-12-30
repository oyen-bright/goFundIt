.PHONY: up down stop restart logs test

up:
	docker-compose --env-file config/config.dev.yaml up

down:
	docker-compose down

stop:
	docker-compose stop

restart: down up

logs:
	docker-compose logs -f

test:
	docker-compose exec app go test -v ./...
