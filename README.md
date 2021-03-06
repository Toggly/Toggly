# Toggly

[![Build Status](https://travis-ci.org/Toggly/core.svg?branch=master)](https://travis-ci.org/Toggly/core)
[![Coverage Status](https://coveralls.io/repos/github/Toggly/core/badge.svg?branch=master)](https://coveralls.io/github/Toggly/core?branch=master)
[![Go Report Card](https://goreportcard.com/badge/github.com/Toggly/core)](https://goreportcard.com/report/github.com/Toggly/core)
[![Version](https://img.shields.io/badge/version-1.0-brightgreen.svg)](https://github.com/Toggly/core/tree/1.0)
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
- In-memory cache
- Redis cache

## REST API Server

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
docker run -d -v mongo:/data --network toggly --name mongo mongo

# Start Redis server
docker run -d -v redis:/data --network toggly --name redis redis

# Start Toggly server
docker run -it -d --network toggly \
    -p 8080:8080 \
    --name toggly-server \
    toggly/toggly-server \
    --store.mongo.url=mongodb://mongo:27017/toggly \
    --cache.type=redis \
    --cache.redis.url=redis://redis:6379
```

### Parameters

| Short | Long                     | Environment       | Default | Description                  |
| ----- | ----------------- | ------------------------ | ------- | ---------------------------- |
| -v    | --version         |                          |         | Show version                 |
| -p    | --port            | `TOGGLY_API_PORT`        | `8080`  | Port                         |
|       | --base-path       | `TOGGLY_API_BASE_PATH`   | `/api`  | Base API Path                |
|       | --no-logo         |                          | `false` | Do not show application logo |
|       | --store.mongo.url | `TOGGLY_STORE_MONGO_URL` |         | Mongo connection url         |
|       | --cache.type      | `TOGGLY_CACHE_TYPE`      |         | Cache type [memory\|redis]   |
|       | --cache.redis.url | `TOGGLY_CACHE_REDIS_URL` |         | Redis connection url         |
| -h    | --help            |                          |         | Show help message            |

### Installation

```bash
cd cmd/toggly-server && go install
```

```bash
toggly-server --version
```

### Build Docker image

```bash
docker build -t toggly/toggly-server .
```

### REST API

See public [OpenAPI specification](https://app.swaggerhub.com/apis-docs/Toggly/toggly-server/1.0.0).

#### Models

##### Project model

```json
{
    "owner": "owner1",
    "code": "project1",
    "description": "Project description",
    "status": "active",
    "reg_date": "2018-10-12T23:47:18.967Z"
}
```

Where:

- `status` can be _active_ or _disabled_
- `reg_date` is date in ISO 8601 format

##### Environment model

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

Where:

- `reg_date` is date in ISO 8601 format

##### Object model

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

Where:

- `parameters.type` can be _boot_, _int_ or _string_
- `parameters.value` depends on type

#### Request headers

Each request has to specify followed headers:

| Name                | Required | Description                                                               |
| ------------------- | -------- | ------------------------------------------------------------------------- |
| X-Toggly-Request-Id | No       | Request ID for request tracking. Automatically generated If not specified |
| X-Toggly-Owner-Id   | Yes      | Owner identifier                                                          |

#### Project

##### `GET /v1/project` - projects list for owner

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

##### `POST /v1/project` - create project

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

##### `PUT /v1/project` - update project

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

##### `GET /v1/project/{project_code}` - get project information

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

##### `DELETE /v1/project/{project_code}` - delete project

#### Environment

##### `GET /v1/project/{project_code}/env` - environments list for project

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

##### `POST /project/{project_code}/env` - create environment

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

##### `PUT /project/{project_code}/env` - update environment

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

##### `GET /project/{project_code}/env/{env_code}` - get environment information

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

##### `DELETE /project/{project_code}/env/{env_code}` - delete environment

#### Object

##### `GET /project/{project_code}/env/{env_code}/object` - get objects list

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

##### `POST /project/{project_code}/env/{env_code}/object` - create object

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

##### `PUT /project/{project_code}/env/{env_code}/object` - create object

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

##### `GET /project/{project_code}/env/{env_code}/object/{obj_code}` - get object information

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

##### `DELETE /project/{project_code}/env/{env_code}/object/{obj_code}` - delete object

##### `GET /project/{project_code}/env/{env_code}/object/{obj_code}/inheritors` - get all object inheritors as flat list

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
