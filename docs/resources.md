# Resources

The scope column indicates which scope need to be included in the access token
to have a permission on that field.


## User

| Field          | Type        | Scope            | Description       |
| -------------- | ----------- | ---------------- | ----------------- |
| `id`           | `snowflake` |                  | User ID           |
| `username`     | `string`    |                  | Username (unique) |
| `display_name` | `string`    |                  | User display name |
| `role`         | `string`    | `read:user.role` | User role         |


## OAuth2 Client

| Field           | Type        | Scope                       | Description                          |
| --------------- | ----------- | --------------------------- | ------------------------------------ |
| `client_id`     | `snowflake` |                             | Client ID                            |
| `owner_id`      | `snowflake` |                             | User ID of the client's owner        |
| `name`          | `string`    |                             | Client name                          |
| `allowed_scope` | `string`    | `read:client.allowed_scope` | The maximum scope client can request |
