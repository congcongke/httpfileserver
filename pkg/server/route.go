package server

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"github.com/congcongke/httpfileserver/pkg/config"
	"github.com/congcongke/httpfileserver/pkg/middleware"
)

func LoadFromConfig(conf *config.Config) *gin.Engine {
	e := gin.Default()

	e.GET("", func(c *gin.Context) {
		c.JSON(http.StatusOK, "pong")
	})

	lfh := NewLocalFileHandle(conf.RootPath)
	fileGroup := e.Group("/file/v1",
		middleware.ReqLoggerMiddleware(), middleware.BasicAuthHandle(conf.Auth.Username, conf.Auth.Password))

	fileGroup.GET("/:filename", lfh.Get)
	fileGroup.PUT("/:filename", lfh.Put)
	fileGroup.GET("", lfh.List)

	return e
}
