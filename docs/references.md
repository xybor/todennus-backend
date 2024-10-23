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

### Todennus (inhouse) scope

Refers [Resources](./resources.md#resources) to know the scope to read a
particular resource or field.

The general form of todennus scope is:

```
[todennus]action{:resource}
```

For example:

`[todennus]read`: read everything

`[todennus]read:user`: read user

`[todennus]*`: do anything

A scope can cover another scope. For example:

`[todennus]read` can cover `[todennus]read:user`.

`[todennus]read:user` can cover `[todennus]read:user.role`.

`[todennus]*:user` can cover `[todennus]read:user` and `[todennus]write:user`.


| Action   | Description                                                 |
| -------- | ----------------------------------------------------------- |
| `*`      | Grant read and write permission on a resource               |
| `read`   | Grant read-only access to resource or a field in a resource |
| `write`  | Grant create, update, and delete permission on a resource   |
| `create` | Grant ability to create a resource                          |
| `update` | Grant ability to update a resource                          |
| `delete` | Grant ability to delete a resource                          |

### Outside scope

We support use freely define your application scope. The scope is called outside if it violates the general form of todennus scope.

For example:

`read:user`

`read_user`

`write_something`

`[todennus]invalid`
