package nginx_syslog

import (
	"flag"
	"fmt"
	"log"
	"net"
	"net/url"
	"syscall"
	"time"
)

var Conn *net.UDPConn
var isChild bool

type Parser struct {
}

func NewParser() *Parser {
	return &Parser{}
}

func (p *Parser) Handle() {

	flag.BoolVar(&isChild, "graceful", false, "listen on open fd (after forking)")
	flag.Parse()
	if isChild {
		syscall.Kill(syscall.Getppid(), syscall.SIGTERM) //干掉父进程程序结束(terminate)信号
		time.Sleep(time.Second * 1)
	}

	Conn = recriver()
	go signalHandler()

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

			data := recvUDPMsg(Conn)
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

	defer shutdown()
}
