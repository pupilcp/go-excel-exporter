package global

import (
	"github.com/spf13/viper"
	"log"
)

var Config *viper.Viper

func SetConfig() {
	var err error
	Config = viper.New()
	Config.SetConfigType("toml")
	Config.SetConfigFile("config/config.toml")
	err = Config.ReadInConfig()
	if err != nil {
		log.Fatalf("read config failed: %+v", err)
	}
}
