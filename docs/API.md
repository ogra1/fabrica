# Fabrica - API
Fabrica provides an open API that can be used to perform operations from an external
script or application. The API returns a JSON message with some common attributes:

- `code`: error code (empty string when the operation is successful)
- `message`: descriptive error message (empty string when the operation is successful)
- `records` (optional): a list of objects e.g. a list of repositories
- `record` (optional): a single object e.g. the details of a build

### Examples
> **List the repositories**
> ```
> curl http://localhost:8000/v1/repos
> ```
> ```
> {
>     "code": "",
>     "message": "",
>     "records": [
>         {
>             "id": "btqvgbp105lvsv3ubhd0",
>             "name": "fabrica",
>             "repo": "https://github.com/ogra1/fabrica",
>             "branch": "master",
>             "keyId": "",
>             "hash": "",
>             "created": "2020-10-01T15:39:27Z",
>             "modified": "2020-10-01T15:39:27Z"
>         },
>         {
>             "id": "btqvh19105lvsv3ubhdg",
>             "name": "logsync",
>             "repo": "https://github.com/slimjim777/logsync",
>             "branch": "master",
>             "keyId": "",
>             "hash": "",
>             "created": "2020-10-01T15:40:53Z",
>             "modified": "2020-10-01T15:40:53Z"
>         }
>     ]
> }
> ```

> **Create a repository**
> ```
> curl -X POST -d '{"repo":"https://github.com/ogra1/fabrica", "branch":"master", "keyId":""}' http://localhost:8000/v1/repos
> ```
> ```
> {
>     "code": "",
>     "message": "btqvgbp105lvsv3ubhd0"
> }
> ```

### List repositories
`GET /v1/repos`

Retrieve a list of the watched repositories.

**Request**

```
curl http://localhost:8000/v1/repos
```

**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | Error description        |
| records   | array  | List of repositories     |



### Create a repository
`POST /v1/repos`

Create a new watched repository.

**Request**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| repo      | string | URL of the repository    |
| branch    | string | Branch of the repository |
| keyId     | string | The ID of the ssh key    |


**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | ID of the repository     |



### Delete a repository
`POST /v1/repos/delete`

Delete a repository and, optionally, delete its builds.

**Request**

| Attribute    | Type   | Description                            |
| ---------    | ------ | ---------------------                  |
| id           | string | ID of the repository                   |
| deleteBuilds | bool   | `true` if the builds are to be deleted |


**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | Error description        |


### List builds
`GET /v1/builds`

Lists the build records.

**Request**
```
curl http://localhost:8000/v1/builds
```

**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | Error description        |
| records   | array  | List of builds           |


### Submit a build
`POST /v1/build`

Launches a new build for a repository.

**Request**

| Attribute    | Type   | Description                            |
| ---------    | ------ | ---------------------                  |
| repo         | object | the repository that us to be built     |

The repository object is similar to the records that are returned from the List 
Repository command.

E.g.
```
{
    "repo": "https://github.com/ogra1/fabrica",
    "branch": "master",
    "keyId": ""
}
```

**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | ID of the build          |




### Fetch an existing build
`GET /v1/builds/{buildId}`

Fetches the details of a build.

**Request**
```
curl http://localhost:8000/v1/builds/{buildId}
```
Where `buildId` is the ID of the build (as seen in the List builds command).

**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | Error description        |
| record    | object | The build record         |

The build record includes details of the build including:

- `status`: the status of the build e.g. `complete`.
- `logs`: an array of log objects, one for each line of the build log.
- `container`: the name of the LXD container.


### Delete a build
`DELETE /v1/builds/{buildId}`

Deletes a build and its generated assets.

**Request**
```
curl -X DELETE http://localhost:8000/v1/builds/{buildId}
```
Where `buildId` is the ID of the build (as seen in the List builds command).

**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | Error description        |


### List ssh keys
`GET /v1/keys`

**Request**
```
curl http://localhost:8000/v1/keys
```

**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | Error description        |
| records   | array  | List of ssh keys         |



### Register an ssh key
`POST /v1/keys`

Register an ssh key with the service.

**Request**

| Attribute    | Type   | Description                            |
| ---------    | ------ | ---------------------                  |
| name         | string | Descriptive name of the key            |
| data         | string | Base64-encoded ssh private key         |

**Response**

| Attribute | Type   | Description              |
| --------- | ------ | ---------------------    |
| code      | string | Error code               |
| message   | string | ID of the ssh-key record |
