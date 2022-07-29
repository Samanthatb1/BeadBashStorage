startPostgresContainer:
	docker run --name BB-postgres -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres:12-alpine

createDB:
	docker exec -it BB-postgres createdb --username=root --owner=root BB-DB

dropDB:
	docker exec -it BB-postgres dropdb BB-DB

migrate_up:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/BB-DB?sslmode=disable" -verbose up

migrate_down:
	migrate -path db/migrations -database "postgres://root:secret@localhost:5432/BB-DB?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -v -cover ./...

.PHONY: createDB dropDB startPostgresContainer migrate_up migrate_down sqlc