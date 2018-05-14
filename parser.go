package nginx_syslog

import (
	"net"
	"fmt"
)

type Parser struct {
	*Rfc3164
	*Log
	UDPConn *net.UDPConn
}

func NewParser() *Parser {
	return &Parser{UDPConn: recriver()}
}

func (p *Parser) Handle() {
	for {
		c := make(chan struct{}, GONUM)
		for i := 0; i < GONUM; i++ {
			go func(c chan struct{}) {
				//		fmt.Println(recvUDPMsg(p.UDPConn))
				p.Rfc3164 = NewRfc3164(recvUDPMsg(p.UDPConn))
				p.Log = NewLog(p.Rfc3164.Content)
				fmt.Printf("%#v\r\n%#v\r\n", p.Rfc3164, p.Log)
				c <- struct{}{}
			}(c)
		}
		<-c
	}
	defer p.end()
}

func (p *Parser) end() {
	p.UDPConn.Close()
}
