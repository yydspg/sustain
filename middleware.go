package sustain

import (
	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	secure2 "github.com/unrolled/secure"
	"github.com/yydspg/sustain/cache"
	"net/http"
	"strings"
)

func CORSMiddleware() PeroHandlerFunc {
	return func(c *PeroContext) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, token, accept, origin, Cache-Control, X-Requested-With, appid, noncestr, sign, timestamp")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT,DELETE,PATCH")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}
		c.Next()
	}
}

func AuthMiddleware(cache cache.Cache, tokenPrefix string) PeroHandlerFunc {
	return func(c *PeroContext) {
		token := c.GetHeader("token")
		if token == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "token empty",
			})
			return
		}
		uidAndName := GetLoginUID(token, tokenPrefix, cache)
		if uidAndName == "" {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "login first",
			})
			return
		}
		args := strings.Split(uidAndName, "@")
		if len(args) != 2 {
			c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
				"msg": "internal error",
			})
			return
		}
		c.Set("uid", args[0])
		c.Set("name", args[1])
		c.Next()
	}
}

func TLSMiddleware(sslAddr string) PeroHandlerFunc {
	return func(c *PeroContext) {
		secureMiddleware := secure2.New(secure2.Options{
			SSLRedirect: true,
			SSLHost:     sslAddr,
		})
		err := secureMiddleware.Process(c.Writer, c.Request)

		if err != nil {
			return
		}

		c.Next()
	}
}
func TestMiddleware() PeroHandlerFunc {
	return func(c *PeroContext) {
		log.Logger.Log().Msg("sustain : test middleware")
		c.Next()
	}
}
