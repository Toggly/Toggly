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
docker run -it -v mongo:/data --network toggly --name mongo -p 27017:27017 -d mongo
docker run -it --network toggly --name toggly-server -p 9090:8080 -d toggly/toggly-server --store.mongo.url=mongodb://mongo:27017/toggly
```

```bash
Usage:
  toggly-server [OPTIONS]

toggly:
  -p, --port=            port (default: 8080) [$TOGGLY_API_PORT]
      --base-path=       Base API Path (default: /api) [$TOGGLY_API_BASE_PATH]
      --debug            Run in DEBUG mode

mongo:
      --store.mongo.url= mongo connection url [$TOGGLY_STORE_MONGO_URL]

cache:
      --cache.disable    Disable cache [$TOGGLY_CACHE_DISABLE]
      --cache.in-memory  In-memory cache. Do not use for production. Only for development purposes. [$TOGGLY_CACHE_IN_MEMORY]

redis:
      --cache.redis.url= redis connection url [$TOGGLY_CACHE_REDIS_URL]

Help Options:
  -h, --help             Show this help message
```

## Installation

```bash
cd cmd/toggly-server && go install
```

## Build Docker image

```bash
version=$(git describe --always --tags) && \
revision=${version}-$(date +%Y%m%d-%H:%M:%S) && \
GOOS=linux go build -o toggly-server -ldflags "-X main.revision=${revision}" ./cmd/toggly-server && \
docker build -t toggly/toggly-server .
```

## Development

To development run:

```bash
go run cmd/toggly-server/main.go
```
