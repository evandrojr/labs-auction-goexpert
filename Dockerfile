# Stage 1: Test
FROM golang:1.24 AS test

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

# Instala o MongoDB e inicia o serviço
RUN apt-get update && apt-get install -y mongodb
RUN mkdir -p /data/db
RUN mongod --fork --logpath /var/log/mongodb.log

# Executa os testes
RUN go test -v ./...

# Stage 2: Builder
FROM golang:1.24 AS builder

WORKDIR /app

COPY go.mod ./
COPY go.sum ./
RUN go mod download

COPY . .

RUN CGO_ENABLED=0 GOOS=linux go build -a -installsuffix cgo -o /app/auction cmd/auction/main.go

# Stage 3: Final
FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/auction .

EXPOSE 8080

ENTRYPOINT ["/app/auction"]
