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
  + Resource Owner Password Credentials Flow (**completed**).
  + Client Credentials Flow.
  + Refresh Token Flow (**completed**).
  + Device Flow (low priority).

- Support registering by third-party OAuth2 provider (Google, Facebook).

### User traffic

- 100M users.
- 1M new users per day.
- 10M OAuth2 requests per day.

## Usage

1. Install [Golang 1.23](https://go.dev/doc/install).

2. Install [Postgres](https://www.postgresql.org/download/).

3. Setup secret values at `config/.env` (or environment variables) and
   configurations at `config/default.ini`.

4. Start the server.

```shell
$ go run ./cmd/rest/main.go
```

5. Create the first user.

```
POST /users

{
  "username": "admin",
  "password": "P@ssw0rd"
}
```

6. Generate a temporary access token by the admin secret key (currently we don't have any OAuth2 Client, therefore we cannot generate access token by the normal flow).

```
POST /oauth2/token

Authorization: Admin $ADMIN_SECRET_KEY$

grant_type=password&
username=admin&
password=P@ssw0rd
```

7. Create the first OAuth2 Client. Note that you must save the `client_secret` in the response. The secret will never be retrieved by anyway.
```
POST /oauth2_clients

Authorization: Bearer $ACCESS_TOKEN$

{
  "name": "Admin Client",
  "is_confidential": true
}
```

8. Now you can use the normal OAuth2 now.

```
POST /oauth2/token

grant_type=password&
client_id=CLIENT_ID&
client_secret=CLIENT_SECRET&
username=admin&
password=P@ssw0rd
```
