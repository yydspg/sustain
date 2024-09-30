package sustain

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog"
	"github.com/yydspg/sustain/cache"
	"net/http"
	"os"
	"strconv"
	"sync"
)

type PeroHttp struct {
	r    *gin.Engine
	pool sync.Pool
}
type PeroContext struct {
	*gin.Context
	log zerolog.Logger
}
type PeroHandlerFunc func(c *PeroContext)

type PeroRouterGroup struct {
	routers *gin.RouterGroup
	P       *PeroHttp
}

func NewPeroHttp() *PeroHttp {
	p := &PeroHttp{
		pool: sync.Pool{},
		r:    gin.New(),
	}
	p.r.Use(gin.Recovery())
	p.pool.New = func() interface{} { return allocatePeroContext() }
	return p
}
func allocatePeroContext() *PeroContext {
	return &PeroContext{Context: nil, log: zerolog.New(os.Stdout).With().Logger()}
}
func newPeroRouterGroup(router *gin.RouterGroup, p *PeroHttp) *PeroRouterGroup {
	return &PeroRouterGroup{
		routers: router,
		P:       p,
	}
}
func (p *PeroHttp) Group(path string, handlers ...PeroHandlerFunc) *PeroRouterGroup {
	return newPeroRouterGroup(p.r.Group(path, p.adaptor(handlers)...), p)
}
func (r *PeroRouterGroup) Post(path string, handlers PeroHandlerFunc) {
	r.routers.POST(path, r.P.convert(handlers))
}
func (r *PeroRouterGroup) Get(path string, handlers PeroHandlerFunc) {
	r.routers.GET(path, r.P.convert(handlers))
}
func (p *PeroHttp) Use(handlers ...PeroHandlerFunc) {
	p.r.Use(p.adaptor(handlers)...)
}
func (p *PeroHttp) Any(path string, handlers ...PeroHandlerFunc) {
	p.r.Any(path, p.adaptor(handlers)...)
}
func (p *PeroHttp) Run(addr string) error {
	return p.r.Run(addr)
}
func (p *PeroHttp) RunTLS(addr, certFile, keyFile string) error {
	return p.r.RunTLS(addr, certFile, keyFile)
}
func GetLoginUID(token string, tokenPrefix string, cache cache.Cache) string {
	uid, err := cache.Get(tokenPrefix + token)
	if err != nil {
		return ""
	}
	return uid
}

// adaptor transmit PeroHandlerFunc -> gin.HandlerFunc
func (p *PeroHttp) adaptor(handlers []PeroHandlerFunc) []gin.HandlerFunc {
	ginHandlers := make([]gin.HandlerFunc, 0, len(handlers))
	if handlers != nil {
		for _, v := range handlers {
			ginHandlers = append(ginHandlers, p.convert(v))
		}
	}
	return ginHandlers
}

// convert
func (p *PeroHttp) convert(handler PeroHandlerFunc) gin.HandlerFunc {
	return func(c *gin.Context) {
		hc := p.pool.Get().(*PeroContext)
		hc.reset()
		hc.Context = c
		defer p.pool.Put(hc)
		// process
		handler(hc)
	}
}

// Response Response
func (p *PeroContext) Response(data interface{}) {
	p.JSON(http.StatusOK, data)
}
func (p *PeroContext) ResponseOK() {
	p.JSON(http.StatusOK, gin.H{
		"status": http.StatusOK,
	})
}

func (p *PeroContext) ResponseWithStatus(status int, data interface{}) {
	p.JSON(status, data)
}
func (p *PeroContext) GetPage() (pageIndex int64, pageSize int64) {
	pageIndex, _ = strconv.ParseInt(p.Query("page_index"), 10, 64)
	pageSize, _ = strconv.ParseInt(p.Query("page_size"), 10, 64)
	if pageIndex <= 0 {
		pageIndex = 1
	}
	if pageSize <= 0 {
		pageSize = 15
	}
	return
}

func (p *PeroContext) reset() {
	p.Context = nil
}
