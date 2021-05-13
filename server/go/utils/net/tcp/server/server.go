package main

import (
	"bufio"
	"fmt"
	"io"
	"net"
	"time"
)

func Server(addr string) {

	var tcpAddr *net.TCPAddr

	tcpAddr, _ = net.ResolveTCPAddr("tcp", addr)

	tcpListener, _ := net.ListenTCP("tcp", tcpAddr)

	defer tcpListener.Close()

	fmt.Println("Server ready to read ...")
	for {
		tcpConn, err := tcpListener.AcceptTCP()
		if err != nil {
			fmt.Println("accept error:", err)
			continue
		}
		fmt.Println("A client connected : " + tcpConn.RemoteAddr().String())
		go tcpPipe(tcpConn)
	}

}

func tcpPipe(conn *net.TCPConn) {
	ipStr := conn.RemoteAddr().String()

	defer func() {
		fmt.Println(" Disconnected : " + ipStr)
		conn.Close()
	}()

	reader := bufio.NewReader(conn)
	i := 0

	for {
		message, err := reader.ReadString('\n') //将数据按照换行符进行读取。
		if err != nil || err == io.EOF {
			break
		}

		fmt.Println(string(message))

		time.Sleep(time.Second * 3)

		msg := time.Now().String() + conn.RemoteAddr().String() + " Server Say hello! \n"
		b := []byte(msg)

		conn.Write(b)
		i++

		if i > 10 {
			break
		}
	}
}
