# Toggly

- [Toggly](#toggly)
  - [Description](#description)
  - [Features](#features)
  - [API](#api)
  - [Usage](#usage)
    - [Docker](#docker)
      - [Docker Compose](#docker-compose)
      - [Docker Run](#docker-run)
    - [Parameters](#parameters)
  - [Installation](#installation)
  - [Build Docker image](#build-docker-image)
  - [Development](#development)
    - [Plugins](#plugins)

## Description

Toggle provides service for feature-flag/parameters based applications.
Toggly API allows organizing your project configuration as standalone flexible and reliable service.

## Features

- Multiple projects
- Multiple environment for each project
- Different parameter types (bool, int, string, enum)
- Flags/Properties inheritance
- MongoDB as a storage
- Cache layer plugins

## API

See public [OpenAPI specification](https://app.swaggerhub.com/apis-docs/Toggly/Core/1.0.0).

## Usage

### Docker

See [Docker Store](https://store.docker.com/community/images/toggly/toggly-server) for details.

#### Docker Compose

See [`docker-compose.yml`](docker-compose.yml) file for details

```bash
docker-compose up -d
```

#### Docker Run

```bash
# Start MongoDB server
docker run -it -v mongo:/data --network toggly --name mongo -p 27017:27017 -d mongo

# Start Toggly server
docker run -it --network toggly --name toggly-server -p 8080:8080 -d toggly/toggly-server --store.mongo.url=mongodb://mongo:27017/toggly --cache-plugin=in-memory
```

### Parameters

| Short | Long              | Environment              | Default | Description                                                                                                                |
| ----- | ----------------- | ------------------------ | ------- | -------------------------------------------------------------------------------------------------------------------------- |
| -v    | --version         |                          |         | Version                                                                                                                    |
| -p    | --port            | `TOGGLY_API_PORT`        | `8080`  | Port                                                                                                                       |
|       | --base-path       | `TOGGLY_API_BASE_PATH`   | `/api`  | Base API Path                                                                                                              |
|       | --no-logo         |                          | `false` | Do not show application logo                                                                                               |
|       | --cache-plugin    | `TOGGLY_CACHE_PLUGIN`    |         | Cache plugin file. Skip `-cache.so` suffix. For example: `--cache-plugin=in-memory` will lookup `in-memory-cache.so` file. |
|       | --store.mongo.url | `TOGGLY_STORE_MONGO_URL` |         | Mongo connection url                                                                                                       |
| -h    | --help            |                          |         | Show help message                                                                                                          |

## Installation

```bash
cd cmd/toggly-server && go install
```

```bash
toggly-server --version
```

## Build Docker image

```bash
docker build -t toggly/toggly-server .
```

## Development

To development run:

```bash
go run cmd/toggly-server/main.go
```

### Plugins

Toggly supports plugins for caching layer. Plugins implementation base on native [Go plugin system](https://golang.org/pkg/plugin/)

By default `in-memory` plugin is available.

To use in-memory cache plugin `.so` file has to be compiled:

```bash
go build -buildmode=plugin -o in-memory-cache.so ./internal/plugin/in-memory-cache/cache.go
```

Than use `--cache-plugin` option to enable caching:

```bash
toggly-server --cache-plugin=in-memory
```

To create your own plugin (for example for using Redis or Memcache) you have to implement [DataCache](internal/pkg/cache/cache.go) interface:

```go
type DataCache interface {
	Get(key string) ([]byte, error)
	Set(key string, data []byte) error
	Flush(scopes ...string)
}
```

Plugin package has to export `func GetCache() DataCache` function. See [in-memory-cache](internal/plugin/in-memory-cache/cache.go) as a reference.
