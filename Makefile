DB_URL=postgresql://postgres:user@localhost:5432/simple_bank?sslmode=disable
PWD := ${CURDIR}

dbup:
	docker compose -f db/docker-compose.yml up -d

dbdown:
	docker compose -f db/docker-compose.yml down

postgres:
	docker run --name postgres --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createdb:
	docker exec -it simple_bank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simple_bank dropdb simple_bank

migrateup:
	migrate -path db/migration -database "$(DB_URL)" -verbose up

migrateup1:
	migrate -path db/migration -database "$(DB_URL)" -verbose up 1

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

migratedown1:
	migrate -path db/migration -database "$(DB_URL)" -verbose down 1

db_docs:
	dbdocs build doc/db.dbml

db_schema:
	dbml2sql --postgres -o doc/schema.sql doc/db.dbml

sqlc:
	docker run --rm -v "$(PWD):/src" -w /src sqlc/sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/tonisco/simple-bank-go/db/sqlc Store

.PHONY: dbup dbdown postgres createdb dropdb migrateup migrateup1 migratedown migratedown1 db_docs db_schema sqlc test server mock