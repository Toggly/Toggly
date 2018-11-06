# Toggly API

All responses will have the form

```json
{
    "data": "Mixed type"
}
```

## Project

### Projects list

**Definition**

`GET /api/v1/project`

**Response**

- `200 OK` on success

```json
[
    {
        "code": "my_project",
        "description": "Project description",
        "status": 0,
        "reg_date": "2018-07-08T21:42:31.236594677-07:00"
    },
    ...
]
```

### Get project

**Definition**

`GET http://HOST/api/v1/project/{project_code}`

**Response**

Example

```http
GET http://HOST/api/v1/project/my_project
```

Response:

```json
{
    "code": "my_project",
    "description": "Project description",
    "status": 0,
    "reg_date": "2018-07-08T21:42:31.236594677-07:00"
}
```

#### Save project

```http
POST http://HOST/api/v1/project
```

Example:

```http
POST http://HOST/api/v1/project

{
    "code": "my_project",
    "description": "",
    "status": 0,
    "reg_date": "2018-07-08T21:42:31.236594677-07:00"
}
```

Response

```json
{
    "status": "ok"
}
```

### Environment

#### Environments list

```http
GET http://HOST/api/v1/project/{project_id}/env
```

Example:

```http
GET http://HOST/api/v1/project/my_project/env
```

Response:

```json
{
    "environments": [
        {
            "code": "dev",
            "description": "Development environment",
            "protected": false
        },
        {
            "code": "prod",
            "description": "Production environment",
            "protected": true
        }
    ]
}
```

#### Get environment

```http
GET http://HOST/api/v1/project/{project_code}/env/{env_code}
```

Example:

```http
GET http://HOST/api/v1/project/my_project/env/prod
```

Response:

```json
{
    "code": "prod",
    "description": "Production environment",
    "protected": true
}
```

#### Save environment

```http
POST http://HOST/api/v1/dict/project/{project_code}/env
```

Example:

```http
POST http://HOST/api/v1/dict/project/my_project/env

{
    "code": "test",
    "description": "Test environment",
    "protected": false
}
```

Response

```json
{
    "status": "ok"
}
```

### Object

#### Get objects list

```http
GET http://HOST/api/v1/project/{project_code}/object
```

Example:

```http
GET http://HOST/api/v1/project/my_project/object
```

Response:

```json
{
    "objects": [
        {
            "code": "user",
            "inherits": "group",
            "description": ""
        }
    ]
}
```

#### Get object description

```http
GET http://HOST/api/v1/project/{project_code}/object/{object_code}
```

Example:

```http
GET http://HOST/api/v1/project/my_project/object/user
```

Response:

```json
{
    "code": "user",
    "inherits": "group",
    "description": "",
    "parameters": []
}
```

#### Save object description

```http
POST http://HOST/api/v1/dict/project/{project_id}/object
```

Example:

```http
POST http://HOST/api/v1/dict/project/my_project/object

{
    "code": "user",
    "inherits": "group",
    "description": "",
    "parameters": []
}
```

Response

```json
{
    "status": "ok"
}
```

#### Get object value

```http
GET http://HOST/api/v1/project/{project_code}/object/{object_code}?id=object_id&env=default
```

Example

```http
GET http://HOST/api/v1/project/my_project/object/user?id=1234&env=prod
```

Response:

```json
{
    "parameters": [
        {
            "code": "active",
            "value": false
        }
    ]
}
```

#### Save object value

```http
POST http://HOST/api/v1/project/{project_code}/object/{object_code}
```

Example:

```http
POST http://HOST/api/v1/project/my_project/object/user

{
    "id": 1234,
    "env": "prod",
    "parameters": [
        {
            "code": "active",
            "value": true
        }
    ]
}
```

### Parameter

#### Get parameters list

```http
GET http://HOST/api/v1/project/{project_code}/object/{object_code}/parameter
```

Example:

```http
GET http://HOST/api/v1/project/my_project/object/user/parameter
```

Response:

```json
{
    "parameters": [
        {
            "code": "active",
            "type": "bool",
            "value": true
        },
        {
            "code": "role",
            "type": "enum",
            "value": "user",
            "enum": "user,editor,admin"
        }
    ]
}
```

#### Get parameter

```http
GET http://HOST/api/v1/project/{project_code}/object/{object_code}/parameter/{param_code}
```

Example:

```http
GET http://HOST/api/v1/project/my_project/object/user/parameter/role
```

Response:

```json
{
    "code": "role",
    "type": "enum",
    "value": "user",
    "enum": "user,editor,admin"
}
```

#### Save parameter

```http
POST http://HOST/api/v1/project/{project_code}/object/{object_code}/parameter
```

Example:

```http
POST http://HOST/api/v1/project/my_project/object/user/parameter

{
    "code": "role",
    "type": "enum",
    "value": "user",
    "enum": "user,editor,admin"
}
```

Response

```json
{
    "status": "ok"
}
```

## Admin API : TBD