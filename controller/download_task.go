package controller

import (
	"github.com/gin-gonic/gin"
	"github.com/pupilcp/go-excel-exporter/controller/request"
	"github.com/pupilcp/go-excel-exporter/service"
	"net/http"
	"strconv"
)

type downloadTaskController struct {
}

var DownloadTaskController downloadTaskController

func (c *downloadTaskController) Create(ctx *gin.Context) {
	var err error
	var params request.CreateTaskReq
	if err = ctx.ShouldBindJSON(&params); err != nil {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 1001,
			"info": err.Error(),
		})
		return
	}
	_, err = service.DownloadTaskService.CreateTask(params.UserId, params.FileName, params.RequestUrl, params.RequestParams, "POST")
	if err != nil {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 1002,
			"info": err.Error(),
		})
	} else {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 1,
			"info": "success",
		})
	}
}

func (c *downloadTaskController) GetList(ctx *gin.Context) {
	uid := ctx.DefaultQuery("user_id", "0")
	p := ctx.DefaultQuery("page", "1")
	ps := ctx.DefaultQuery("page_size", "10")
	if uid == "0" {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code":    1001,
			"message": "Params error",
		})
		return
	}
	userId, _ := strconv.Atoi(uid)
	page, _ := strconv.Atoi(p)
	pageSize, _ := strconv.Atoi(ps)
	list, total := service.DownloadTaskService.GetTaskList(userId, -1, page, pageSize)
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code": 1,
		"info": "success",
		"data": map[string]interface{}{
			"list":  list,
			"total": total,
		},
	})
}

func (c *downloadTaskController) Retry(ctx *gin.Context) {
	taskId := ctx.PostForm("task_id")
	if taskId == "" {
		ctx.JSON(http.StatusOK, map[string]interface{}{
			"code": 1001,
			"info": "Params error",
		})
		return
	}
	tid, _ := strconv.Atoi(taskId)
	service.TaskChannel <- tid
	ctx.JSON(http.StatusOK, map[string]interface{}{
		"code":    1,
		"message": "success",
	})
}
