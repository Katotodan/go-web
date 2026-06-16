# Go Web

A small Go web app demonstrating a MySQL-backed HTTP server with basic routing, static file serving, and user/table management endpoints.

## Features

- Gorillla mux router
- Static file serving from `static/`
- MySQL connection in `db/db.go`
- Create, drop, list tables
- Insert, query, list, and delete users
- Basic route examples and protected endpoints

## Prerequisites

- Go 1.26 or later
- MySQL server running and accessible at `127.0.0.1:3306`
- Database `book_store` available
- MySQL user: `root`
- MySQL password: `my-secret-pw`

> Update `db/db.go` if your database connection settings differ.

## Run locally

```bash
git clone https://github.com/Katotodan/go-web.git
cd go-web
go run main.go
```

The server listens on `http://localhost:3000`.

## Routes

- `GET /` - root route
- `GET /about` - serves `static/about.html`
- `GET /static/{path}` - serves static assets from `static/`
- `GET /books/{title}/page/{page}` - example route with path parameters

### Table management

- `POST /new/table` - create table
  - body: `{ "tableName": "users" }`
- `DELETE /delete/table/{tableName}` - drop table
- `GET /all/table` - list all tables

### User management

- `POST /insert/user` - insert user
  - body: `{ "username": "alice", "password": "secret", "table": "users" }`
- `GET /user/{table}/{id}` - fetch single user by id
- `GET /all/users/{table}` - list all users in a table
- `DELETE /delete/{table}/{id}` - delete user by id

### Example middleware routes

- `GET /foo`
- `GET /bar`
- `GET /secret`
- `GET /login`
- `GET /logout`

## Project structure

- `main.go` - application entrypoint and router setup
- `controller/controller.go` - HTTP handlers and request processing
- `db/db.go` - database connection logic
- `static/` - static HTML and CSS files

## Notes

- Passwords are hashed before storage using bcrypt.
- The app is configured for MySQL with `parseTime=true`.
- Adjust route handlers and database config for production readiness.
