startPostgresContainer:
	docker run --name BB-postgres --network bb-storage-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:14-alpine

createDB:
	docker exec -it BB-postgres createdb --username=root --owner=root BB-DB

dropDB:
	docker exec -it BB-postgres dropdb BB-DB

migrate_up:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/BB-DB?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/BB-DB?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createDB dropDB startPostgresContainer migrate_up migrate_down sqlc server