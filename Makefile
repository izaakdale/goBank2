postgres:
	docker run --name postDb --network=banking-backend -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres
createdb:
	docker exec -it postDb createdb --username=root --owner=root goBank
dropdb:
	docker exec -it postDb dropdb goBank
migrateup:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goBank?sslmode=disable" -verbose up
migrateup1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goBank?sslmode=disable" -verbose up 1
migratedown:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goBank?sslmode=disable" -verbose down
migratedown1:
	migrate -path db/migration -database "postgresql://root:secret@localhost:5432/goBank?sslmode=disable" -verbose down 1
sqlc:
	sqlc generate
test:
	go test -v -cover ./...
server:
	go run main.go
mock:
	mockgen -destination db/mock/store.go -package mockdb github.com/izaakdale/goBank2/db/sqlc Store
dockerbuild:
	docker build -t gobank:latest .
dockerrun:
	docker run --name gobank --network=banking-backend -p 8080:8080 -e DB_SOURCE="postgresql://root:secret@postDb:5432/goBank?sslmode=disable" gobank:latest
_PHONY:
	postgres createdb dropdb migrateup migratedown sqlc test main migrateup1 migratedown1