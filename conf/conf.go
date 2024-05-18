package conf

import (
	"fmt"
	"os"
	"strings"

	"gopkg.in/ini.v1"
)

type Config struct {
	UserName string
	Password string
	Host     string
	Port     string

	LogLevel string
}

func (c Config) ConnectString() string {

	dsn0 := c.UserName + ":" +
		c.Password + "@(" +
		c.Host + ":" +
		c.Port + ")/" + "?charset=utf8mb4&parseTime=True&loc=Local"
	return dsn0
}

var DBConfig = Config{}

func InitConf() {
	cfg, err := ini.Load("./easysearch.ini")
	if err != nil {
		fmt.Printf("Fail to read file: %v", err)
		os.Exit(1)
	}

	DBConfig.UserName = cfg.Section("mysql").Key("username").String()
	DBConfig.Password = cfg.Section("mysql").Key("password").String()
	DBConfig.Host = cfg.Section("mysql").Key("host").String()
	DBConfig.Port = cfg.Section("mysql").Key("port").String()
	DBConfig.LogLevel = strings.ToLower(cfg.Section("").Key("loglevel").String())
}
