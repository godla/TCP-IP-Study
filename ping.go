package main

import (
	"errors"
	"fmt"
	"net"
	"os"
	// "runtime"
	// "sync"
	//"testing"
	"time"

	"./xnet/xnternal/iana" //use of internal package not allowed
	"./xnet/xnternal/nettest"
	"golang.org/x/net/icmp"
	"golang.org/x/net/ipv4"
	"golang.org/x/net/ipv6"
)

//-ldflags=-linkmode=internal
func main() {
	//使用go get 安装x/net 不能成功
	//	planB: git clone git@github.com:golang/net.git $GOPATH./src/golang.org/x/net/

	//这是golang net test 代码用来学习 删除ip6 的支持数据

	if m, ok := nettest.SupportsRawIPSocket(); !ok {
		fmt.Println(m)
		return
	} else {
		fmt.Println(m, ok)
	}

	for i, tt := range privilegedPingTests {
		fmt.Println(i, tt)
		if err := doPing(tt, i); err != nil {
			fmt.Println(err)
			return
		}
	}

}

func googleAddr(c *icmp.PacketConn, protocol int) (net.Addr, error) {
	const host = "www.baidu.com"
	ips, err := net.LookupIP(host)
	if err != nil {
		return nil, err
	}
	fmt.Println(ips)

	netaddr := func(ip net.IP) (net.Addr, error) {
		switch c.LocalAddr().(type) {
		case *net.UDPAddr:
			return &net.UDPAddr{IP: ip}, nil
		case *net.IPAddr:
			return &net.IPAddr{IP: ip}, nil
		default:
			return nil, errors.New("neither UDPAddr nor IPAddr")
		}
	}

	for _, ip := range ips {
		switch protocol {
		case iana.ProtocolICMP:
			fmt.Println("ICMP")
			fmt.Println(ip.To4())
			//os.Exit(1)
			if ip.To4() != nil {
				return netaddr(ip)
			}
		case iana.ProtocolIPv6ICMP:
			if ip.To16() != nil && ip.To4() == nil {
				return netaddr(ip)
			}
		}
	}
	return nil, errors.New("no A or AAAA record")
}

type pingTest struct {
	network, address string
	protocol         int       //协议
	mtype            icmp.Type //协议类型
}

var privilegedPingTests = []pingTest{

	{"ip4:icmp", "0.0.0.0", iana.ProtocolICMP, ipv4.ICMPTypeEcho},
	//{"ip6:ipv6-icmp", "::", iana.ProtocolIPv6ICMP, ipv6.ICMPTypeEchoRequest},
}

func doPing(tt pingTest, seq int) error {
	//net.ListenPacket(network, address)
	c, err := icmp.ListenPacket(tt.network, tt.address)
	if err != nil {
		return err
	}
	defer c.Close()

	dst, err := googleAddr(c, tt.protocol)
	if err != nil {
		return err
	}
	fmt.Println(dst)

	// if tt.network != "udp6" && tt.protocol == iana.ProtocolIPv6ICMP {
	// 	var f ipv6.ICMPFilter
	// 	f.SetAll(true)
	// 	f.Accept(ipv6.ICMPTypeDestinationUnreachable)
	// 	f.Accept(ipv6.ICMPTypePacketTooBig)
	// 	f.Accept(ipv6.ICMPTypeTimeExceeded)
	// 	f.Accept(ipv6.ICMPTypeParameterProblem)
	// 	f.Accept(ipv6.ICMPTypeEchoReply)
	// 	if err := c.IPv6PacketConn().SetICMPFilter(&f); err != nil {
	// 		return err
	// 	}
	// }

	wm := icmp.Message{
		Type: tt.mtype, Code: 0,
		Body: &icmp.Echo{
			ID: os.Getpid() & 0xffff, Seq: 1 << uint(seq),
			Data: []byte("HELLO-R-U-THERE"),
		},
	}

	//返回校验和消息
	wb, err := wm.Marshal(nil)
	if err != nil {
		return err
	}
	fmt.Println(wm.Body)
	if n, err := c.WriteTo(wb, dst); err != nil {
		return err
	} else if n != len(wb) {
		return fmt.Errorf("got %v; want %v", n, len(wb))
	}

	rb := make([]byte, 1500)
	if err := c.SetReadDeadline(time.Now().Add(3 * time.Second)); err != nil {
		return err
	}

	n, peer, err := c.ReadFrom(rb)
	if err != nil {
		return err
	}

	rm, err := icmp.ParseMessage(tt.protocol, rb[:n])
	if err != nil {
		return err
	}

	fmt.Printf("read from %v\n", peer)
	//获取我们发送出去的数据 在这里可以写入时间来计算往返时间
	dd, err := rm.Body.Marshal(iana.ProtocolICMP)
	//fmt.Println(icmp.Echo(rm.Body))
	fmt.Printf("%s\n", dd[4:])
	/**
	* 或者使用unsafe 转换
	*	t := T {
	        A: 10,
	        B: "abc",
	    }
	    l := unsafe.Sizeof(t)
	    pb := (*[1024]byte)(unsafe.Pointer(&t))
	    fmt.Println("Struct:", t)
	    fmt.Println("Bytes:", (*pb)[:l])
	 *
	 *
	 *
	*/
	switch rm.Type {
	case ipv4.ICMPTypeEchoReply, ipv6.ICMPTypeEchoReply:
		return nil
	default:
		return fmt.Errorf("got %+v from %v; want echo reply", rm, peer)
	}
}

// func TestConcurrentNonPrivilegedListenPacket(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("avoid external network")
// 	}
// 	switch runtime.GOOS {
// 	case "darwin":
// 	case "linux":
// 		t.Log("you may need to adjust the net.ipv4.ping_group_range kernel state")
// 	default:
// 		t.Skipf("not supported on %s", runtime.GOOS)
// 	}

// 	network, address := "udp4", "127.0.0.1"
// 	if !nettest.SupportsIPv4() {
// 		network, address = "udp6", "::1"
// 	}
// 	const N = 1000
// 	var wg sync.WaitGroup
// 	wg.Add(N)
// 	for i := 0; i < N; i++ {
// 		go func() {
// 			defer wg.Done()
// 			c, err := icmp.ListenPacket(network, address)
// 			if err != nil {
// 				t.Error(err)
// 				return
// 			}
// 			c.Close()
// 		}()
// 	}
// 	wg.Wait()
// }

// var nonPrivilegedPingTests = []pingTest{
// 	{"udp4", "0.0.0.0", iana.ProtocolICMP, ipv4.ICMPTypeEcho},

// 	{"udp6", "::", iana.ProtocolIPv6ICMP, ipv6.ICMPTypeEchoRequest},
// }

// func TestNonPrivilegedPing(t *testing.T) {
// 	if testing.Short() {
// 		t.Skip("avoid external network")
// 	}
// 	switch runtime.GOOS {
// 	case "darwin":
// 	case "linux":
// 		t.Log("you may need to adjust the net.ipv4.ping_group_range kernel state")
// 	default:
// 		t.Skipf("not supported on %s", runtime.GOOS)
// 	}

// 	for i, tt := range nonPrivilegedPingTests {
// 		if err := doPing(tt, i); err != nil {
// 			t.Error(err)
// 		}
// 	}
// }
