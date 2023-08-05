build:
	docker-compose build

run:
	docker-compose up app

migrate:
	docker-compose run --rm migrate

swag:
	swag init -g cmd/main.go