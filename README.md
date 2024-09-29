# sustain
an easy way to create http service,based on [gin](https://github.com/gin-gonic/gin)
## recommend 
project structure
> project
├── docs
├── internal
├── module
│   ├── example1
│   │   ├── api.go
│   │   ├── db.go
│   │   ├── service.go
│   │   ├── sql
│   │   └── swagger
│   └── example2
│       ├── api.go
│       ├── db.go
│       ├── service.go
│       ├── sql
│       └── swagger
└── pkg
> 
## how to use
- create 1moudle.go which aim to register http sevice 
```go
//go:embed sql
var sqlFS embed.FS

//go:embed swagger/api.yaml
var swaggerContent string
func init() {
	sustain.AddModule(func(ctx interface{}) register.Module {

		fmt.Println("register......")
		api := New(ctx.(*config.Context))
		return register.Module{
			Name: "group",
			SetupAPI: func() register.APIRouter {
				return api
			},
			SQLDir:  register.NewSQLFS(sqlFS),
			Swagger: swaggerContent,
        }
    }
}
```
- create api.go
```go

```