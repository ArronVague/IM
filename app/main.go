package main

import (
	_ "IM/app/model"
	"fmt"
	"github.com/gorilla/websocket"
	"html/template"
	"log"
	"net/http"
	"time"
)

func registerView() {
	tpl, err := template.ParseGlob("./app/view/**/*")
	if err != nil {
		log.Fatal(err.Error())
	}
	for _, v := range tpl.Templates() {
		tplName := v.Name()
		http.HandleFunc(tplName, func(writer http.ResponseWriter, request *http.Request) {
			tpl.ExecuteTemplate(writer, tplName, nil)
		})
	}
}

var (
	//完成握手操作
	upgrade = websocket.Upgrader{
		//允许跨域(一般来讲,websocket都是独立部署的)
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}
)

func wsHandler(w http.ResponseWriter, r *http.Request) {
	var (
		conn *websocket.Conn
		err  error
		data []byte
	)
	//服务端对客户端的http请求(升级为websocket协议)进行应答，应答之后，协议升级为websocket，http建立连接时的tcp三次握手将保持。
	if conn, err = upgrade.Upgrade(w, r, nil); err != nil {
		return
	}

	//启动一个协程，每隔1s向客户端发送一次心跳消息
	go func() {
		var (
			err error
		)
		for {
			if err = conn.WriteMessage(websocket.TextMessage, []byte("heartbeat")); err != nil {
				return
			}
			time.Sleep(1 * time.Second)
		}
	}()
	//得到websocket的长链接之后,就可以对客户端传递的数据进行操作了
	for {
		//通过websocket长链接读到的数据可以是text文本数据，也可以是二进制Binary
		if _, data, err = conn.ReadMessage(); err != nil {
			goto ERR
		}
		if err = conn.WriteMessage(websocket.TextMessage, data); err != nil {
			goto ERR
		}
	}
ERR:
	//出错之后，关闭socket连接
	conn.Close()
}

func main() {
	fmt.Println("websocket server start")
	//http.HandleFunc 是一个方便的函数，它创建一个默认的 http.ServeMux（HTTP 请求的多路复用器），并且将给定的函数注册到这个多路复用器上。这个函数需要有特定的签名，即 func(ResponseWriter, *Request)。
	http.HandleFunc("/ws", wsHandler)
	http.Handle("/asset/", http.FileServer(http.Dir(".")))
	http.Handle("/resource/", http.FileServer(http.Dir(".")))

	//http.ListenAndServe 函数接受两个参数：一个是服务器的地址，另一个是处理请求的 http.Handler。如果这个 http.Handler 是 nil，那么就会使用默认的 http.ServeMux，也就是我们在 http.HandleFunc 中注册的那个。
	log.Fatal(http.ListenAndServe(":1060", nil))

}
