# shortener

/internal/
    /repository/
        /domain/     ← модели бизнес-сущностей (Link)
        /postgres/   ← реализация интерфейсов (SQL)
    /usecase/
        usecase.go   ← бизнес-логика (сценарии)
        contract.go  ← интерфейс Repository
        /dto/        ← DTO запросов/ответов
    /transport/      ← HTTP
