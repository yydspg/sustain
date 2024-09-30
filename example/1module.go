package example

import (
	"embed"
	"github.com/rs/zerolog/log"
	"github.com/yydspg/sustain"
)

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
