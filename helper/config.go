package helper

import (
	"github.com/jinzhu/configor"
)

type Config struct {
	SshHost []SshConf `json:"sshHost"`
}

func LoadConf() Config {
	conf := Config{}
	configor.Load(&conf, "/Users/chenrenhuan/server/config.json")
	return conf
}
