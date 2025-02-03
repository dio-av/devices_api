# Project devices_api

Devices API implements a CRUD operation for devices. It uses PostgreSQL as a persistent data, a REST API for interaction with the server.

## Project Structure and dependencies

- The project uses the Respository design Pattern for the PostgreSQL usage.
- Docker 
- A Mock for the database
- Chi Mux for the routing
- Go-swagger for the API documentation



## MakeFile

Run build make command with tests
```bash
make all
```

Build the application
```bash
make build
```

Run the application
```bash
make run
```
Create DB container
```bash
make docker-run
```

Shutdown DB Container
```bash
make docker-down
```

DB Integrations Test:
```bash
make itest
```

Live reload the application:
```bash
make watch
```

Run the test suite:
```bash
make test
```

Clean up binary from the last build:
```bash
make clean
```
