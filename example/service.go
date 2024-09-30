package example

import (
	"github.com/rs/zerolog"
	"github.com/yydspg/sustain"
)

type IService interface {
	do6()
}
type Service struct {
	ctx *sustain.PeroModuleContext
	log zerolog.Array
}

func (s *Service) do6() {

}
func NewService(ctx *sustain.PeroModuleContext) IService {
	return &Service{}
}
