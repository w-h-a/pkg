# pkg

## Features

| Package   | Choices          | Use Case                   |
| --------- | ---------------- | -------------------------- |
| api       | http             | build gateway servers      |
| broker    | nats             | asynchronous communication |
| client    | grpc             | synchronous communication  |
| runtime   | kubernetes       | service info               |
| security  | TBD              | auth and encryption        |
| server    | grpc             | build backend servers      |
| store     | cockroach, redis | data persistence           |
| telemetry | memory           | logs, metrics, and traces  |
