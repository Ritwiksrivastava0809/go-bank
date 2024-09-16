# Build Stage
FROM golang:1.23.1-alpine3.20 AS builder
WORKDIR /app
COPY . .
RUN go build -o main main.go

# Run Stage
FROM alpine:3.20
WORKDIR /app
COPY --from=builder /app/main .
COPY environment/development.yaml /app/environment/


EXPOSE 8080

CMD [ "/app/main" ]