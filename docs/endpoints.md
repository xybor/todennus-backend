# API Endpoints

Almost all API endpoints require the [Authentication](./references.md#authentication).

API Endpoint which doesn't need Authentication will be marked as (*no-auth*).

## User

### Register (*no-auth*)

`POST /users`

Creates a new user. Returns a [user](./resources.md#user) object.

JSON Params:

| Field       | Type     | Description |
| ----------- | -------- | ----------- |
| `username`  | `string` | Username    |
| `passsword` | `string` | Password    |


### Get user

`GET /users/{user_id}`

`GET /users/username/{username}`

Returns a [user](./resources.md#user) object.


### Validate user
`POST /users/validate`

Returns a [user](./resources.md#user) object.

JSON Params:

| Field       | Type     | Description |
| ----------- | -------- | ----------- |
| `username`  | `string` | Username    |
| `passsword` | `string` | Password    |

## Clients

### Create first client (*no-auth*)

`POST /oauth2_clients/first`

Uses admin user to create the first oauth2 (confidential) client. Returns a
[client](./resources.md#resources) object and the `client_secret`.

Why this API? When todennus is started, there is no existed Client, we don't
have any flow to authenticate a user (all authentication flows require a
Client). This API is only valid if there is no existing Client and the user is
administrator.

JSON Params:

| Field      | Type     | Description    |
| ---------- | -------- | -------------- |
| `username` | `string` | Admin username |
| `password` | `string` | Admin password |
| `name`     | `string` | Client name    |

### Create client

`POST /oauth2_clients`

*Require `create:client` scope*.

Create a new oauth2 client. Returns a [client](./resources.md#resources)
object. If `is_confidential` is true, the result includes the `client_secret`
(`client_secret` can be retrieved by this API only).

JSON Params:

| Field             | Type     | Description                                     |
| ----------------- | -------- | ----------------------------------------------- |
| `name`            | `string` | Client name                                     |
| `is_confidential` | `bool`   | `true` if client can hold secret confidentially |

### Get client

`GET /oauth2_clients/{client_id}`

Returns a [client](./resources.md#oauth2-client) object.
