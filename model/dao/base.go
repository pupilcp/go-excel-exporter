package dao

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
	"github.com/spf13/viper"
	"log"
)

var DB *gorm.DB

func SetDB(config *viper.Viper) {
	var err error
	DB, err = gorm.Open("mysql",
		fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?charset=utf8&parseTime=true&timeout=%ds",
			config.GetString("mysql.user"),
			config.GetString("mysql.password"),
			config.GetString("mysql.host"),
			config.GetString("mysql.port"),
			config.GetString("mysql.database"),
			int(config.GetInt64("mysql.timeout")),
		))
	if err != nil {
		errMsg := fmt.Sprintf("db connect err:%+v", err)
		log.Fatal(errMsg)
	}
	DB.SingularTable(true)
	DB.DB().SetMaxOpenConns(int(config.GetInt64("mysql.maxOpenConn")))
	DB.DB().SetMaxIdleConns(int(config.GetInt64("mysql.maxIdleConn")))
}
