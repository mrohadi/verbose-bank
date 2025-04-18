# include .envrc

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
	docker run --name postgres12 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=P@ss0wrd -d postgres:12-alpine

## create/db: create database on postgres container
create/db:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

## drop/db: drop database on postgres container
drop/db:
	docker exec -it postgres12 dropdb simple_bank

MIGRATE_VERSION=""
RESULT=""
ifneq ($(MIGRATE_VERSION),)
	RESULT := $(MIGRATE_VERSION)
endif
## migrate/up: apply database migration up
## 			MIGRATE_VERSION - parameter migrate version
migrate/up:
	migrate -path db/migration -database "postgresql://root:P@ssw0rd@localhost:5432/simple_bank?sslmode=disable" -verbose up

## migrate/down: apply database migration down
## 			MIGRATE_VERSION - parameter migrate version
migrate/down:
	migrate -path db/migration -database "postgresql://root:P@ssw0rd@localhost:5432/simple_bank?sslmode=disable" -verbose down
	
## migrate/up/version: apply database migration up
## 			MIGRATE_VERSION - parameter migrate version
migrate/up/version:
	migrate -path db/migration -database "postgresql://root:P@ssw0rd@localhost:5432/simple_bank?sslmode=disable" -verbose up $(MIGRATE_VERSION)

## migrate/down/version: apply database migration down
## 			MIGRATE_VERSION - parameter migrate version
migrate/down/version:
	migrate -path db/migration -database "postgresql://root:P@ssw0rd@localhost:5432/simple_bank?sslmode=disable" -verbose down $(MIGRATE_VERSION)

## sqlc: generate sqlc file
sqlc:
	@sqlc generate

## server/run: run the development server
server/run:
	go run main.go

#==================================================================================== #
# QUALITY CONTROL
#==================================================================================== #

FOLDER=cmd/api

## test: running specific unit test
test:
	go test -v -cover -run ^$(FUNC)$$ github.com/mrohadi/simplebank/$(FOLDER)

## test/all: running all unit test
test/all:
	go test -v -cover ./...
	
## test/all/profile: running all unit test with generate cover profile
test/all/profile:
	go test -coverprofile=/tmp/profile.out -v -cover ./...
	
## test/profile/show: show generated cover profile on web based
test/profile/show:
	go tool cover -html=/tmp/profile.out

## mock/gen: generate mock database code using mockgen tool
mock/gen:
	mockgen -package mockdb -destination db/mock/store.go github.com/mrohadi/simplebank/db/sqlc Store

.PHONY: postgres create/db drop/db migrate/up migrate/down test/all test/all/profile server/run mock/gen
