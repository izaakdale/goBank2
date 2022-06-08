postgres:
	docker run --name postDb -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it postDb createdb --username=root --owner=root goBank
dropdb:
	docker exec -it postDb dropdb goBank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goBank?sslmode=disable" -verbose up
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goBank?sslmode=disable" -verbose down
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
_PHONY:
	postgres createdb dropdb migrateup migratedown sqlc test main