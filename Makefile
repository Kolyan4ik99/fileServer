
all: swag
	@go run cmd/main.go

swag:
	@swag init -g internal/app/fileServer.go

.PHONY: swag