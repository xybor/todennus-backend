# References

## Snowflake

Todennus utilizes Twitter's snowflake format for uniquely identifiable
descriptors (IDs).

Because Snowflake IDs are up to 64 bits in size (e.g. a int64), they are always
returned as strings in the HTTP API to prevent integer overflows in some
languages.

Snowflake ID example: `328184286924505088`.

## Authentication

Authentication is performed with the Authorization HTTP header as the following
format:

```
Authorization: Bearer {ACCESS_TOKEN}
```

## @me

`@me` can be used to replace `{user_id}` in all API requests requiring authentication to represent the `{user_id}` of the authorized user.

## Scope

Refers [Resources](./resources.md#resources) to know the scope to read a
particular resource or field.

Refers [API Endpoints](./endpoints.md#api-endpoints) to know the scope which
each API needs to execute.

A scope can cover another scope. For example:

`read` can cover `read:user`.

`read:user` can cover `read:user.role`.

`*:user` can cover `read:user` and `write:user`.


| Action   | Description                                                 |
| -------- | ----------------------------------------------------------- |
| `*`      | Grant read and write permission on a resource               |
| `read`   | Grant read-only access to resource or a field in a resource |
| `write`  | Grant create, update, and delete permission on a resource   |
| `create` | Grant ability to create a resource                          |
| `update` | Grant ability to update a resource                          |
| `delete` | Grant ability to delete a resource                          |
