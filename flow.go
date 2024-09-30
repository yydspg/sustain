package sustain

import (
	"github.com/gin-gonic/gin"
	"github.com/judwhite/go-svc"
	"github.com/rs/zerolog"
)

type PeroFlow interface {
	Prepare() *PeroModuleContext
	Run(p *PeroModuleContext)
}

func setup(ctx *PeroModuleContext) error {
	ms := GetModules(ctx)
	for _, m := range ms {
		if m.SetupAPI != nil {
			a := m.SetupAPI()
			if a != nil {
				ctx.log.Debug().Msg("sustain:" + m.Name + "register http route")
				a.Route(ctx.GetHttpRoute())
			}
		}
		if ctx.SetupTask {
			if m.SetupTask != nil {
				t := m.SetupTask()
				if t != nil {
					t.RegisterTasks()
				}
			}
		}
	}
	return nil
}
func Run(ctx *PeroModuleContext) {
	s := NewPeroServer(ctx)
	ctx.SetHttpRoute(s.p)
	err := setup(ctx)
	if err != nil {
		panic(err)
	}
	err = svc.Run(s)
	if err != nil {
		panic(err)
	}
}
func Prepare() *PeroModuleContext {

	gin.SetMode(gin.DebugMode)
	zerolog.SetGlobalLevel(zerolog.DebugLevel)
	var config PeroConfig
	config = &DefaultConfig{
		Addr:    ":9001",
		SSLAddr: "",
	}
	return NewPeroModuleContext(config, true)
}
