# Shortener

A service for shortening links and collecting click-through statistics

### Features

- Clean Architecture
- Gracefully shutdown
- Table-driven testing, integration tests
- PostgreSQL DB
- Goose Migrations

## Prerequisites

Docker and Docker Compose

Goose (migrations) - [installation guide](https://github.com/pressly/goose)

Fast command
```
go install github.com/pressly/goose/v3/cmd/goose@latest
```

## Quick Start

- Start the Database

```
docker-compose up -d --build
```

- Run Migrations

```
goose -dir=db/migrations postgres "postgres://user:password@localhost:54302/shortener?sslmode=disable" up
```

- Start the Application

```
 cmd/main.go
```


## API Testing
After starting the application, you can test API endpoints using the [http.http](http.http) file.