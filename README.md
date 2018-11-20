# Toggly

[![Build Status](https://travis-ci.org/Toggly/core.svg?branch=master)](https://travis-ci.org/Toggly/core)
[![Coverage Status](https://coveralls.io/repos/github/Toggly/core/badge.svg?branch=master)](https://coveralls.io/github/Toggly/core?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Toggly/core)](https://goreportcard.com/report/github.com/Toggly/core)
[![Version](https://img.shields.io/badge/version-1.0.0-brightgreen.svg)](https://github.com/Toggly/core/tree/1.0.0)
[![Docker](https://img.shields.io/badge/-Docker-blue.svg)](https://store.docker.com/community/images/toggly/toggly-server)

## Description

Toggle provides service for feature-flag/parameters based applications.
Toggly API allows organizing your project configuration as standalone flexible and reliable service.

See [GitHub repository](https://github.com/Toggly/core) for source code.

## Features

- Multiple projects
- Multiple environments for each project
- Multiple transports applicable ([REST API implemented](#rest-api))
- Different parameter types (bool, int, string)
- Flags/Properties inheritance
- MongoDB as a storage
- Cache layer plugins

## Usage

### Docker

See [Docker Store](https://store.docker.com/community/images/toggly/toggly-server) for details.

#### Docker Compose

See [`docker-compose.yml`](https://github.com/Toggly/core/blob/master/docker-compose.yml) file for details

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

To create your own plugin (for example for using Redis or Memcache) you have to implement [DataCache](https://github.com/Toggly/core/blob/master/internal/pkg/cache/cache.go) interface:

```go
type DataCache interface {
    Get(key string) ([]byte, error)
    Set(key string, data []byte) error
    Flush(scopes ...string)
}
```

Plugin package has to export `func GetCache() DataCache` function. See [in-memory-cache](https://github.com/Toggly/core/blob/master/internal/pkg/cache/in_memory.go) as a reference.

## REST API

See public [OpenAPI specification](https://app.swaggerhub.com/apis-docs/Toggly/toggly-server/1.0.0).

### Models

#### Project model

```json
{
    "owner": "owner1",                     //string
    "code": "project1",                    //string
    "description": "Project description",  //string
    "status": "active",                    //string `active` or `disabled`
    "reg_date": "2018-10-12T23:47:18.967Z" //string date in ISO 8601 format
}
```

#### Environment model

```json
{
    "owner": "owner1",                        //string
    "project_code": "project1",               //string
    "code": "dev",                            //string
    "description": "Development environment", //string
    "protected": false,                       //bool
    "reg_date": "2018-10-12T23:47:18.967Z"    //string date in ISO 8601 format
}
```

#### Object model

```json
{
    "code": "user",                        //string
    "owner": "owner1",                     //string
    "project_code": "project1",            //string
    "env_code": "env1",                    //string
    "description": "Object 1 description", //string
    "inherits": {
        "project_code": "project1",        //string
        "env_code": "env1",                //string
        "object_code": "obj0"              //string
    },
    "parameters": [
        {
            "code": "parameter1",          //string
            "description": "Parameter 1",  //string
            "type": "bool",                //string `bool`, `int`, `string`
            "value": false                 //depends on type
        }
    ]
}
```

### Request headers

Each request has to specify followed headers:

| Name                | Required | Description                                                               |
| ------------------- | -------- | ------------------------------------------------------------------------- |
| X-Toggly-Request-Id | No       | Request ID for request tracking. Automatically generated If not specified |
| X-Toggly-Owner-Id   | Yes      | Owner identifier                                                          |

### Project

#### `GET /v1/project` - projects list for owner

Response:

```json
[
    {
        "owner": "owner1",
        "code": "project1",
        "description": "Project description",
        "status": "active",
        "reg_date": "2018-10-12T23:47:18.967Z"
    }
]
```

#### `POST /v1/project` - create project

Request:

```json
{
    "code": "project1",
    "description": "Project description",
    "status": "active"
}
```

Response:

```json
{
    "owner": "owner1",
    "code": "project1",
    "description": "Project description",
    "status": "active",
    "reg_date": "2018-10-12T23:47:18.967Z"
}
```

#### `PUT /v1/project` - update project

Request:

```json
{
    "code": "project1",
    "description": "Project description",
    "status": "active"
}
```

Response:

```json
{
    "owner": "owner1",
    "code": "project1",
    "description": "Project description",
    "status": "active",
    "reg_date": "2018-10-12T23:47:18.967Z"
}
```

#### `GET /v1/project/{project_code}` - get project information

Response:

```json
{
    "owner": "owner1",
    "code": "project1",
    "description": "Project description",
    "status": "active",
    "reg_date": "2018-10-12T23:47:18.967Z"
}
```

#### `DELETE /v1/project/{project_code}` - delete project

### Environment

#### `GET /v1/project/{project_code}/env` - environments list for project

Response:

```json
{
  "environments": [
    {
      "owner": "owner1",
      "project_code": "project1",
      "code": "dev",
      "description": "Development environment",
      "protected": false,
      "reg_date": "2018-10-12T23:47:18.967Z"
    }
  ]
}
```

#### `POST /project/{project_code}/env` - create environment

Request:

```json
{
  "code": "dev",
  "description": "Development environment",
  "protected": false
}
```

Response:

```json
{
  "owner": "owner1",
  "project_code": "project1",
  "code": "dev",
  "description": "Development environment",
  "protected": false,
  "reg_date": "2018-10-12T23:47:18.967Z"
}
```

#### `PUT /project/{project_code}/env` - update environment

Request:

```json
{
  "code": "dev",
  "description": "Development environment",
  "protected": false
}
```

Response:

```json
{
  "owner": "owner1",
  "project_code": "project1",
  "code": "dev",
  "description": "Development environment",
  "protected": false,
  "reg_date": "2018-10-12T23:47:18.967Z"
}
```

#### `GET /project/{project_code}/env/{env_code}` - get environment information

Response:

```json
{
    "owner": "owner1",
    "project_code": "project1",
    "code": "dev",
    "description": "Development environment",
    "protected": false,
    "reg_date": "2018-10-12T23:47:18.967Z"
}
```

#### `DELETE /project/{project_code}/env/{env_code}` - delete environment

### Object

#### `GET /project/{project_code}/env/{env_code}/object` - get objects list

Response:

```json
[
    {
        "code": "user",
        "owner": "owner1",
        "project_code": "project1",
        "env_code": "env1",
        "description": "Object 1 description",
        "inherits": {
            "project_code": "project1",
            "env_code": "env1",
            "object_code": "obj0"
        },
        "parameters": [
            {
                "code": "parameter1",
                "description": "Parameter 1",
                "type": "bool",
                "value": false
            }
        ]
    }
]
```

#### `POST /project/{project_code}/env/{env_code}/object` - create object

Request:

```json
{
    "code": "obj1",
    "description": "Object 1 description",
    "inherits": {
        "project_code": "project1",
        "env_code": "env1",
        "object_code": "obj0"
    },
    "parameters": [
        {
            "code": "parameter1",
            "description": "Parameter 1",
            "type": "bool",
            "value": false
        }
    ]
}
```

Response:

```json

{
    "code": "user",
    "owner": "owner1",
    "project_code": "project1",
    "env_code": "env1",
    "description": "Object 1 description",
    "inherits": {
        "project_code": "project1",
        "env_code": "env1",
        "object_code": "obj0"
    },
    "parameters": [
        {
            "code": "parameter1",
            "description": "Parameter 1",
            "type": "bool",
            "value": false
        }
    ]
}

```

#### `PUT /project/{project_code}/env/{env_code}/object` - create object

Request:

```json
{
    "code": "obj1",
    "description": "Object 1 description",
    "inherits": {
        "project_code": "project1",
        "env_code": "env1",
        "object_code": "obj0"
    },
    "parameters": [
        {
            "code": "parameter1",
            "description": "Parameter 1",
            "type": "bool",
            "value": false
        }
    ]
}
```

Response:

```json
{
    "code": "user",
    "owner": "owner1",
    "project_code": "project1",
    "env_code": "env1",
    "description": "Object 1 description",
    "inherits": {
        "project_code": "project1",
        "env_code": "env1",
        "object_code": "obj0"
    },
    "parameters": [
        {
            "code": "parameter1",
            "description": "Parameter 1",
            "type": "bool",
            "value": false
        }
    ]
}
```

#### `GET /project/{project_code}/env/{env_code}/object/{obj_code}` - get object information

Response:

```json
{
    "code": "user",
    "owner": "owner1",
    "project_code": "project1",
    "env_code": "env1",
    "description": "Object 1 description",
    "inherits": {
        "project_code": "project1",
        "env_code": "env1",
        "object_code": "obj0"
    },
    "parameters": [
        {
            "code": "parameter1",
            "description": "Parameter 1",
            "type": "bool",
            "value": false
        }
    ]
}
```

#### `DELETE /project/{project_code}/env/{env_code}/object/{obj_code}` - delete object

#### `GET /project/{project_code}/env/{env_code}/object/{obj_code}/inheritors` - get all object inheritors as flat list

Response:

```json
[
    {
        "code": "user",
        "owner": "owner1",
        "project_code": "project1",
        "env_code": "env1",
        "description": "Object 1 description",
        "inherits": {
            "project_code": "project1",
            "env_code": "env1",
            "object_code": "obj0"
        },
        "parameters": [
            {
                "code": "parameter1",
                "description": "Parameter 1",
                "type": "bool",
                "value": false
            }
        ]
    }
]
```
