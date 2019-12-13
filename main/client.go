package main

import (
	"log"
	"net"
)

func establishConn(i int) net.Conn {
	conn, err := net.Dial("tcp", ":8888")
	if err != nil {
		log.Printf("%d: dial error: %s", i, err)
		return nil
	}
	log.Println(i, ":connect to server ok")
	return conn
}

func main() {
	var ch = make(chan string, 10)
	ch <- "d4df6c21557d5900fc53a169179a862f"

	//var sl []net.Conn
	//for i := 1; i < 1000; i++ {
	//	conn := establishConn(i)
	//	if conn != nil {
	//		sl = append(sl, conn)
	//	}
	//}
	//
	//time.Sleep(time.Second * 10000)
}
