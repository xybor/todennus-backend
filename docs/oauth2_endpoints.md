# OAuth2 Endpoints

## Resource Owner Password Flow

`POST /oauth2/token`


| Field           | Type     | Description                                  |
| --------------- | -------- | -------------------------------------------- |
| `grant_type`    | `string` | Must `password`                              |
| `client_id`     | `string` | Client ID                                    |
| `client_secret` | `string` | Client Secret                                |
| `username`      | `string` | User's username                              |
| `password`      | `string` | User's password                              |
| `scope`         | `string` | Client Secret [scope](./references.md#scope) |

For example:

```json
{
    "grant_type": "password",
    "client_id": "308994132968210433",
    "client_secret": "xnHjds...",
    "username": "admin",
    "password": "password",
    "scope": "read"
}
```


## Refresh Token Flow

`POST /oauth2/token`


| Field           | Type     | Description                           |
| --------------- | -------- | ------------------------------------- |
| `grant_type`    | `string` | Must `refresh_token`                  |
| `client_id`     | `string` | Client ID                             |
| `client_secret` | `string` | Client Secret (depend on client type) |
| `refresh_token` | `string` | Refresh token                         |

For example:

```json
{
    "grant_type": "password",
    "client_id": "308994132968210433",
    "client_secret": "xnHjds...",
    "refresh_token": "eyJhbGciOiJIUzI..."
}
```
