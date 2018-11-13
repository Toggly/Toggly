# Toggly

Toggle provides service for feature-flag/parameters based applications.
Toggly API allows organizing your project configuration as standalone flexible and reliable service.

## Features

- Multiple projects
- Multiple environment for each project
- Different parameter types (bool, int, string, enum)
- Flags/Properties inheritance

## API

See public [OpenAPI specification](https://app.swaggerhub.com/apis-docs/Toggly/Core/1.0.0)

## Usage

### Docker

```bash
docker run -it --rm toggly/toggly-server --help
```

```bash
Usage:
  toggly_server [OPTIONS]

toggly:
  -p, --port=            port (default: 8080) [$TOGGLY_API_PORT]
      --base-path=       Base API Path (default: /api) [$TOGGLY_API_BASE_PATH]
      --auth-token=      Consumer auth token [$TOGGLY_API_AUTH_TOKEN]
      --debug            Run in DEBUG mode

mongo:
      --store.mongo.url= mongo connection url [$TOGGLY_STORE_MONGO_URL]

Help Options:
  -h, --help             Show this help message
```

## Installation

```bash
cd cmd/toggly-server && go install
```

## Build Docker image

```bash
cd cmd/toggly-server && GOOS=linux go build -o ../../toggly-server && cd ../.. && docker build -t toggly/toggly-server .
```

## Development

To development run:

```bash
go run cmd/toggly-server/main.go
```
