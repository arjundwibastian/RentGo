# RentGo API

RentGo is a backend API for a vehicle rental application. This project
was built using Go and provides features such as user authentication,
vehicle management, and rental transactions.

The project uses Clean Architecture to separate the API layer, business
logic, and database operations.

## Tech Stack

- Go
- Echo Framework
- PostgreSQL
- GORM
- JWT Authentication
- Swagger

## Project Structure

    .
    ├── main.go
    ├── internal
    │   ├── config
    │   ├── domain
    │   ├── handler
    │   ├── middleware
    │   ├── repository
    │   ├── usecase
    │   └── validation
    ├── docs
    └── go.mod

## Features

- User registration and login
- JWT authentication
- User management
- Vehicle management
- Rental transaction handling
- Swagger API documentation

## Installation

Clone the repository:

```bash
git clone https://github.com/arjundwibastian/RentGo.git
```

Go to the project directory:

```bash
cd RentGo
```

Install dependencies:

```bash
go mod tidy
```

## Environment Setup

Create a `file.env` file based on `file.env.example`.

Example:

```env
APP_PORT=3000

DB_HOST=localhost
DB_PORT=5432
DB_USER=postgres
DB_PASSWORD=password
DB_NAME=rentgo

JWT_SECRET=your_secret_key
```

Make sure PostgreSQL is running and the database configuration matches
your environment.

## Running the Project

Start the API:

```bash
go run main.go
```

The API will run on:

    http://localhost:3000

## API Documentation

Swagger documentation is available at:

    http://localhost:3000/swagger/index.html

## Testing

Run tests using:

```bash
go test ./...
```

## Notes

This project is still under development and improvements may be added in
the future.
