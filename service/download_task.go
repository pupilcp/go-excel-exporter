package service

import (
	"encoding/json"
	"github.com/pupilcp/go-excel-exporter/model"
	"github.com/pupilcp/go-excel-exporter/model/dao"
	"time"
)

type downloadTaskService struct {
}

var TaskChannel = make(chan int, 100)

var DownloadTaskService downloadTaskService

func (s *downloadTaskService) GetTaskList(userId, status, page, pageSize int) ([]model.DownloadTask, int) {

	list := dao.DownloadTaskDao.GetList(userId, status, page, pageSize)
	total := dao.DownloadTaskDao.GetCount(userId, status)

	return list, total
}

func (s *downloadTaskService) CreateTask(userId int, fileName, requestUrl string, requestParams interface{}, requestMethod string) (bool, error) {
	requestParamsStr, _ := json.Marshal(requestParams)
	taskData := model.DownloadTask{
		UserId:        userId,
		FileName:      fileName,
		RequestUrl:    requestUrl,
		RequestParams: string(requestParamsStr),
		RequestMethod: requestMethod,
		DownloadFile:  "",
		TaskStatus:    0,
		CreatedAt:     int(time.Now().Unix()),
		UpdatedAt:     int(time.Now().Unix()),
	}
	taskId, err := dao.DownloadTaskDao.CreateTask(taskData)
	if err != nil {
		return false, err
	}
	TaskChannel <- taskId
	return true, err
}

func (s *downloadTaskService) UpdateTask(taskId int, file, remark string, taskStatus int) error {
	err := dao.DownloadTaskDao.UpdateTask(taskId, file, remark, taskStatus)
	if err != nil {
		return err
	}
	return err
}
