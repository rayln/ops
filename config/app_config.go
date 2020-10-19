package config

import (
	"encoding/json"
	"io/ioutil"
	"os"
)

type AppConfig struct {
	AppName      string
	Port         string
	WsPort       string
	RelativePath string
	StaticPath   string
	Mode         string
	IsSign       bool //是否开启验签
}

func (that *AppConfig) Init() *AppConfig {
	path, _ := os.Getwd()
	configDir := path + "/file/config"
	configPath := configDir + "/app_config.json"
	//打开config文件
	file, errfile := os.Open(configPath)
	defer file.Close()
	if errfile != nil {
		//如果不存在，则创建文件夹和config.json文件，并且把初始struct放入config里面
		initconfig := AppConfig{
			AppName:      "ops",
			Port:         ":8009",
			RelativePath: "/",
			StaticPath:   "/static",
			Mode:         "dev",
			IsSign:       true,
		}
		os.MkdirAll(configDir, os.ModeDir)
		b, _ := json.Marshal(initconfig)
		ioutil.WriteFile(configPath, b, 0666)
		return &initconfig
	}
	//把config文件内容放入AppConfig中，并且返回
	decoder := json.NewDecoder(file)
	//conf:=AppConfig{}
	err := decoder.Decode(that)
	if err != nil {
		panic(err.Error())
	}
	return that
}
