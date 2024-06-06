package config

import (
	"errors"
	"fmt"
	"github.com/unknwon/goconfig"
	"log"
	"os"
)

const configFile = "/conf/conf.ini"

var File *goconfig.ConfigFile

/*
*
程序加载此文件的时候会先走初始化
*/
func init() {
	//拿到当前程序的目录
	currentDir, err := os.Getwd()
	if err != nil {
		panic(err)
	}
	configPath := currentDir + configFile

	//如果这个文件不存在
	if !fileExist(configPath) {
		panic(errors.New("配置文件不存在"))
	}

	//参数 sangserver.exe D:/xxx
	i := len(os.Args)
	if i > 1 {
		dir := os.Args[1]
		if dir != "" {
			configPath = dir + configFile
		}
	}

	//加载文件 文件系统的读取
	File, err = goconfig.LoadConfigFile(configPath)
	if err != nil {
		log.Fatal("读取配置文件出错:", err)
	}
	fmt.Println(File)
}

/*
*
判断文件是否存在
*/
func fileExist(fileName string) bool {
	_, err := os.Stat(fileName)
	return err == nil || os.IsExist(err)

}

func A() {
	fmt.Println("aa")
}
