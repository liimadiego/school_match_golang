# SchoolMatch

Microservice in Go that allows users to register schools and submit student reviews. The application provides a REST API with JWT authentication.

## Features

- REST API built with Gin framework
- PostgreSQL database integration with GORM
- User authentication with JWT
- CRUD operations for schools and reviews
- Docker

## Getting Started

### Running with Docker

```bash
git clone https://github.com/liimadiego/school_match_golang
cd school_match_golang

docker-compose up -d

# Run the following commands inside the container or adjust the .env file to point to the correct host
# When using Docker, DB_HOST=postgres
# When running locally, DB_HOST=localhost
go run cmd/migrate/main.go -up
```

The API will be available at http://localhost:8080.