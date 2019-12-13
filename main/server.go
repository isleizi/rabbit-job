package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"log"
	"net"
	"os"
	"time"
)

var connList = make(map[string]*net.Conn)
var ipp string

func main() {
	// 表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"
	addr := "0.0.0.0:8888"
	listener, err := net.Listen("tcp", addr)
	CheckError(err)
	defer listener.Close()

	//开启多个协程
	for {
		//用conn接收链接
		conn, err := listener.Accept()
		if err != nil {
			log.Fatal(err)
			continue
		}
		go handleConnection(conn, 10)
	}
}

//长连接入口
func handleConnection(conn net.Conn, timeout int) {
	conn.Write([]byte("socket服务连接成功的消息..."))
	buffer := make([]byte, 112048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			Log(conn.RemoteAddr().String(), " connection error: ", err)
			conn.Close()
			return
		}

		client := md5V(conn.RemoteAddr().String())
		connList[client] = &conn

		data := buffer[:n]
		msg := make(chan byte)

		//心跳计时
		go HeartBeating(conn, msg, timeout)

		//检测每次Client是否有数据传来
		go GravelChannel(data, msg)

		Log("receive data length:", n)
		Log(conn.RemoteAddr().String(), "receive data string:", string(data))
	}
}

// 心跳计时，根据GravelChannel判断Client是否在设定时间内发来信息
func HeartBeating(conn net.Conn, readerChannel chan byte, timeout int) {
	select {
	case fk := <-readerChannel:
		Log(conn.RemoteAddr().String(), "receive data string:", string(fk))
		Log("注册时间")
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		conn.Write([]byte("socket服务心跳成功的消息2222..."))
		break
	case <-time.After(time.Second * 5):
		Log("It's really weird to get Nothing!!!")
		conn.Close()
	}
}

// 写入通道
func GravelChannel(n []byte, mess chan byte) {
	for _, v := range n {
		mess <- v
	}
	close(mess)
}

func LoopClient() {
	for range time.Tick(time.Millisecond * 3000) {
		fmt.Println("dddd")
		for k := range connList {
			fmt.Println("dddd2222")
			if k != "" {
				fmt.Println("dddd33333")
				client := connList[k]
				(*client).Write([]byte("socket服务端发送心跳成功1111..."))
			} else {
				fmt.Println("没有客户端")
			}
		}
	}
}

// 日志
func Log(v ...interface{}) {
	log.Println(v...)
}

// IO异常
func CheckError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}

// md5
func md5V(str string) string {
	h := md5.New()
	h.Write([]byte(str))
	return hex.EncodeToString(h.Sum(nil))
}
