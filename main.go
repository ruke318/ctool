package main

import (
	"ctool/command"
	"flag"
)

// 定义一些参数
var runCmd = flag.Bool("cmd", false, "批量执行 cmd")
var num = flag.Int("n", 0, "登录机器编号")
var clear = flag.Bool("c", false, "不执行配置中的 cmd 项")

func main() {
	flag.Parse()
	if *runCmd {
		cmd := &command.SshCmd{}
		cmd.Run()
		return
	}

	c := &command.SshLogin{}
	c.Login(*num, *clear)
	return
}
