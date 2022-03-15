package dao

import (
	"github.com/pupilcp/go-excel-exporter/model"
)

type downloadTaskDao struct {
}

var DownloadTaskDao = downloadTaskDao{}

func (d *downloadTaskDao) GetList(userId, status, page, pageSize int) []model.DownloadTask {
	var taskList []model.DownloadTask
	command := DB.Model(model.DownloadTask{}).Select("*")
	if userId > 0 {
		command = command.Where("user_id=?", userId)
	}
	if status != -1 {
		command = command.Where("task_status=?", status)
	}
	command = command.Limit(pageSize).Offset((page - 1) * pageSize)
	command.Scan(&taskList)
	return taskList
}

func (d *downloadTaskDao) GetCount(userId, status int) int {
	command := DB.Model(model.DownloadTask{})
	if userId > 0 {
		command = command.Where("user_id=?", userId)
	}
	if status != -1 {
		command = command.Where("task_status=?", status)
	}
	total := 0
	command.Count(&total)
	return total
}

func (d *downloadTaskDao) GetOne(taskId int) model.DownloadTask {
	var task model.DownloadTask
	command := DB.Model(model.DownloadTask{}).Select("*")
	if taskId > 0 {
		command = command.Where("task_id=?", taskId)
	}
	command.Scan(&task)
	return task
}

func (d *downloadTaskDao) CreateTask(data model.DownloadTask) (int, error) {
	DB.Model(model.DownloadTask{}).Create(&data)
	if DB.Error != nil {
		return 0, DB.Error
	}
	return data.TaskId, nil
}

func (d *downloadTaskDao) UpdateTask(taskId int, fileUrl, remark string, taskStatus int) error {
	//data := model.DownloadTask{
	//	TaskId: taskId,
	//	DownloadFile:  fileUrl,
	//	TaskStatus:    taskStatus,
	//	Remark: remark,
	//	//UpdatedAt: int(time.Now().Unix()),
	//}
	data := make(map[string]interface{})
	data["download_file"] = fileUrl
	data["remark"] = remark
	data["task_status"] = taskStatus
	DB.Model(model.DownloadTask{}).Where("task_id = ?", taskId).UpdateColumn(data)
	if DB.Error != nil {
		return DB.Error
	}
	return nil
}
