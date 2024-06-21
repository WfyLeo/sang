package config

import (
	"errors"
	"github.com/unknwon/goconfig"
	"log"
	"os"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

// 程序加载的时候，会先走初始化方法
func init() {
	//程序的当前目录
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile
	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}
	//参数 newGoSang.exe
	len := len(os.Args)
	if len > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}

	//文件系统读取
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Println("读取配置文件出错:", err)
	}
}

func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)
}

func A() {

}
