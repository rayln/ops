package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type DatabaseConfig struct {
	MaxIdleConns int
	MaxOpenConns int
	ShowSQL      bool
	LogLevel     int //1 debug 2 info
	Master       string
	Select       string
	Redis        string
	RedisPwd     string
}

func (that *DatabaseConfig) Init() *DatabaseConfig {
	path, _ := os.Getwd()
	configDir := path + "/file/config"
	configPath := configDir + "/database_config.json"
	//打开config文件
	file, errfile := os.Open(configPath)
	defer file.Close()
	if errfile != nil {
		//如果不存在，则创建文件夹和config.json文件，并且把初始struct放入config里面
		initconfig := DatabaseConfig{
			MaxIdleConns: 5,
			MaxOpenConns: 10,
			ShowSQL:      true,
			LogLevel:     1,
			Master:       "",
			Select:       "",
			Redis:        "",
		}
		os.MkdirAll(configDir, os.ModeDir)
		b, _ := json.Marshal(initconfig)
		ioutil.WriteFile(configPath, b, 0666)
		return &initconfig
	}
	//把config文件内容放入AppConfig中，并且返回
	decoder := json.NewDecoder(file)
	err := decoder.Decode(that)
	if err != nil {
		panic(err.Error())
	}
	return that
}
