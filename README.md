<div align="center">

  <h1>GameClub</h1>

  <h4>gameclub-backend</h4>
  <h6>backend for gameclub application</h6>

[![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)
[![Go](https://img.shields.io/badge/Go-00ADD8.svg?style=for-the-badge&logo=go&logoColor=white)](https://go.dev)

</div>

## About

> [!NOTE]
> Planned

## Project Structure

```
gameclub-backend
├── config
├── db
├── docker
├── dtos
├── handlers
├── interfaces
├── mappers
├── middleware
├── mocks
├── models
├── observer
├── redis
├── repositories
├── routes
├── security
├── strategy
│   ├── christmas
│   ├── normal
│   └── summer
└── utils
```

## Getting Started

### 1. Export database variables

```bash
export DB_HOST=localhost
export DB_USER=secret_very_very
export DB_PASSWORD=secret_top_secret
export DB_NAME=gameclub
export DB_PORT=5432
```

### 2. Spin up PostgreSQL with Docker

```bash
docker compose -f docker/docker-compose.yml up -d
```

### 3. Run the application

```bash
go run main.go
```

### 4. Run tests and generate test report

```bash
# Run all tests
go test ./...

# Run tests with verbose output
go test -v ./...

# Run tests with coverage
go test -cover ./...

# Generate HTML coverage report
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html
```

## LICENSE

```
This is free and unencumbered software released into the public domain.

Anyone is free to copy, modify, publish, use, compile, sell, or
distribute this software, either in source code form or as a compiled
binary, for any purpose, commercial or non-commercial, and by any
means.

In jurisdictions that recognize copyright laws, the author or authors
of this software dedicate any and all copyright interest in the
software to the public domain. We make this dedication for the benefit
of the public at large and to the detriment of our heirs and
successors. We intend this dedication to be an overt act of
relinquishment in perpetuity of all present and future rights to this
software under copyright law.

THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT.
IN NO EVENT SHALL THE AUTHORS BE LIABLE FOR ANY CLAIM, DAMAGES OR
OTHER LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE,
ARISING FROM, OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR
OTHER DEALINGS IN THE SOFTWARE.

For more information, please refer to <https://unlicense.org>
```
