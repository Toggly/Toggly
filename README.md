# Toggly

to run:
```bash
go run app/*.go
```


## API

### Account

Data hierarchy

```
Account
 |
 - Role
 |
 - User
```

### Dictionary

Data hierarchy

```
Project
 |
 - Object
 | |
 | - Parameter 
 |
 - Environment
   | 
   - Object
     |
     - Parameter
```

#### Project

Get project

```
GET http://HOST/api/v1/dict/project/{project_id}
```

Response:
```
{
    "name": "project1",
    "environments": ["dev", "test", "prod"]
}
```

Save project

```
POST http://HOST/api/v1/dict/project

{
    "name": "project1",
    "environments": ["dev", "test", "prod"]
}
```

#### Object

Get object

```
GET http://HOST/api/v1/dict/project/{project_id}/object/{object_code}
```

Response:

```
{
    "code": "user",
    "parameters": []
    "overrides": Null
}
```

Save object

```
POST http://HOST/api/v1/dict/project/{project_id}/object

{
    "code": "user",
    "parameters": []
    "overrides": Null
}
```

#### Environment

Get environment

```
GET http://HOST/api/v1/dict/project/{project_id}/env/{env_code}
```

Response:

```
{
    "code": "dev"
    "description": "Description",
    "protected": true
}
```

Save environment

```
POST http://HOST/api/v1/dict/project/{project_id}/env

{
    "code": "dev",
    "description": "Description",
    "protected": true
}
```


#### Environment Object

Get object

```
GET http://HOST/api/v1/dict/project/{project_id}/env/{env_code}/object/{object_code}
```

Response:

```
{
    "code": "user",
    "parameters": []
}
```

Save object

```
POST http://HOST/api/v1/dict/project/{project_id}/env/{env_code}/object

{
    "code": "user",
    "parameters": []
}
```
