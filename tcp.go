package main

import (
	"flag"
	"fmt"
	"net"
	"os"
)

/**
 *  tcp 启动 链接 三次握手  关闭四次
 *  慢启动 + 拥塞窗口
 *  门限控制 < 拥塞窗口 进入拥塞避免
 *  在发送数据后 检测 ack 确认超时 或者 收到重复ack 确认 操作门限值 和 拥塞窗口值 来限流
 *  接收方 使用通告窗口 告知发送可接受多少字节
 *
 *  发送方：滑动窗口协议
 */

//结构体内嵌接口
type Man struct {
	iTest
	name string
}

type iTest interface {
	hello()
}

type Em struct {
	good string
}

func (em *Em) hello() {
	fmt.Println(em.good)
}

func (man *Man) hello() {
	fmt.Println(man.name)
}

func main() {

	// vman := Man{
	// 	&Em{
	// 		good: "goodman",
	// 	},
	// 	"helloman"}
	// vman.hello()
	// fmt.Println(vman)

	nf := flag.String("name", "server", "input server or client")
	flag.Parse()
	fmt.Println(*nf)
	if *nf == "server" {
		server()
	} else {
		client()
	}

	var wait string
	fmt.Scanln(&wait)
}

func server() {
	laddr, aerr := net.ResolveTCPAddr("tcp", ":999")
	if aerr != nil {
		os.Exit(1)
	}
	listen, lerr := net.ListenTCP("tcp", laddr)
	defer listen.Close()
	if lerr != nil {
		os.Exit(2)
	}
	for {
		tcpConn, cerr := listen.AcceptTCP()
		if cerr == nil {
			go sworker(tcpConn)
		}
	}
}

/**
 * 工作协程
 */
func sworker(conn *net.TCPConn) {
	fmt.Println("sworker")
	var b [1500]byte
	for {
		len, err := conn.Read(b[0:])

		if err == nil {
			fmt.Println(string(b[0:len]))
		}

		conn.Write([]byte("HELLO WORLD"))
	}
}

func client() {
	raddr, rerr := net.ResolveTCPAddr("tcp", "127.0.0.1:999")
	if rerr != nil {
		os.Exit(3)
	}
	tcpConn, dterr := net.DialTCP("tcp", nil, raddr)

	if dterr != nil {
		os.Exit(4)
	}

	var b [1500]byte

	wlen, err := tcpConn.Write([]byte("ARE-U-THERE"))
	if err == nil && wlen > 0 {
		rlen, _ := tcpConn.Read(b[0:])
		fmt.Println(string(b[0:rlen]))
	}

}
