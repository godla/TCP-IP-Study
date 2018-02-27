package main

import (
	"fmt"
	"net"
)

/**
UDP是无链接！ UDP是无链接！ UDP是无链接！

DialUDP 是 pre-connected
ListenUDP 是 unconnect

如果*UDPConn是connected,读写方法是Read和Write。
如果*UDPConn是unconnected,读写方法是ReadFromUDP和WriteToUDP（以及ReadFrom和WriteTo)。

如果使用dail
你将失去SetKeepAlive或TCPConn和UDPConn的SetReadBuffer 这些函数， 除非做类型转换
*/
func main() {
	var wait string
	Dc()
	fmt.Scanln(&wait)
}

func Dc() {
	var raddr *net.UDPAddr
	raddr, err := net.ResolveUDPAddr("udp", "127.0.0.1:9123")
	conn, err := net.DialUDP("udp", nil, raddr)

	if err != nil {
		fmt.Println(err.Error())
	}

	_, err = conn.Write([]byte("ARE-U-THERE"))

	if err != nil {
		fmt.Println(err.Error())
	}

	var buf [1500]byte

	rlen, err := conn.Read(buf[0:])

	if err != nil {
		fmt.Println(err.Error())
	}
	fmt.Println(string(buf[0:rlen]))
}

func Ds() {
	var laddr *net.UDPAddr

	laddr, err := net.ResolveUDPAddr("udp", ":9123")

	conn, err := net.ListenUDP("udp", laddr)

	if err != nil {
		fmt.Println(err.Error())
	}

	go func(conn *net.UDPConn) {
		var rbuf [1500]byte

		for {
			rlen, raddr, err := conn.ReadFromUDP(rbuf[0:])

			if err != nil {
				fmt.Println(err.Error())
			}

			fmt.Println(string(rbuf[0:rlen]))

			conn.WriteToUDP([]byte("I-AM-HERE"), raddr)
		}
	}(conn)
}
