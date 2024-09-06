postgres :
	docker run --name postgres16 -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb -U root bank

dropdb:
	docker exec -it postgres16 dropdb -U root bank

migrateup :
	migrate -path pkg/db/migration -database  "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose up

migratedown :
	migrate -path pkg/db/migration -database  "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose down

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go -e development	
	
.PHONY: postgres createdb dropdb migrateup migratedown sqlc test server
