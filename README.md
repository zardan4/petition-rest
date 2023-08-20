
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/zardan4/petition-rest/linter.yml)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/zardan4/petition-rest)
[![Go Report Card](https://goreportcard.com/badge/github.com/zardan4/petition-rest)](https://goreportcard.com/report/github.com/zardan4/petition-rest)
![GitHub Repo stars](https://img.shields.io/github/stars/zardan4/petition-rest)
[![Go Reference](https://pkg.go.dev/badge/github.com/zardan4/petition-rest.svg)](https://pkg.go.dev/github.com/zardan4/petition-rest)

# REST for petitions
Handles requests for petitions, users and signatures. Can be used for writing small petitions interfaces<br>
For this petitions branch use audit module: https://github.com/zardan4/petition-audit-rabbitmq
## Enpoints. Requests/Responses
<details>
<summary>
<h3>Auth</h3>
</summary>

#### POST /signup. Create new user
```go
CreateUser(user petitions.User) (int, error)
```
Request bodyad
```json
{
    "name": "mark zuckerberg",
    "grade": "3",
    "password": "secretpassword123"
}
```
Response
```json
{
    "id": newUserId
}
```

<br>

#### POST /signin. Sign in as old user
```go
GenerateTokens(name, password, fingerprint string) (core.JWTPair, error)
```
Request body
```json
{
    "name": "mark zuckerberg",
    "password": "secretpassword123",
    "fingerprint": "your_unique_device_fingerprint"
}
```
Response
```json
{
    "access_token": usersJWT,
    "refresh_token": userRefreshToken
}
```

<br>

#### POST /refresh. Refresh user's tokens by refresh token. Delete previous refresh token
```go
RefreshTokens(refreshToken, fingerprint string) (core.JWTPair, error)
```
Request body
```json
{
    "fingerprint": "your_unique_device_fingerprint"
}
```
Cookie
```json
{
    "refresh_token": "refresh_token_cookie"
}
```
Response
```json
{
    "access_token": newUsersJWT,
    "refresh_token": newUserRefreshToken
}
```

<br>

#### POST /logout. Delete user's refresh session
```go
Logout(refreshToken string) error
```
Cookie
```json
{
    "refresh_token": "refresh_token_cookie"
}
```
Response
```json
{
    "status": "ok"
}
```

### Additional
- Refresh session depends on fingerprint too so make unique refresh session from each user's device and don't use the same fingerprint(generate it [here](https://www.npmjs.com/package/fingerprint-generator))
- Follow [this scheme](https://www.figma.com/file/0KyFbPgCpoIK4BXovODFDl/Auth-JWT-scheme) to better understand how to use auth
</details>

<!-- petitions -->
<details>
<summary>
<h3>Petitions. Only auth</h3>
</summary>

#### GET /petitions. Get all petitions
```go
GetAllPetitions() ([]petitions.Petition, error)
```
Response
```json
{
    "data": [
        {
            "id": "id",
            "title": "title",
            "date": "date",
            "timeend": "timeend",
            "text": "text",
            "answer": "answer"
        }
    ]
}
```

<br>

#### POST /petitions. Create petition
```go
CreatePetition(title, text string, authorId int) (int, error)
```
Request body
```json
{
    "title": "title_example",
    "text": "text_example"
}
```
Response
```json
{
    "id": "id"
}
```

<br>

#### GET /petitions/{id}. Get petition by id
```go
GetPetition(petitionId int) (petitions.Petition, error)
```
Response
```json
{
    "id": "id",
    "title": "title",
    "date": "date",
    "timeend": "timeend",
    "text": "text",
    "answer": "answer"
}
```

<br>

#### PUT /petitions/{id}. Update petition by id
```go
UpdatePetition(petition petitions.UpdatePetitionInput, petitionId, userId int) error
```
Request body. Optional fields but at least one
```json
{
    "id": "id",
    "title": "title",
    "date": "date",
    "timeend": "timeend",
    "text": "text",
    "answer": "answer"
}
```
Response
```json
{
    "status": "ok"
}
```

<br>

#### DELETE /petitions/{id}. Delete petition by id
```go
DeletePetition(petitionId, userId int) error
```
Response
```json
{
    "status": "ok"
}
```

<br>

#### GET /petitions/{id}/signed. Get petition signed status by user
```go
CheckSignatorie(petitionId, userId int) (bool, error)
```
Response
```json
{
    "signed": bool
}
```
</details>

<!-- signatures -->
<details>
<summary>
<h3>Signatures. Only auth</h3>
</summary>

#### GET /petitions/{id}/subs. Get all signatures for petition
```go
GetAllSubs(petitionId int) ([]petitions.Sub, error)
```
Response
```json
{
    "data": [
        {
            "id": "id",
            "date": "date",
            "userId": "userId",
            "name": "username"
        }
    ]
}
```

<br>

#### POST /petitions/{id}/subs. Create signature for petition
```go
CreateSub(petitionId, userId int) (int, error)
```
Request body
```json
{}
```
Response
```json
{
    "id": "signatureId"
}
```

<br>

#### DELETE /petitions/{id}/subs. Delete signature for petition by user
```go
DeleteSub(subId, petitionId, userId int) error
```
Response
```json
{
    "status": "ok"
}
```
</details>

## Running:
Firstly, configure your .env
```makefile
make run # run containers
```
```makefile
make migrate # init tables
```
```makefile
make swag # init swagger
```
## TODO
- ~~docker-compose~~
- ~~unit tests~~
- ~~swagger~~
- ~~log out endpoint~~