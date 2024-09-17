# Build Stage
FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go
RUN apk add curl

RUN curl -L https://github.com/golang-migrate/migrate/releases/download/v4.14.1/migrate.linux-amd64.tar.gz | tar xvz
RUN mv migrate.linux-amd64 /app/pkg/db/migrate

# Run Stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY --from=builder /app/pkg/db/migrate ./pkg/db/migrate
COPY environment/development.yaml /app/environment/
COPY start.sh .
COPY wait-for.sh /app/wait-for.sh
COPY pkg/db/migration /app/migration   

EXPOSE 8080

CMD [ "/app/main" ]
ENTRYPOINT [ "/app/start.sh" ]
