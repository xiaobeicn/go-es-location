package main

import (
	"github.com/gin-gonic/gin"
	"github.com/xiaobeicn/go-es-location/global"
	"github.com/xiaobeicn/go-es-location/initialize"
	"github.com/xiaobeicn/go-es-location/router"
	"net/http"
	"time"
)

func init() {
	initialize.Init()
}

func main() {
	gin.SetMode(global.GConfig.App.Mode)
	r := gin.Default()

	// 创建自定义配置服务
	httpServer := &http.Server{
		Addr:           global.GConfig.App.Addr,
		Handler:        r,
		ReadTimeout:    time.Second * 10,
		WriteTimeout:   time.Second * 10,
		MaxHeaderBytes: 1 << 20,
	}
	
	router.InitRouter(r)

	_ = httpServer.ListenAndServe()
}
