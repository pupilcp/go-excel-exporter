package global

import (
	"github.com/pupilcp/go-excel-exporter/model/dao"
)

func InitConfig() {
	SetConfig()
	SetLogger()
	dao.SetDB(Config)
}
