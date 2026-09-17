# RentGo API

RentGo is a backend REST API service for a vehicle rental platform. The
application provides functionality for user authentication, vehicle
management, and rental transactions.

The project is built using Go with a Clean Architecture approach to
maintain separation between business logic, data access, and API
handling.

---

## Features

- User authentication with JWT
- User registration and login
- Role-based access control
- Vehicle management
- Vehicle rental management
- Booking and rental transaction handling
- Request validation
- PostgreSQL database integration
- Swagger API documentation
- Unit testing with mocks

---

## Tech Stack

### Backend

- Go (Golang)
- Echo Framework
- GORM ORM
- PostgreSQL

### Authentication

- JWT Authentication
- bcrypt password hashing

### Documentation

- Swagger / OpenAPI

### Testing

- Go Testing
- Mock-based testing

---

## Architecture

RentGo follows a Clean Architecture pattern:

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
    │   └── swagger
    ├── file.env.example
    └── go.mod

### Layer Explanation

**Handler** - Receives HTTP requests - Parses request data - Returns API
responses

**Usecase** - Contains business logic - Handles authentication and
rental workflows

**Repository** - Handles database operations - Abstracts data
persistence

**Domain** - Contains entities, interfaces, and business rules

---

# Installation

## Clone Repository

```bash
git clone https://github.com/arjundwibastian/RentGo.git
cd RentGo
```

## Install Dependencies

```bash
go mod tidy
```

---

# Environment Configuration

Create:

    file.env

based on:

    file.env.example

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

---

# Database Setup

RentGo uses PostgreSQL.

Create database:

```sql
CREATE DATABASE rentgo;
```

Configure your database credentials inside:

    file.env

---

# Running the Application

Start the API:

```bash
go run main.go
```

The API runs on:

    http://localhost:3000

---

# API Documentation

Swagger documentation:

    http://localhost:3000/swagger/index.html

Swagger provides: - Available endpoints - Request parameters - Response
formats - API testing interface

---

# Authentication

RentGo uses JWT authentication.

Protected requests require:

```http
Authorization: Bearer <token>
```

---

# API Endpoints

## Authentication

### Register User

    POST /register

Create a new user account.

### Login

    POST /login

Authenticate user and receive JWT token.

---

## User

### Get User Profile

    GET /users/profile

Returns authenticated user information.

---

## Vehicle

### Create Vehicle

    POST /vehicles

Create a new rental vehicle.

### Get Vehicles

    GET /vehicles

Retrieve available vehicles.

### Update Vehicle

    PUT /vehicles/:id

Update vehicle information.

### Delete Vehicle

    DELETE /vehicles/:id

Remove a vehicle.

---

## Rental

### Create Rental

    POST /rentals

Create a rental transaction.

### Get Rental History

    GET /rentals

Retrieve rental records.

---

# Example Request

## Login

Request:

```json
{
  "email": "user@mail.com",
  "password": "password"
}
```

Response:

```json
{
  "message": "login success",
  "token": "jwt_token"
}
```

---

# Testing

Run tests:

```bash
go test ./...
```

Includes testing for: - Authentication usecase - User usecase -
Repository mocking

---

# Environment Security

Do not commit sensitive files:

    file.env
    .env

Use:

    file.env.example

as a template.

---

# License

This project is developed for educational and development purposes.
