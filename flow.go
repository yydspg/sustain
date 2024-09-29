package sustain

import (
	"github.com/gin-gonic/gin"
	"github.com/judwhite/go-svc"
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
	var config PeroConfig
	config = &DefaultConfig{
		Addr:    "127.0.0.1:80",
		SSLAddr: "127.0.0.1:443",
	}
	return NewPeroModuleContext(&config, true)
}
