# go-request-middleware

Gin middleware utilities for request binding and validation.

## Installation

```bash
go get github.com/Serdiev/go-request-middleware
```

## Usage

### ValidatePath

Binds URL path parameters to a struct.

```go
import "github.com/Serdiev/go-request-middleware"

type UserQuery struct {
    Id   string
    Name string
}

type Service struct{}

func handler(c *gin.Context, query UserQuery, svc *Service) {
    c.JSON(200, gin.H{"id": query.Id, "name": query.Name})
}

r.GET("/users/:id", gin_utils.ValidatePath(handler, &Service{}))
```

### ValidateRequest

Binds JSON request body to a struct. Also injects `id` path parameter if the struct has an `Id` field.

```go
type CreateRequest struct {
    Id   string
    Name string
}

func handler(c *gin.Context, req CreateRequest, svc *Service) {
    c.JSON(200, gin.H{"id": req.Id, "name": req.Name})
}

r.POST("/users", gin_utils.ValidateRequest(handler, &Service{}))
```

### ValidateQuery

Binds query parameters to a struct. Also injects `id` path parameter if the struct has an `Id` field.

```go
type SearchRequest struct {
    Id    string
    Query string
}

func handler(c *gin.Context, req SearchRequest, svc *Service) {
    c.JSON(200, gin.H{"id": req.Id, "query": req.Query})
}

r.GET("/search", gin_utils.ValidateQuery(handler, &Service{}))
```
