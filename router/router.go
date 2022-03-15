package router

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/pupilcp/go-excel-exporter/controller"
	"net/http"
)

func InitRoute() *gin.Engine {
	fmt.Println("InitRoute")
	r := gin.Default()
	// 2.绑定路由规则，执行的函数
	// gin.Context，封装了request和response
	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "hello World!")
	})
	api := r.Group("/api")
	{
		api.POST("/task/create", controller.DownloadTaskController.Create)
		api.POST("/task/retry", controller.DownloadTaskController.Retry)
		api.GET("/task/get-list", controller.DownloadTaskController.GetList)
	}

	return r
}
