package example

import (
	"embed"
	"github.com/yydspg/sustain"
)

//go:embed sql
var sqlFS embed.FS

//go:embed swagger/api.yaml
var swaggerContent string

func init() {
	sustain.AddModule(func(ctx interface{}) sustain.Module {
		return sustain.Module{
			Name: "example",
			SetupAPI: func() sustain.ApiRouter {
				return nil
			},
		}
	})
}
