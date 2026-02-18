# REST API Contact Person

A RESTful API for managing contacts and groups built with Go.

## Database Setup

### Run Migrations

Apply all migrations to create database tables:
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/go_contact_person?sslmode=disable" up
```

### Rollback Migrations

Revert the last migration:
```bash
migrate -path migrations -database "postgres://postgres:postgres@localhost:5433/go_contact_person?sslmode=disable" down
```

### Create New Migration

Create a new migration file:
```bash
migrate create -ext sql -dir migrations -seq <migration_name>
```

Example:
```bash
migrate create -ext sql -dir migrations -seq create_users_table
```

## Running the Application

### Development Mode

Run the API server:
```bash
go run cmd/api/main.go
```

### Seed Database

Populate the database with sample data:
```bash
go run cmd/seeder/main.go
```

### Build for Production

Build the executable:
```bash
go build -o bin/api cmd/api/main.go
```

Then run the built executable:
```bash
./bin/api
```

## Available Commands

| Command | Description |
|---------|-------------|
| `go run cmd/api/main.go` | Run the API server |
| `go run cmd/seeder/main.go` | Seed the database with sample data |
| `migrate ... up` | Apply all pending migrations |
| `migrate ... down` | Rollback the last migration |
| `migrate create ...` | Create a new migration file |
| `go build -o bin/api cmd/api/main.go` | Build the application |

