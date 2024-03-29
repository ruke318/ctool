package command

import (
	"ctool/helper"
	"fmt"
	"io"
	"os"
	"os/signal"
	"syscall"
	"time"

	"golang.org/x/crypto/ssh"
	"golang.org/x/crypto/ssh/terminal"
)

type SshLogin struct {
	curr   helper.SshConf
	notCmd bool
}

type SSHTerminal struct {
	Session *ssh.Session
	exitMsg string
	stdout  io.Reader
	stdin   io.Writer
	stderr  io.Reader
}

//当前登录机器
func (s *SshLogin) Login(number int, clear bool) {
	var conf []helper.SshConf
	config := helper.LoadConf()
	conf = config.SshHost
	//自定义接收参数
	if number == 0 {
		number = SelectHost(conf) + 1
	}
	if number <= 0 {
		fmt.Errorf("exit: (%d) \n", -2)
		return
	}
	// for {
	// 	if number > 0 && number <= len(conf) {
	// 		break
	// 	}
	// 	number = SelectHost(conf) + 1
	// }
	s.curr = conf[number-1]
	s.notCmd = clear
	fmt.Printf("正在登录: %v 】%v(%v)... \n\n", number, s.curr.Desc, s.curr.Host)
	client, cerr := helper.SshConnect(s.curr)
	if cerr != nil {
		fmt.Errorf("登录错误: %s \n", cerr)
		return
	}
	defer client.Close()

	s.New(client.Client)
}

// 实例化一个session连接
func (sl *SshLogin) New(client *ssh.Client) error {

	session, err := client.NewSession()
	if err != nil {
		return err
	}
	defer session.Close()

	s := SSHTerminal{
		Session: session,
	}

	return s.interactiveSession(sl.notCmd, sl.curr.Cmd)
}

// copy 信息展示
func (t *SSHTerminal) interactiveSession(notCmd bool, cmdString string) error {

	defer func() {
		if t.exitMsg == "" {
			fmt.Fprintln(os.Stdout, "the connection was closed on the remote side on ", time.Now().Format(time.RFC822))
		} else {
			fmt.Fprintln(os.Stdout, t.exitMsg)
		}
	}()

	fd := int(os.Stdin.Fd())
	state, err := terminal.MakeRaw(fd)
	if err != nil {
		return err
	}
	defer terminal.Restore(fd, state)

	termWidth, termHeight, err := terminal.GetSize(fd)
	if err != nil {
		return err
	}

	termType := os.Getenv("TERM")
	if termType == "" {
		termType = "xterm-256color"
	}

	err = t.Session.RequestPty(termType, termHeight, termWidth, ssh.TerminalModes{})
	if err != nil {
		return err
	}

	t.updateTerminalSize()

	t.stdin, err = t.Session.StdinPipe()
	if err != nil {
		return err
	}
	t.stdout, err = t.Session.StdoutPipe()
	if err != nil {
		return err
	}
	t.stderr, err = t.Session.StderrPipe()

	go io.Copy(os.Stderr, t.stderr)
	go io.Copy(os.Stdout, t.stdout)
	go func() {
		buf := make([]byte, 128)
		for {
			n, err := os.Stdin.Read(buf)
			if err != nil {
				fmt.Println(err)
				return
			}
			if n > 0 {
				_, err = t.stdin.Write(buf[:n])
				if err != nil {
					fmt.Println(err)
					t.exitMsg = err.Error()
					return
				}
			}
		}
	}()

	//---执行自定义脚本
	if !notCmd && cmdString != "" {
		t.Session.Run(cmdString)
	}

	err = t.Session.Shell()
	if err != nil {
		return err
	}
	err = t.Session.Wait()
	if err != nil {
		return err
	}
	return nil
}

// 窗口大小
func (t *SSHTerminal) updateTerminalSize() {

	go func() {
		// SIGWINCH is sent to the process when the window size of the terminal has
		// changed.
		sigwinchCh := make(chan os.Signal, 1)
		signal.Notify(sigwinchCh, syscall.SIGWINCH)

		fd := int(os.Stdin.Fd())
		termWidth, termHeight, err := terminal.GetSize(fd)
		if err != nil {
			fmt.Println(err)
		}

		for {
			select {
			// The client updated the size of the local PTY. This change needs to occur
			// on the server side PTY as well.
			case sigwinch := <-sigwinchCh:
				if sigwinch == nil {
					return
				}
				currTermWidth, currTermHeight, err := terminal.GetSize(fd)

				// Terminal size has not changed, don't do anything.
				if currTermHeight == termHeight && currTermWidth == termWidth {
					continue
				}

				t.Session.WindowChange(currTermHeight, currTermWidth)
				if err != nil {
					fmt.Printf("Unable to send window-change reqest: %s.", err)
					continue
				}

				termWidth, termHeight = currTermWidth, currTermHeight

			}
		}
	}()

}
