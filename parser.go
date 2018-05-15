package nginx_syslog

import (
	"net"
	"fmt"
)

type Parser struct {
	/*	*Rfc3164
		*Log*/
}

type NetUDP struct {
	Conn *net.UDPConn
}

func NewParser() *Parser {
	return &Parser{}
}

func NewNetUDP() *NetUDP {
	return &NetUDP{Conn: recriver()}
}

func (p *Parser) Handle() {
	NetUDP := NewNetUDP()
	for {
		c := make(chan struct{}, GONUM)
		go func(c chan struct{}) {
			data := recvUDPMsg(NetUDP.Conn)
			Rfc3164 := NewRfc3164(data)
			Log := NewLog(Rfc3164.Content)
			fmt.Printf("%#v\r\n%#v\r\n", Rfc3164, Log)
			c <- struct{}{}
		}(c)
		<-c
	}
	defer NetUDP.Conn.Close()
}
