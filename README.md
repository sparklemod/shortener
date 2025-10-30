# Shortener

A service for shortening links and collecting click-through statistics

### Features

- Clean Architecture
- Gracefully shutdown
- Unit tests
- PostgreSQL/MongoDB (choose in config)
- Goose Migrations

## Quick Start

- Start the Database

```
docker-compose up -d --build
```

- Start the Application

```
go run cmd/main.go
```


## API Testing
After starting the application, you can test API endpoints using the [http.http](http.http) file.