network:
	docker network create bank-network

postgres :
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine

createdb:
	docker exec -it postgres16 createdb -U root bank

dropdb:
	docker exec -it postgres16 dropdb -U root bank

migrateup :
	migrate -path pkg/db/migration -database  "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose up

migratedown :
	migrate -path pkg/db/migration -database  "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose down

migrateup1 :
	migrate -path pkg/db/migration -database  "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose up 1

migratedown1 :
	migrate -path pkg/db/migration -database  "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose down 1

sqlc: 
	sqlc generate

test:
	go test -v -cover ./...

server: 
	go run main.go -e development	

docker-build:
	docker build -t go-bank:latest .

docker-run:
	docker run --name go-bank --network bank-network -p 8080:8080 go-bank:latest

docker-stop:
	docker stop go-bank
	docker rm go-bank	

go-bank-up:
	docker network create bank-network 
	docker run --name postgres16 --network bank-network -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=password -d postgres:16-alpine
	sleep 3
	docker exec -it postgres16 createdb -U root bank
	migrate -path pkg/db/migration -database  "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose up
	docker build -t go-bank:latest .
	docker run --name go-bank --network bank-network -p 8080:8080 go-bank:latest


go-bank-down:
	docker stop go-bank
	docker rm go-bank
	migrate -path pkg/db/migration -database  "postgresql://root:password@localhost:5432/bank?sslmode=disable" -verbose down
	docker exec -it postgres16 dropdb -U root bank
	docker stop postgres16
	docker rm postgres16
	docker rmi go-bank:latest
	docker rmi postgres:16-alpine
	docker network rm bank-network


.PHONY: network postgres createdb dropdb migrateup migratedown sqlc test server migrateup1 migratedown1 docker-build docker-run docker-stop go-bank-up go-bank-down