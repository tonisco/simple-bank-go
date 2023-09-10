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

migratedown:
	migrate -path db/migration -database "$(DB_URL)" -verbose down

sqlc:
	docker run --rm -v "$(PWD):/src" -w /src sqlc/sqlc generate

test:
	go test -v -cover ./...

.PHONY: dbup dbdown postgres createdb dropdb migrateup migratedown sqlc test