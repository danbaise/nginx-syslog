package nginx_syslog

import (
	"fmt"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"
)

var (
	NETWORK string = "udp"
	ADDRESS string = ":8888"
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

func signalHandler() {
	signalChan := make(chan os.Signal)

	// 监听指定信号
	signal.Notify(
		signalChan,
		//syscall.SIGHUP,	//终端断线
		syscall.SIGINT,  //Ctrl+C信号
		syscall.SIGTERM, //结束程序(可以被捕获、阻塞或忽略)
		syscall.SIGUSR2, //同SIGUSR1，保留给用户使用的信号
	)

	// 输出当前进程的pid
	log.Println("Pid is:", os.Getpid())

	// 处理信号
	for {
		switch <-signalChan {
		case syscall.SIGUSR2: //
			fmt.Println("Received SIGUSR2.")
			startNewProcess()

		case syscall.SIGTERM:
			log.Println(os.Getpid(), "Received SIGTERM.")
			shutdown()
		}
	}

}

func startNewProcess() {
	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd()},
	}
	var args []string
	for _, arg := range os.Args {
		if arg == "-graceful" {
			break
		}
		args = append(args, arg)
	}
	args = append(args, "-graceful")
	pid, _, _ := syscall.StartProcess(os.Args[0], args, procAttr)
	fmt.Printf("Start new process ... Pid: %d \r\n", pid)

}

func shutdown() {
	Conn.Close()
	log.Println(os.Getpid(), "Shutdown.")
	os.Exit(0)
}
