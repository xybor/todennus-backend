# Todennus backend

An Identity, OpenID Connect, and OAuth2 Provider.

## Documentations

[API Refereneces](./docs/1.references.md)

[Resources](./docs/2.resources.md)

[API Endpoints](./docs/3.endpoints.md)

[OAuth2 Endpoints](./docs/4.oauth2_endpoints.md)

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

- Handle scope (**completed**).
- Allow integrate with custom external IdP.
- Allow integrate with third-party Identity/OAuth2 provider (Google, Discord, etc.).

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
$ make start-rest-server
```

5. The first registered user is always admininistrator.

```
POST /users

{
  "username": "admin",
  "password": "P@ssw0rd"
}
```

6. Create the first OAuth2 Client. This API Endpoint will be blocked after the
first client is created.

```
POST /oauth2_clients/first

{
  "name": "Admin Client",
  "is_confidential": true,
  "username": "{admin_username}",
  "password": "{admin_password}"
}
```

7. You can use the OAuth2 flow now.

```
POST /oauth2/token

grant_type=password&
client_id=CLIENT_ID&
client_secret=CLIENT_SECRET&
username=admin&
password=P@ssw0rd
```
