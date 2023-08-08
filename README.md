
![GitHub Workflow Status (with event)](https://img.shields.io/github/actions/workflow/status/zardan4/petition-rest/linter.yml)
![GitHub go.mod Go version (subdirectory of monorepo)](https://img.shields.io/github/go-mod/go-version/zardan4/petition-rest)
[![Go Report Card](https://goreportcard.com/badge/github.com/zardan4/petition-rest)](https://goreportcard.com/report/github.com/zardan4/petition-rest)
![GitHub Repo stars](https://img.shields.io/github/stars/zardan4/petition-rest)
[![Go Reference](https://pkg.go.dev/badge/github.com/zardan4/petition-rest.svg)](https://pkg.go.dev/github.com/zardan4/petition-rest)

# REST for petitions
Handles requests for petitions, users and signatures. Can be used for writing small petitions interfaces
## Enpoints. Requests/Responses
<details>
<summary>
<h3>Auth</h3>
</summary>

#### POST /signup. Create new user
```go
CreateUser(user petitions.User) (int, error)
```
Request body
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
GetUserByName(name, password string) (petitions.User, error)
```
Request body
```json
{
    "name": "mark zuckerberg",
    "password": "secretpassword123"
}
```
Response
```json
{
    "token": usersJWT
}
```
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
- delete refresh sessions on signing in from the same fingerprint
- ~~log out endpoint~~