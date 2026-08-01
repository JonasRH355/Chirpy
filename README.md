# Chirpy API

A production-ready REST API built in Go as part of the **Boot.dev - Learn HTTP Servers in Go** course.

This project demonstrates how to build an HTTP server from scratch using Go's standard library, following backend best practices such as authentication, routing, middleware, database integration, and webhook handling.

---

## ✅ Features

- [x] HTTP Server built with Go's standard library
- [x] RESTful API
- [x] PostgreSQL integration
- [x] User registration
- [x] JWT Authentication
- [x] Refresh Token flow
- [x] Authorization
- [x] Password hashing
- [x] CRUD operations
- [x] Middleware
- [x] Webhooks
- [x] Environment configuration

---

## 🚀 What this project does

Chirpy is a Twitter-like REST API where users can register, authenticate, and publish short posts ("chirps").

The application includes:

- User registration
- Secure login with JWT authentication
- Refresh token authentication flow
- User authorization
- CRUD operations for chirps
- PostgreSQL database integration
- Middleware for request metrics
- Webhook handling
- JSON request/response APIs
- Password hashing
- Content filtering
- Route parameter handling
- Environment variable configuration

---

## 📚 Topics I mastered

Building this project helped me gain practical experience with:

- HTTP Servers using Go
- REST API design
- HTTP Routing
- JSON encoding & decoding
- Middleware
- PostgreSQL integration
- SQL migrations
- Authentication
- Authorization
- JWT Tokens
- Refresh Tokens
- Password Hashing
- Webhooks
- Environment Variables
- HTTP Status Codes
- Request Validation
- Go's `net/http` package
- Project organization
- Production-ready backend architecture

---

## 💡 Why should someone care?

Most Go tutorials stop after building simple CRUD applications.

This project goes further by implementing many features commonly found in real production APIs, including:

- JWT authentication
- Refresh token flow
- Authorization
- Database persistence
- Middleware
- Webhooks
- Clean API architecture

It serves as a solid reference for developers learning backend engineering with Go while staying close to the standard library instead of relying heavily on external frameworks.

---

## ⚙️ How to install and run

### Clone the repository

```bash
git clone https://github.com/yourusername/chirpy.git
cd chirpy
```

### Install dependencies

```bash
go mod download
```

### Configure environment variables

Create a `.env` file:

```env
DB_URL=postgres://user:password@localhost:5432/chirpy?sslmode=disable
SECRET=your_jwt_secret
POLKAKEY=your_webhook_key
PLATFORM=dev
```

### Run database migrations

Run your PostgreSQL migrations before starting the server.

### Start the server

```bash
go run .
```

The API will be available at:

```
http://localhost:8080
```

---

## 🛠️ Tech Stack

- Go
- net/http
- PostgreSQL
- SQLC
- JWT Authentication
- UUID
- godotenv

---

## 📖 Course

This project was built while completing the **Learn HTTP Servers in Go** course by **Boot.dev**.

The course covers:

- HTTP Servers
- Routing
- REST Architecture
- JSON APIs
- PostgreSQL
- Authentication
- Authorization
- Webhooks
- Documentation

---

## 📜 Certificate

Completed **Learn HTTP Servers in Go** on Boot.dev.
