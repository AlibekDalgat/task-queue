# Запуск
1. Скачать зависимости
```
go mod tidy
```
2. Запустить приложение (порт работы - 8000, указывается в файле .env)
```
go run cmd/main.go
```
# Тестирование
```
go test ./internal/delivery/
```
