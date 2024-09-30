package example

import (
	"fmt"
	"github.com/rs/zerolog"
	"github.com/yydspg/sustain"
)

type Example struct {
	ctx     *sustain.PeroModuleContext
	log     zerolog.Logger
	service IService
}

func New(ctx *sustain.PeroModuleContext) *Example {
	return &Example{
		ctx:     ctx,
		log:     zerolog.Logger{},
		service: NewService(ctx),
	}
}
func (e *Example) Route(p *sustain.PeroHttp) {

	user := p.Group("v1/example/user", sustain.CORSMiddleware())
	user.Get("/do1", e.do1)
	user.Post("/do2", e.do2)

	// block format
	money := p.Group("v1/example/money", sustain.CORSMiddleware())
	{
		money.Get("/do3", e.do3)
		money.Get("/do4", e.do4)
	}
}
func (e *Example) do1(p *sustain.PeroContext) {
	fmt.Println("do1")
}
func (e *Example) do2(p *sustain.PeroContext) {
	fmt.Println("do2")
}
func (e *Example) do3(p *sustain.PeroContext) {
	fmt.Println("do3")
}
func (e *Example) do4(p *sustain.PeroContext) {
	fmt.Println("do4")
}
