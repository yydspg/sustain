package sustain

import (
	"embed"
	"github.com/rs/zerolog"
	"sync"
)

var once sync.Once
var moduleList []Module
var modules = make([]ModuleFunc, 0)

type ModuleFunc func(ctx interface{}) Module

type ApiRouter interface {
	Route(p *PeroHttp)
}
type TaskRouter interface {
	RegisterTasks()
}
type SQLFS struct {
	embed.FS
}
type Module struct {
	// module name
	Name string
	// api router
	SetupAPI func() ApiRouter
	// task router
	SetupTask func() TaskRouter
	// sql pwd
	SQLDir  *SQLFS
	Swagger string
	Service interface{}
	Start   func() error
	Stop    func() error
}
type PeroConfig interface {
	GetAddr() string
	GetSSLAddr() string
}
type DefaultConfig struct {
	Addr    string
	SSLAddr string
}

func (c *DefaultConfig) GetAddr() string {
	return c.Addr
}
func (c *DefaultConfig) GetSSLAddr() string {
	return c.SSLAddr
}

type PeroModuleContext struct {
	httpRouter *PeroHttp
	cfg        *PeroConfig
	log        *zerolog.Logger
	valueMap   sync.Map
	SetupTask  bool
}

func NewPeroModuleContext(cfg *PeroConfig, setupTask bool) *PeroModuleContext {
	return &PeroModuleContext{
		cfg:       cfg,
		log:       zerolog.DefaultContextLogger,
		valueMap:  sync.Map{},
		SetupTask: setupTask,
	}
}
func (p *PeroModuleContext) getConfig() *PeroConfig {
	return p.cfg
}
func (c *PeroModuleContext) GetHttpRoute() *PeroHttp {
	return c.httpRouter
}
func (c *PeroModuleContext) SetHttpRoute(p *PeroHttp) {
	c.httpRouter = p
}
func (c *PeroModuleContext) SetValue(value interface{}, key string) {
	c.valueMap.Store(key, value)
}

func (c *PeroModuleContext) Value(key string) any {
	v, _ := c.valueMap.Load(key)
	return v
}

// AddModule each module init(),call this method
func AddModule(moduleFunc func(ctx interface{}) Module) {
	modules = append(modules, moduleFunc)
}
func GetModules(ctx any) []Module {
	once.Do(func() {
		for _, m := range modules {
			moduleList = append(moduleList, m(ctx))
		}
	})

	return moduleList
}
func StartAllModule(ctx *PeroModuleContext) error {
	// 获取所有模块
	ms := GetModules(ctx)
	for _, m := range ms {
		if m.Start != nil {
			err := m.Start()
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func StopAllModule(ctx *PeroModuleContext) error {
	ms := GetModules(ctx)
	for _, m := range ms {
		if m.Stop != nil {
			err := m.Stop()
			if err != nil {
				return err
			}
		}

	}
	return nil
}
func GetModuleByName(name string, ctx any) Module {
	for _, m := range moduleList {
		if m.Name == name {
			return m
		}
	}
	return Module{}
}
func GetService(name string) interface{} {
	for _, m := range moduleList {
		if m.Name == name {
			return m.Service
		}
	}
	return nil
}
