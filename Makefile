include .envrc

#==================================================================================== #
# HELPERS
#==================================================================================== #

## help: print this help message
help:
	@echo 'Usage.'
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' | sed -e 's/^/ /'

#==================================================================================== #
# DEVELOPMENT
#==================================================================================== #

## postgres: run docker container containing postgres image
postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=P@ss0wrd -d postgres:12-alpine

## create/db: create database on postgres container
create/db:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

## drop/db: drop database on postgres container
drop/db:
	docker exec -it postgres12 dropdb simple_bank

## migrate/up: apply database migration up
migrate/up:
	migrate -path db/migration -database "postgresql://root:P@ssw0rd@localhost:5432/simple_bank?sslmode=disable" -verbose up

## migrate/down: apply database migration down
migrate/down:
	migrate -path db/migration -database "postgresql://root:P@ssw0rd@localhost:5432/simple_bank?sslmode=disable" -verbose down

## sqlc: generate sqlc file
sqlc:
	sqlc generate

#==================================================================================== #
# QUALITY CONTROL
#==================================================================================== #

## test: running unit test
test:
	go test -v -cover ./...

.PHONY: postgres create/db drop/db migrate/up migrate/down test
