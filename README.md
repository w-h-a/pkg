# pkg

## Features

| Package             | Choices            | Use Case                   |
| ------------------- | ------------------ | -------------------------- |
| api                 | http               | build gateway servers      |
| broker              | nats               | asynchronous communication |
| client              | grpc               | synchronous communication  |
| runtime             | kubernetes         | service info               |
| security/token      | basic tokens, jwts | token providers            |
| security/auth{n,z}  | ad hoc             | build ad hoc middlewares   |
| security/encryption | TBD                | encryption                 |
| server              | grpc               | build backend servers      |
| store               | cockroach, redis   | data persistence           |
| telemetry           | memory             | logs, metrics, and traces  |
