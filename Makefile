DB_URL=postgresql://postgres:user@localhost:5432/simple_bank?sslmode=disable
PWD := ${CURDIR}

dbup:
	docker compose -f db/docker-compose.yml up -d

dbdown:
	docker compose -f db/docker-compose.yml down

postgres:
	docker run --name postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

postgresW:
	docker run --name postgres -p 5432:5432 --network bank-network -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:15-alpine

createdb:
	docker exec -it simple_bank createdb --username=root --owner=root simple_bank

dropdb:
	docker exec -it simple_bank dropdb simple_bank

new_migration:
	migrate create -ext sql -dir db/migration -seq $(name)

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
	go test -v -cover -short ./...

server:
	go run main.go

mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/tonisco/simple-bank-go/db/sqlc Store
	mockgen -destination worker/mock/distributor.go -package mockwk github.com/tonisco/simple-bank-go/worker TaskDistributor

proto:
	rm -f pb/*.go
	rm -f doc/swagger/*.swagger.json
	protoc --proto_path=proto --go_out=pb --go_opt=paths=source_relative \
    --go-grpc_out=pb --go-grpc_opt=paths=source_relative \
	--grpc-gateway_out=pb --grpc-gateway_opt paths=source_relative \
	--openapiv2_out=doc/swagger --openapiv2_opt=allow_merge=true,merge_file_name=simple_bank \
    proto/*.proto
	statik -src=./doc/swagger -dest=./doc

evans:
	evans --host localhost --port 8090 -r repl

redis:
	docker run --name redis -p 6379:6379 -d redis:7.2.1-alpine

.PHONY: dbup dbdown postgres postgresW createdb dropdb new_migration migrateup migrateup1 migratedown migratedown1 db_docs db_schema sqlc test server mock proto evans redis