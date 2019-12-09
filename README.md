## 使用go自动登录ssh

开发需要, 要ssh登录的机器太多, 之前一直用`Termius`这个工具在连接第三方连接, 但是这个工具只是提供了简单的登录功能, 多次都是因为登录上去拉代码, 或者执行一些特定操作, 重复劳动, 程序员这种生物就是喜欢偷懒, 就算少敲一个字符都舒服, 于是乎在网上找了找`golang`写的ssh登录, 自己稍微改造了一下

> PS 如果你想 build 后使用, 可自行build


### 使用

> first

首先你得配置一个配置文件, 和`config.json`中一样, 然后修改代码中你的配置文件位置

> 参数

```golang
go run main.go -h

-c    不执行配置中的 cmd 项
-cmd
    批量执行 cmd
-n int
    登录机器编号

```

> 用法

```bash
go run main.go


1 】 测试机(xxx.xxxx.xxxx:22)

2 】 vvvvv(xx.xx.xx.xx:2452)

3 】 ......

选择你要登录的机器的编号:

```


```bash
go run main.go -n 4 -c


登录了第4台机器, 并且不执行配置中的 cmd 命令
```

```bash
go run main.go -cmd

1 】 测试机(xxx.xxxx.xxxx:22)

2 】 vvvvv(xx.xx.xx.xx:2452)

3 】 ......

请选择要执行命令的机器, 用逗号分隔, 0为所有机器:
1,2
你选择了这些机器:
机器列表:

1 】 报名系统测试机(118.24.199.162:22)
------------------------------------------
2 】 报名系统正式环境(123.59.150.140:22)
------------------------------------------
请输入你要批量执行的命令:
ls
结果列表:

------------------------------xxxxxxx:22)------------------------------

mysql57-community-release-el7-10.noarch.rpm

------------------------------报xxxx:22)------------------------------

a.php
confl_back.sql
cwd_user_attribute.ibd
dump.rdb
erlang-solutions-1.0-1.noarch.rpm
foo.sql
go
gogs-repositories
MP_verify_rArw0UFkfm27JTVT.txt
nohup.out
signup-201911060400.sql
```