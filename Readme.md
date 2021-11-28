# Expense Tracker
Is meant to be a simple application that will help in tracking income and expense.

# Dependencies
1. Postgres
2. Redis
3. Go version 1.16.5

# API details.

## Auth endpoints (Un Authenticated)
- POST {{HOST}}/api/v1/register
- POST {{HOST}}/api/v1/login

## App endpoints (Authenticated)
- GET {{HOST}}/api/v1/token/refresh
- DELETE {{HOST}}/api/v1/logout
- GET {{HOST}}/api/v1/expense_types?page=1&per_page=10
- POST {{HOST}}/api/v1/expense_types

# Notes
- Upon logging in successfully we get back `auth_token` and `refresh_token`
- `auth_token` is valid for 15 minutes
- `refresh_token` is valid for 7 days.
- Both tokens are JWT in nature.
- When making a request to Authenticated endpoints pass header with Bearer - `{{auth_token}}`

# Application setup details. 
- Make sure you have the Dependencies installed.
- Create a .env file with the following in it.
```
DB_HOST=
DB_USER=
DB_NAME=
DB_PASSWORD=
REDIS_DSN=
ACCESS_SECRET=
REFRESH_SECRET=
```
- Start server with `go run main.go`

# Binaries
| OS | Latest release |
----------|-----------|
| MAC     | TODO      |
| LINUX   | TODO      |
| Windows | TODO      |


