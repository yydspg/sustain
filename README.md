##  sustain

an easy way to create http service based on [gin](https://github.com/gin-gonic/gin)
You can create it directly through sustain or develop it through secondary development
## usage
see [example](https://github.com/yydspg/sustain/tree/main/example)
```go
//go:embed sql
var sqlFS embed.FS

//go:embed swagger/api.yaml
var swaggerContent string

func init() {
	log.Logger.Log().Msg("initializing example service")
	sustain.AddModule(func(ctx interface{}) sustain.Module {
		return sustain.Module{
			Name:    "example",
			Swagger: swaggerContent,
			SQLDir:  sustain.NewSQLFS(sqlFS),
			SetupAPI: func() sustain.ApiRouter {
				return New(ctx.(*sustain.PeroModuleContext))
			},
			SetupTask: func() sustain.TaskRouter {
				return &Task{}
			},
			Start: func() error {
				log.Logger.Log().Msg("start example service")
				return nil
			},
			Stop: func() error {
				log.Logger.Log().Msg("stop example service")
				return nil
			},
		}
	})
}
```
## recommend
- project structure <br>
  example <br>
  ├── 1module.go  <br>
  ├── api.go  <br>
  ├── assets <br>
  │   └── ssl <br>
  │       ├── ssl.key <br>
  │       └── ssl.pem <br>
  ├── docs  <br>
  ├── service.go  <br>
  ├── sql  <br>
  │   └── test.sql  <br>
  ├── swagger  <br>
  │   └── api.yaml <br>
  └── task.go
## version
- 1.0 /24/09/30