package nginx_syslog

import (
	"net"
	"fmt"
	"os"
	"os/signal"
	"syscall"
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

func signalHandler(udp *net.UDPConn) {
	signalChan := make(chan os.Signal)

	// 监听指定信号
	signal.Notify(
		signalChan,
		//syscall.SIGHUP,	//终端断线
		syscall.SIGINT,  //Ctrl+C信号
		syscall.SIGTERM, //结束程序(可以被捕获、阻塞或忽略)
		//syscall.SIGUSR2,//同SIGUSR1，保留给用户使用的信号
	)

	// 输出当前进程的pid
	fmt.Println("pid is: ", os.Getpid())

	go func(udp *net.UDPConn) {
		// 处理信号
		switch <-signalChan {
		/*		case syscall.SIGUSR2:
					fmt.Println("USER2信号")
					startNewProcess()*/
		case syscall.SIGINT: //
			fmt.Println("Ctrl+C信号")
			//file := os.NewFile(3, "")
			file, _ := udp.File()
			startNewProcess(file)

		case syscall.SIGTERM:
			fmt.Println(syscall.SIGTERM)
		}
	}(udp)

	fmt.Println(os.Args[0])

}

func startNewProcess(file *os.File) {
	procAttr := &syscall.ProcAttr{
		Env:   os.Environ(),
		Files: []uintptr{os.Stdin.Fd(), os.Stdout.Fd(), os.Stderr.Fd(), file.Fd()},
	}
	pid, _, _ := syscall.StartProcess(os.Args[0], os.Args, procAttr)
	fmt.Printf("start new process ... pid: %d", pid)
}
