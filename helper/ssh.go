package helper

import (
	"errors"
	// "fmt"
	"io/ioutil"
	"time"

	"github.com/mitchellh/go-homedir"
	"golang.org/x/crypto/ssh"
)

type SSHTerminal struct {
	Client *ssh.Client
	Msg    string
}

type SshConf struct {
	User string `json:"user"`
	Host string `json:"host"`
	Key  string `json:"key"`
	Pass string `json:"pass"`
	Desc string `json:"desc"`
	Cmd  string `json:"cmd"`
	Num  int
}

// 执行 命令返回结果
func (t *SSHTerminal) RunCmd(cmd string) (string, error) {
	session, err := t.Client.NewSession()
	if err != nil {
		return "", err
	}
	defer session.Close()
	ret, err := session.CombinedOutput(cmd)
	if err != nil {
		return "", err
	}
	return string(ret), nil
}

// 关闭客户端连接
func (t *SSHTerminal) Close() {
	t.Client.Close()
}

// 创建客户端
func SshConnect(conf SshConf) (SSHTerminal, error) {
	//创建sshp登陆配置
	config := &ssh.ClientConfig{
		Timeout:         time.Second, //ssh 连接time out 时间一秒钟, 如果ssh验证错误 会在一秒内返回
		User:            conf.User,
		HostKeyCallback: ssh.InsecureIgnoreHostKey(), //这个可以， 但是不够安全
	}
	if conf.Key == "" {
		config.Auth = []ssh.AuthMethod{ssh.Password(conf.Pass)}
	} else {
		data, err := publicKeyAuthFunc(conf.Key)
		if err != "" {
			eData := SSHTerminal{
				Msg: err,
			}
			return eData, errors.New(err)
		}
		config.Auth = []ssh.AuthMethod{data}
	}

	//dial 获取ssh client
	sshClient, err := ssh.Dial("tcp", conf.Host, config)
	if err != nil {
		// log.Fatal("创建ssh client 失败",err)
		errMsg := SSHTerminal{
			Msg: "创建ssh client 失败",
		}
		return errMsg, err
	}
	// //返回terminal
	terminal := SSHTerminal{
		Client: sshClient,
	}
	return terminal, nil
}

// publickey login
func publicKeyAuthFunc(kPath string) (ssh.AuthMethod, string) {
	keyPath, err := homedir.Expand(kPath)
	if err != nil {
		return nil, "find key's home dir failed"
	}
	key, err := ioutil.ReadFile(keyPath)
	if err != nil {
		return nil, "ssh key file read failed"
	}
	// Create the Signer for this private key.
	signer, err := ssh.ParsePrivateKey(key)
	if err != nil {
		return nil, "ssh key signer failed"
	}
	return ssh.PublicKeys(signer), ""
}
