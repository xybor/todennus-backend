# Todennus backend

An Identity, OpenID Connect, and OAuth2 Provider.

## Documentations

[API Refereneces](./docs/references.md)

[Resources](./docs/resources.md)

[Swagger Endpoint](#start-swagger-server)

## Tech stack

- Architecture: Clean architecture, Domain Driven Development.
- Database: [gorm](https://github.com/go-gorm/gorm), [go-migrate](https://github.com/golang-migrate/migrate), [postgreSQL](https://www.postgresql.org/), [redis](https://redis.io/).
- Mux: [go-chi](https://github.com/go-chi/chi).
- Docs: [swaggo](https://github.com/swaggo/swag).
- Container: [Docker](https://www.docker.com/).
- Deployment: [Docker compose](https://docs.docker.com/compose/).

## Target

### Architecture

Strictly follow Clean Architecture and DDD.

### Features

- OAuth2 Provider with:
  + Authorization Code Flow ***\*completed\****.
  + Authorization Code Flow With PKCE ***\*completed\****.
  + Implicit Flow.
  + Resource Owner Password Credentials Flow ***\*completed\****.
  + Client Credentials Flow.
  + Refresh Token Flow ***\*completed\****.
  + Device Flow (low priority).

- Support Open ID Connect.
- Allow integrate with external Identity/OAuth2 Provider ***\*completed\****.

### User traffic

- 100M users.
- 1M new users per day.
- 10M OAuth2 requests per day.

## Get started

### Start swagger server

```shell
$ make start-swagger
```

### Start System

Please refer [todennus-orchestration](https://github.com/xybor/todennus-orchestration) for starting our system.

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
