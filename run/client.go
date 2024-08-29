package run

import (
	"bufio"
	"context"
	"flag"
	"fmt"
	"io"
	"log"
	"net"
	"os"
	"time"
)

var ServerIp string
var ServerPort int

func init() {
	flag.StringVar(&ServerIp, "ip", "127.0.0.1", "Server IP")
	flag.IntVar(&ServerPort, "port", 8888, "Server Port")
}
func RunClient() {
	flag.Parse()
	ctx, cancelFunc := context.WithCancel(context.Background())
	client := NewClient(ServerIp, ServerPort, ctx, cancelFunc)
	//go client.LMReadServer()
	defer client.conn.Close()

	for {

		fmt.Println("______________")
		fmt.Println("0 退出")
		fmt.Println("1 私聊")
		fmt.Println("2 公聊")
		fmt.Println("3 修改用户名")
		fmt.Println("______________")
		cmd := BestIn("")
		switch cmd {
		case "1":
			fmt.Println("进入私聊模式")
			client.PivateChat()
			break
		case "2":
			fmt.Println("进入公聊模式")
			client.ChatAll()
			break
		case "3":
			fmt.Println("请输入新用户名")
			client.Rename()
			break
		case "0":
			fmt.Println("退出当前模式")
			break
		default:
			fmt.Println("请输入正确的模式代码")

		}
	}

}

type Client struct {
	ServerIp   string
	ServerPort int
	Name       string
	conn       net.Conn
	ctx        context.Context
	cancelFunc context.CancelFunc
}

func NewClient(ServerIp string, ServerPort int, ctx context.Context, cancelFunc context.CancelFunc) *Client {

	cilent := &Client{
		ServerIp:   ServerIp,
		ServerPort: ServerPort,
		Name:       ServerIp,
		ctx:        ctx,
		cancelFunc: cancelFunc,
	}
	conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))
	if err != nil {
		log.Fatalln("net.Dial err", err)
		return nil
	}
	cilent.conn = conn
	return cilent
}

func (this *Client) LMReadServer() {
	io.Copy(os.Stdout, this.conn)
}

func (this *Client) Rename() {
	name := BestIn("输入用户名")
	fmt.Println("新用户名是", name)
	_, err := this.conn.Write([]byte("Rename|" + name + "|"))
	if err != nil {
		log.Println("conn Write err", err)
	}
	buf := make([]byte, 4096)
	// 设置读取操作的超时时间为5秒

	for i := 0; i < 100; i++ {
		err = this.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			fmt.Println("Error setting read deadline:", err)
			return
		}
		n, err := this.conn.Read(buf)
		if err != nil {
			//Timeout 强制下线
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))
				if err != nil {
					log.Fatalln("net.Dial err", err)
				}
				this.conn = conn
			}
			if err != io.EOF {
				log.Println("coon Read err ", err)
				log.Println("server close coon")
			} else {
				log.Println("No data read")
			}
		}

		buf_str := string(buf[:n])
		if buf_str == "|User name is used" {
			fmt.Println("User name is used")
			return
		} else if buf_str == "|Rename ok" {
			this.Name = name
			fmt.Println("Rename to", name)
			return
		}

	}

}

func (this *Client) PivateChat() {
	name := BestIn("输入用户名")
	msg := BestIn("请输私聊内容")
	//PivateChat|lcs|hello|
	_, err := this.conn.Write([]byte("PivateChat|" + name + "|" + "[" + this.Name + "]:" + msg))
	if err != nil {
		log.Println("conn Write err", err)
	}
	buf := make([]byte, 4096)
	// 设置读取操作的超时时间为5秒

	for i := 0; i < 100; i++ {
		err = this.conn.SetReadDeadline(time.Now().Add(60 * time.Second))
		if err != nil {
			fmt.Println("Error setting read deadline:", err)
			return
		}
		n, err := this.conn.Read(buf)
		if err != nil {
			//Timeout 强制下线
			if netErr, ok := err.(net.Error); ok && netErr.Timeout() {
				conn, err := net.Dial("tcp", fmt.Sprintf("%s:%d", ServerIp, ServerPort))
				if err != nil {
					log.Fatalln("net.Dial err", err)
				}
				this.conn = conn
			}
			if err != io.EOF {
				log.Println("coon Read err ", err)
				log.Println("server close coon")
			} else {
				log.Println("No data read")
			}
		}

		buf_str := string(buf[:n])
		if buf_str == "|Pivate chat grammatical errors" {
			fmt.Println("|Pivate chat grammatical errors")
			return
		} else if buf_str == "|PivateChat ok" {
			fmt.Println("|PivateChat ok", name)
			return
		} else if buf_str == "|Pivate chat no user" {
			fmt.Println("|Pivate chat no user", name)
			return
		}

	}

}

func (this *Client) ChatAll() {
	msg := BestIn("输入公聊消息")
	_, err := this.conn.Write([]byte(msg))
	if err != nil {
		log.Println("conn Write err", err)
	}
}

func BestIn(prompt string) (in string) {
	reader := bufio.NewReader(os.Stdin)
	if prompt != "" {
		fmt.Println(prompt)
	}
	in, err := reader.ReadString('\r')
	if err != nil {
		fmt.Println("读取输入时出错:", err)
		return
	}
	return in[:len(in)-1]
}
