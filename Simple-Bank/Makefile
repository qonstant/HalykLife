postgres:
	docker run --name postgres12 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:12-alpine

createdb:
	docker exec -it postgres12 createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it postgres12 dropdb simple_bank

migrateup:
	 migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up

migrateup1:
	 migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose up 1

migratedown:
	 migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down

migratedown1:
	 migrate -path db/migration -database "postgresql://root:password@localhost:5432/simple_bank?sslmode=disable" -verbose down 1

sqlc:
	sqlc generate

tests:
	cd api && go test -v -cover ./...

server:
	go run main.go
	
mock:
	mockgen -package mockdb -destination db/mock/store.go  Simple-Bank/db/sqlc Store

storetest:
	go test -timeout 30s -run ^TestTransferTx$ -v -cover ./api

coverfile:
	go test -coverprofile=c.out\
	go tool cover -html="c.out"

.PHONY: postgres createdb dropdb migrateup migratedown sqlc tests mock migrateup1 migratedown
