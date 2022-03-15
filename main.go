package main

import (
	"github.com/pupilcp/go-excel-exporter/global"
	"github.com/pupilcp/go-excel-exporter/router"
	"github.com/pupilcp/go-excel-exporter/service"
	"net/http"
)

func main() {
	global.InitConfig()
	r := router.InitRoute()
	r.StaticFS("/download", http.Dir(global.Config.GetString("system.downloadPath")))
	go service.DownloadService.HandleTask()
	// 监听端口，默认在8080
	r.Run()
}
