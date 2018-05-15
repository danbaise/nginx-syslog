package nginx_syslog

import (
	"net"
	"fmt"
	"net/url"
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
			Request := NewRequest(Log.Request)
			m, err := url.Parse(Request.Path)
			if err != nil {
				checkError(err)
			}
			fmt.Println(m.Query())
			c <- struct{}{}
		}(c)
		<-c
	}
	defer NetUDP.Conn.Close()
}
