package helper

import (
	"encoding/json"
	"io"
	"log"
	"os"
	"strings"

	"github.com/jinzhu/configor"
)

type Config struct {
	SshHost []SshConf `json:"sshHost"`
}

func LoadConf() Config {
	appConf := ReadConf(AppConfPath)

	basePath := strings.TrimRight(appConf.HostListFile, "/") + "/"
	conf := Config{}
	if len(ListConfigName) == 0 {
		ListConfigName = "config"
	}
	filePath := basePath + ListConfigName + ".json"
	configor.Load(&conf, filePath)
	return conf
}

type AppConf struct {
	HostListFile string `json:"hostListFile"`
}

var AppConfPath string    // app 配置路径
var ListConfigName string // 列表配置文件名前缀

func ReadConf(path string) AppConf {
	conf := AppConf{}
	configor.Load(&conf, path)
	return conf
}

func SetConfigDir(appPath string, configDir string) {
	conf := AppConf{}
	conf.HostListFile = configDir
	af, err := os.OpenFile(appPath, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}
	defer af.Close()

	js, _ := json.Marshal(conf)
	io.WriteString(af, string(js))
}
