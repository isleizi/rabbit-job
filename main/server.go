package main

import (
	"crypto/md5"
	"encoding/hex"
	"fmt"
	"io"
	"log"
	"net"
	"net/http"
	"os"
	"time"
)

var connList = make(map[string]*net.Conn)

func main() {
	// 表示监听本地所有ip的8080端口，也可以这样写：addr := ":8080"
	addr := "0.0.0.0:8889"
	listener, err := net.Listen("tcp", addr)
	CheckError(err)
	defer listener.Close()

	go webServer()

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

// hello world, the web server
// w: 给客户端回复数据， req: 读取客户端发送的数据
func HelloServer(w http.ResponseWriter, req *http.Request) {
	// 打印客户端头信息
	fmt.Println(req.Method)
	fmt.Println(req.Header)
	fmt.Println(req.Body)
	fmt.Println(req.URL)
	fmt.Println(req.Form.Get("code"))

	req.ParseForm()
	d := req.Form
	fmt.Println(d)
	fmt.Println(req.PostForm)

	//fmt.Println(req.Form.Get("code"))
	//if len(req.Form) > 0 {
	//	for k,v := range req.Form {
	//		fmt.Printf("%s=%s\n", k, v[0])
	//	}
	//}

	//aa
	fmt.Println(req.URL.Query().Get("name"))

	// 给客户端回复数据
	io.WriteString(w, "hello, world!\n")
	w.Write([]byte("lisa"))
}
func webServer() {
	// 注册函数，用户连接， 自动调用指定处理函数
	http.HandleFunc("/hello", HelloServer)

	// 监听绑定
	err := http.ListenAndServe(":12345", nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
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
func GravelChannel(n []byte, mess chan byte) {
	for _, v := range n {
		mess <- v
	}
	close(mess)
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
