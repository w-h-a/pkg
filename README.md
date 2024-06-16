# pkg

## Features

| Package             | Choices                  | Use Case                   |
| ------------------- | ------------------------ | -------------------------- |
| api                 | http                     | build gateway servers      |
| broker              | nats                     | asynchronous communication |
| client              | grpc                     | synchronous communication  |
| runtime             | kubernetes               | service info               |
| security/authn      | basic tokens, jwts       | authentication             |
| security/authz      | TBD                      | authorization              |
| security/encryption | TBD                      | encryption                 |
| server              | grpc                     | build backend servers      |
| store               | cockroach, redis, memory | data persistence           |
| telemetry           | memory                   | logs, metrics, and traces  |
