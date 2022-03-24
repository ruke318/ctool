package main

import (
	"ctool/command"
	"ctool/helper"
	"flag"
	"fmt"
	"log"
	"os"
	"os/user"
)

// 定义一些参数
var runCmd = flag.Bool("cmd", false, "批量执行 cmd")
var num = flag.Int("n", 0, "登录机器编号")
var clear = flag.Bool("c", false, "不执行配置中的 cmd 项")

var configPath = flag.String("confPath", "", "设置配置文件目录")
var listFile = flag.String("f", "", "列表配置文件名称")

var APP_PATH = "/CtoolDir/c.json"

func getAppPath() string {
	usr, err := user.Current()
	if err != nil {
		log.Fatal(err)
		os.Exit(-2)
		return ""
	}
	path := usr.HomeDir + APP_PATH
	if !fileExists(path) {
		os.Create(path)
	}
	return path
}

// 初始化应用文件位置
func init() {
	flag.Parse()
	path := getAppPath()
	helper.AppConfPath = path
	// 检测是否存在配置文件目录
	c := helper.ReadConf(path)
	if len(c.HostListFile) == 0 && len(*configPath) == 0 {
		fmt.Println("plase run -- `ctool -confPath=your confgi dir`  set config dir")
		os.Exit(-4)
		return
	}
}

func main() {
	// 设置配置文件目录
	if len(*configPath) > 0 {
		helper.SetConfigDir(getAppPath(), *configPath)
		fmt.Println("set config path success !", *configPath)
		return
	}

	if len(*listFile) > 0 {
		helper.ListConfigName = *listFile
	}

	if *runCmd {
		cmd := &command.SshCmd{}
		cmd.Run()
		return
	}

	c := &command.SshLogin{}
	c.Login(*num, *clear)
	return
}

func fileExists(name string) bool {
	if _, err := os.Stat(name); err != nil {
		if os.IsNotExist(err) {
			return false
		}
	}
	return true
}
