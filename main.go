package main

import (
	"flag"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

var (
	port int

	username string
	password string
)

func init() {
	flag.StringVar(&username, "user", "", "username")
	flag.StringVar(&password, "pwd", "", "password")
	flag.IntVar(&port, "p", 1080, "port on listen, must be greater than 0")
	flag.Parse()
}

func main() {
	if port <= 0 {
		flag.Usage()
		os.Exit(1)
	}
	var serverAddr *net.TCPAddr
	if addr, err := net.ResolveTCPAddr("tcp", ":"+strconv.Itoa(port)); err != nil {
		fmt.Println(err.Error())
		os.Exit(1)
	} else {
		serverAddr = addr
	}
	
	currentTime := time.Now().Format("2006/01/02 15:04:05")

	server := NewServer()
	server.EnableUDP()
	server.OnStarted(func(listener *net.TCPListener) {
		fmt.Printf("MAIN    : %v: INFO     SOCKS5 Proxy started on port: %s\n", currentTime, listener.Addr().String())
	})
	server.OnConnected(func(network, address string, port int) {
		fmt.Printf("SOCKS5  : %v: INFO     [%s] connect to: %s:%d\n", currentTime, network, address, port)
	})
	if username != "" || password != "" {
		server.SetAuthHandle(handlerAuth)
	}
	if err := server.Run(serverAddr); err != nil {
		fmt.Println("MAIN    : %v: ERROR    Running SOCKS5 Proxy server error: %s\n", currentTime, err.Error())
		os.Exit(1)
	}

	fmt.Println("MAIN    : %v: INFO     SOCKS5 server normal exit.\n", currentTime)
	os.Exit(0)
}

func handlerAuth(u, p string) bool {
	return u == username && p == password
}
