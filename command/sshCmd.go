package command

import (
	"bufio"
	"ctool/helper"
	"fmt"
	"github.com/jinzhu/configor"
	"os"
	"strconv"
	"strings"
)

type SshCmd struct {
}

type result struct {
	keyTag int
	ret string
	info helper.SshConf
}

func (s *SshCmd) Run() {
	var conf []helper.SshConf
	configor.Load(&conf, "/Users/chenrenhuan/server/config.json")
	ShowList(conf)
	nums := ""
	for {
		if nums != "" {
			break
		}
		fmt.Println("请选择要执行命令的机器, 用逗号分隔, 0为所有机器:")
		fmt.Scanln(&nums)
	}

	//去除不存在的机器
	numArrStr := strings.Split(nums,",")
	var selected []helper.SshConf
	for _, value := range numArrStr {
		if value != "" {
			if value != "0" {
				key, _:= strconv.Atoi(value)
				if key > 0 {
					selected = append(selected, conf[key - 1])
				}
			} else {
				selected = conf
				break
			}
		}
	}
	if len(selected) == 0 {
		fmt.Println("你没有选择机器")
		return
	}
	fmt.Println("你选择了这些机器:")
	ShowList(selected)

	//输入命令
	cmdString := ""
	inputReader := bufio.NewReader(os.Stdin)
	for {
		if cmdString != "" {
			break
		}
		fmt.Println("请输入你要批量执行的命令:")
		cmdString, _ = inputReader.ReadString('\n')
	}

	//执行命令
	c := make(chan *result, len(conf))
	for key, item := range selected {
		//用协程去不同机器超找
		go s.getData(c, key, item, cmdString)
	}

	var data []*result
	for i := 0; i < len(selected); i++ {
		item := <-c
		item.info = selected[item.keyTag]
		data = append(data, item)
	}
	s.showData(data)
	return
}

//获取数据
func (s *SshCmd) getData(c chan *result, key int, conf helper.SshConf, cmdString string) {
	line := &result{keyTag: key, ret: ""}
	client, error := helper.SshConnect(conf)
	if error != nil {
		c <- line
		return
	}
	defer client.Close()

	res, err := client.RunCmd(cmdString)
	if err == nil {
		line.ret = res
	}
	c <- line
}

//结果展示
func (SshCmd) showData(data []*result) {
	fmt.Println("结果列表:\n")
	for _, value := range data {
		fmt.Print(strings.Repeat("-", 30))
		fmt.Printf("%v(%v)", value.info.Desc, value.info.Host)
		fmt.Println(strings.Repeat("-", 30) + "\n")
		fmt.Println(value.ret)
	}
}