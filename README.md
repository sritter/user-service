# user-service

A Go REST API for user CRUD operations with auto-generated Swagger (OpenAPI) documentation.

## Features
- RESTful CRUD for users
- Swagger docs at `/swagger/index.html`

## Quickstart
```sh
# Install swag CLI if not present
# go install github.com/swaggo/swag/cmd/swag@latest

# Generate Swagger docs
swag init

# Run the server
go run main.go
```

Open [http://localhost:8080/swagger/index.html](http://localhost:8080/swagger/index.html) to view the Swagger UI.


## Build and test with Docker:
```sh
docker build -t user-service:latest .  --no-cache  --progress=plain
```

Run the container:
```sh
docker run -d -p 8080:8080 --name user-service user-service:latest
``` 
