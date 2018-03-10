package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"time"
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

var workerChan = make(chan *net.TCPConn, 10000)
var graddr *string

func main() {

	// vman := Man{
	// 	&Em{
	// 		good: "goodman",
	// 	},
	// 	"helloman"}
	// vman.hello()
	// fmt.Println(vman)

	nf := flag.String("name", "server", "input server or client")
	graddr = flag.String("ip", "127.0.0.1:999", "ip:port")
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

	//plan2 code start
	for i := 0; i < 10; i++ {
		go sworker2(workerChan)
	}

	st := time.Now()
	var cn int
	cn = 0
	//plan2 code end
	for {

		tcpConn, cerr := listen.AcceptTCP()
		if cerr == nil {
			cn++
			//go sworker(tcpConn) //plan 1
			workerChan <- tcpConn //plan 2 减少了工作协程创建
		}
		if time.Since(st).Seconds() >= 2 {
			fmt.Println("accept tcp client :", cn)
			cn = 0
			st = time.Now()
		}
	}
}

/**
 * 通过chan，的工作线程
 */
func sworker2(conns <-chan *net.TCPConn) {
	fmt.Println("sworker2")
	var b [1500]byte
	for conn := range conns {
		//for {
		_, err := conn.Read(b[0:])

		if err == nil {
			//fmt.Println(string(b[0:len]))
		}

		conn.Write([]byte("HELLO WORLD"))
		//conn.Close()
		//}
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
	//linux 默认文件数打开1024
	//ulimit -n 20000 修改
	for i := 0; i < 40000; i++ {
		fmt.Println("create clint", i)
		go func() {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println(err)
				}
			}()
			raddr, rerr := net.ResolveTCPAddr("tcp", *graddr)
			if rerr != nil {
				fmt.Println(rerr)
				os.Exit(3)
			}
			tcpConn, dterr := net.DialTCP("tcp", nil, raddr)

			if dterr != nil {
				fmt.Println(dterr)
				os.Exit(4)
			}

			var b [1500]byte

			wlen, err := tcpConn.Write([]byte("ARE-U-THERE"))
			if err == nil && wlen > 0 {
				rlen, _ := tcpConn.Read(b[0:])
				fmt.Println(string(b[0:rlen]))
			} else {
				fmt.Println(err)
			}
		}()
	}
}
