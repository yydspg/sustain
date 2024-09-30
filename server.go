package sustain

import (
	"github.com/judwhite/go-svc"
	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"net/http"
	"os"
	"path/filepath"
	"strings"
)

type PeroServer struct {
	p       *PeroHttp
	log     zerolog.Logger
	addr    string
	sslAddr string
	ctx     *PeroModuleContext
}

func NewPeroServer(ctx *PeroModuleContext) *PeroServer {
	r := NewPeroHttp()
	//r.Use(CORSMiddleware())
	return &PeroServer{
		p:       r,
		ctx:     ctx,
		addr:    ctx.cfg.GetAddr(),
		sslAddr: ctx.cfg.GetSSLAddr(),
	}
}
func (s *PeroServer) Init(env svc.Environment) error {
	if env.IsWindowsService() {
		dir := filepath.Dir(os.Args[0])
		return os.Chdir(dir)
	}
	return nil
}
func (s *PeroServer) run(sslAddr string, addr string) error {
	s.p.Any("v1/ping", func(p *PeroContext) {
		p.ResponseOK()
	})
	s.p.Any("/swagger/:module", func(c *PeroContext) {
		m := c.Param("module")
		module := GetModuleByName(m, s.ctx)
		if strings.TrimSpace(module.Swagger) == "" {
			c.Status(http.StatusNotFound)
			return
		}
		c.String(http.StatusOK, module.Swagger)
	})
	if len(addr) != 0 {
		if sslAddr != "" {
			go func() {
				err := s.p.Run(addr)
				if err != nil {
					panic(err)
				}
			}()
		} else {
			log.Logger.Log().Msg("sustain server run")
			err := s.p.Run(addr)
			if err != nil {
				return err
			}
		}
	}
	// https
	if sslAddr != "" {
		s.p.Use(TLSMiddleware(sslAddr))
		currDir, _ := os.Getwd()
		return s.p.RunTLS(sslAddr, currDir+"/assets/ssl/ssl.pem", currDir+"/assets/ssl/ssl.key")
	}
	return nil
}
func (s *PeroServer) Start() error {
	go func() {
		err := s.run(s.sslAddr, s.addr)
		if err != nil {
			panic(err)
		}
	}()

	err := StartAllModule(s.ctx)
	if err != nil {
		return err
	}
	return nil
}
func (s *PeroServer) Stop() error {

	return StopAllModule(s.ctx)
}
