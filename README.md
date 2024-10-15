# Todennus backend

An Identity, OpenID Connect, and OAuth2 Provider.

## Documentations

[API Refereneces](./docs/references.md)

[Resources](./docs/resources.md)

[API Endpoints](./docs/endpoints.md)

[OAuth2 Endpoints](./docs/oauth2_endpoints.md)

## Tech stack

- Architecture: Clean architecture, Domain Driven Development.
- Database: [gorm](https://github.com/go-gorm/gorm), [go-migrate](https://github.com/golang-migrate/migrate), [postgreSQL](https://www.postgresql.org/), [redis](https://redis.io/).
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

- Allow integrate with external Identity/OAuth2 Provider.

### User traffic

- 100M users.
- 1M new users per day.
- 10M OAuth2 requests per day.

## Get started

You need to setup secret values at `.env` (or export environment variables).
Please refer the [.env.example](./.env.example).

###  Run without Docker

1. Install [Golang 1.23](https://go.dev/doc/install).

2. Install [Postgres](https://www.postgresql.org/download/).

3. Install [Redis](https://redis.io/docs/latest/operate/oss_and_stack/install/install-redis/).

4. Start the server.

```shell
$ make start-rest-server
```

### Run with Docker

1. Build dockerfile.

```shell
$ make docker-build
```

2. Start docker compose.

```shell
$ make docker-compose-up
```

### Create the first user and client

1. Create the first user. The first registered user is always admininistrator.

```
POST /users

{
  "username": "admin",
  "password": "P@ssw0rd"
}
```

2. Create the first OAuth2 Client. This API Endpoint will be blocked after the
first client is created.

```
POST /oauth2_clients/first

{
  "name": "Admin Client",
  "is_confidential": true,
  "username": "admin",
  "password": "P@ssw0rd"
}
```

3. You can use the OAuth2 flow now.

```
POST /oauth2/token

grant_type=password&
client_id=CLIENT_ID&
client_secret=CLIENT_SECRET&
username=admin&
password=P@ssw0rd
```
