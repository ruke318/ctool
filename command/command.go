package command

import (
	"ctool/helper"
	"fmt"
	"strings"
)

// 现实机器列表
func ShowList(list []helper.SshConf) {
	fmt.Println("机器列表:\n")
	for key, value := range list {
		fmt.Printf("%v 】 %v(%v) \n", key+1, value.Desc, value.Host)
		fmt.Println(strings.Repeat("-", 42))
	}
}
