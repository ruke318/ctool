package command

import (
	"ctool/helper"
	"fmt"
	"strconv"
	"strings"

	"github.com/manifoldco/promptui"
)

// 现实机器列表
func SelectHost(list []helper.SshConf) int {
	helpList := make([]helper.SshConf, 0)
	for kk, item := range list {
		item.Num = kk + 1
		helpList = append(helpList, item)
	}

	templates := &promptui.SelectTemplates{
		Label:    "{{ . }}?",
		Active:   "\U0001F336 {{ .Num | cyan }} => {{ .Desc | cyan }} ({{ .Host | red }})",
		Inactive: "  {{ .Num | cyan }} {{ .Desc | cyan }} ({{ .Host | red }})",
		Selected: "\U0001F336 {{ .Num | cyan }} {{ .Desc | red | cyan }} ({{ .Host | red }})",
		Details: `
--------- detail ----------
{{ "IP端口:"| faint }}	{{ .Host }}
{{ "简  介:"| faint }}	{{ .Desc }}
{{ "用户名:"| faint }}	{{ .User }}
{{ "默认命令:" | faint }}	{{ .Cmd }}`,
		Help: `{{ "支持如下操作:" | faint }} {{ .NextKey | faint }} ` +
			`{{ .PrevKey | faint }} {{ .PageDownKey | faint }} {{ .PageUpKey | faint }} {{"选择机器" | faint}}` +
			`{{ if .Search }} {{ "或者" | faint }} {{ .SearchKey | faint }} {{ "进行搜索" | faint }}{{ end }}`,
	}

	searcher := func(input string, index int) bool {
		strLine, _ := strconv.Atoi(input)
		if strLine == index+1 {
			return true
		} else if strLine > 0 {
			return false
		}

		pepper := helpList[index]
		name := strings.Replace(strings.ToLower(pepper.Desc), " ", "", -1)
		input = strings.Replace(strings.ToLower(input), " ", "", -1)

		return strings.Contains(name, input)
	}

	prompt := promptui.Select{
		Label:     "选择机器",
		Items:     helpList,
		Templates: templates,
		Size:      10,
		Searcher:  searcher,
	}

	i, _, err := prompt.Run()

	if err != nil {
		fmt.Printf("Prompt failed %v\n", err)
		return -1
	}
	return i
}

func ShowList(list []helper.SshConf) {
	fmt.Println("机器列表:\n")
	for key, value := range list {
		fmt.Printf("%v 】 %v(%v) \n", key+1, value.Desc, value.Host)
		fmt.Println(strings.Repeat("-", 42))
	}
}
