package router

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaobeicn/go-es-location/es"
	"github.com/xiaobeicn/go-es-location/router/handler"
	"net/http"
)

func InitRouter(r *gin.Engine) *gin.Engine {
	r.NoRoute(func(c *gin.Context) {
		if c.Request.RequestURI == "/favicon.ico" {
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该路由",
		})
		return
	})
	r.NoMethod(func(c *gin.Context) {
		if c.Request.RequestURI == "/favicon.ico" {
			return
		}
		c.JSON(http.StatusNotFound, gin.H{
			"code": 404,
			"msg":  "找不到该方法",
		})
		return
	})

	r.GET("", handler.HomeHandler)
	r.GET("/near", handler.NearHandler)
	r.GET("/ip", handler.IpHandler)

	r.GET("/es", func(c *gin.Context) {
		index, err := es.CreateIndex(c)
		if err != nil {
			c.JSON(500, gin.H{
				"err": err.Error(),
			})
		} else {
			c.JSON(200, gin.H{
				"index": index,
			})
		}
	})

	return r
}
