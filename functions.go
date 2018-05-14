package nginx_syslog

import (
	"net"
	"fmt"
	"os"
)

var (
	NETWORK string = "udp"
	ADDRESS string = ":514"
	//用UDP协议发送时，用sendto函数最大能发送数据的长度为：65535－20－8＝65507字节，其中20字节为IP包头长度，8字节为UDP包头长度。用sendto函数发送数据时，如果指的的数据长度大于该值，则函数会返回错误。
	MAXDATALENGTH int = 65507
	GONUM         int = 50
)

func recriver() *net.UDPConn {
	udpAddress, err := net.ResolveUDPAddr(NETWORK, ADDRESS)
	ln, err := net.ListenUDP(udpAddress.Network(), udpAddress)
	checkError(err)
	fmt.Println("listening " + NETWORK + " port " + ADDRESS)
	return ln
}

func checkError(err error) {
	if err != nil {
		fmt.Println("Error: %s", err.Error())
		os.Exit(1)
	}
}

func recvUDPMsg(conn *net.UDPConn) string {
	buf := make([]byte, MAXDATALENGTH)
	n, _, err := conn.ReadFromUDP(buf[0:])
	checkError(err)
	return string(buf[0:n])
}
