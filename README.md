# Todennus backend

A centralized Authentication Server and OAuth2 Provider.

## Tech stack

- Architecture: Clean architecture, Domain Driven Development.
- Database: [gorm](https://github.com/go-gorm/gorm), [go-migrate](https://github.com/golang-migrate/migrate), [postgreSQL](https://www.postgresql.org/).
- Mux: [go-chi](https://github.com/go-chi/chi).

## Target

### Architecture

Strictly follow Clean Architecture and DDD.

### Features

- OAuth2 Provider with:
 + Authorization Code Flow.
 + Authorization Code Flow With PKCE.
 + Implicit Flow.
 + Resource Owner Password Credentials Flow (**Completed**).
 + Client Credentials Flow.
 + Refresh Token Flow.
 + Device Flow (low priority).

- Support registering by OAuth2.

### User traffic

- 100M users.
- 1M new users per day.
- 10M OAuth2 requests per day.
