package model

type DownloadTask struct {
	TaskId        int    `gorm:"column:task_id;primary_key;AUTO_INCREMENT" json:"task_id"`
	UserId        int    `gorm:"column:user_id;default:0;NOT NULL" json:"user_id"`
	FileName      string `gorm:"column:file_name;NOT NULL" json:"file_name"`               // 文件名
	RequestUrl    string `gorm:"column:request_url;NOT NULL" json:"request_url"`           // 请求url
	RequestParams string `gorm:"column:request_params;NOT NULL" json:"request_params"`     // 请求参数
	RequestMethod string `gorm:"column:request_method;NOT NULL" json:"request_method"`     // 请求方法，POST/GET
	DownloadFile  string `gorm:"column:download_file;NOT NULL" json:"download_file"`       // 下载文件地址
	TaskStatus    int    `gorm:"column:task_status;default:0;NOT NULL" json:"task_status"` // 任务状态，0：未处理，1：处理成功，2：处理失败
	Remark        string `gorm:"column:remark;NOT NULL" json:"remark"`                     // 备注信息
	CreatedAt     int    `gorm:"column:created_at;default:0;NOT NULL" json:"created_at"`   // 任务创建时间
	UpdatedAt     int    `gorm:"column:updated_at;default:0;NOT NULL" json:"updated_at"`   // 任务更新时间
}
