package nginx_syslog

import (
	"net"
	"fmt"
	"net/url"
	"log"
)

type Parser struct {
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
	signalHandler(NetUDP.Conn)

	for {
		c := make(chan struct{}, GONUM)
		go func(c chan struct{}) {
			defer func() {
				c <- struct{}{}
				if err := recover(); err != nil {
					// 这里可以对异常进行一些处理和捕获
					log.Print(err)
				}
			}()

			data := recvUDPMsg(NetUDP.Conn)
			Rfc3164 := NewRfc3164(data)
			Log := NewLog(Rfc3164.Content)
			Request := NewRequest(Log.Request)
			m, err := url.Parse(Request.Path)
			if err != nil {
				panic(err)
			}
			query := m.Query()
			fmt.Println(query)
		}(c)
		<-c
	}
	defer NetUDP.Conn.Close()
}
