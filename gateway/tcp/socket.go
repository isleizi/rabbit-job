package tcp

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net"
	"os"
	"time"
)

type socket struct {
	// 超时时间
	timeout int
	// 地址:端口
	addr string
}

func (s *socket) run() {
	listener, err := net.Listen("tcp", s.addr)
	s.checkError(err)
	defer listener.Close()

	group := gin.RouterGroup{}
	group.Use()

	//开启多个协程
	for {
		//用conn接收链接
		conn, err := listener.Accept()
		if err != nil {
			gin.Logger()
			logg.Error(err)
			continue
		}
		go s.handleConnection(conn, s.timeout)
	}
}

//长连接入口
func (s *socket) handleConnection(conn net.Conn, timeout int) {
	conn.Write([]byte("socket服务连接成功的消息..."))
	buffer := make([]byte, 112048)
	for {
		n, err := conn.Read(buffer)
		if err != nil {
			log.Println(conn.RemoteAddr().String(), " connection error: ", err)
			conn.Close()
			return
		}

		// 接收数据
		data := buffer[:n]
		msg := make(chan byte)

		//心跳计时
		go s.heartBeating(conn, msg, timeout)

		//检测每次Client是否有数据传来
		go s.gravelChannel(data, msg)

		log.Println("receive data length:", n)
		log.Println(conn.RemoteAddr().String(), "receive data string:", string(data))
	}
}

// 心跳计时，根据GravelChannel判断Client是否在设定时间内发来信息
func (s *socket) heartBeating(conn net.Conn, readerChannel chan byte, timeout int) {
	select {
	case fk := <-readerChannel:
		Log(conn.RemoteAddr().String(), "receive data string:", string(fk))
		Log("心跳重新时间")
		conn.SetDeadline(time.Now().Add(time.Duration(timeout) * time.Second))
		conn.Write([]byte("socket服务心跳成功的消息2222..."))
		break
	case <-time.After(time.Second * 5):
		Log("It's really weird to get Nothing!!!")
		conn.Close()
	}
}

// 写入通道
func (s *socket) gravelChannel(n []byte, mess chan byte) {
	for _, v := range n {
		mess <- v
	}
	close(mess)
}

// IO异常
func (s *socket) checkError(err error) {
	if err != nil {
		fmt.Fprintf(os.Stderr, "Fatal error: %s", err.Error())
		os.Exit(1)
	}
}
