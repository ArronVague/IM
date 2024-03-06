package controller

import (
	"fmt"
	"github.com/gorilla/websocket"
	"gopkg.in/fatih/set.v0"
	"log"
	"net/http"
	"strconv"
	"sync"
)

// Node 本核心在于形成userid和Node的映射关系
type Node struct {
	Conn *websocket.Conn
	//并行转串行
	DataQueue chan []byte
	GroupSets set.Interface
}

// userid和Node映射关系表
var clientMap = make(map[int64]*Node)

// sync.RWMutex 是 Go 语言中的一种读写锁，它允许多个 goroutine 同时读取某个资源，但是在写入（更新）资源时，只允许一个 goroutine 访问，且在写入时，其他任何读取或写入操作都会被阻塞。适合读操作多于写操作。
var rwLocker sync.RWMutex

func Chat(writer http.ResponseWriter, request *http.Request) {
	query := request.URL.Query()
	id := query.Get("id")
	token := query.Get("token")
	userId, _ := strconv.ParseInt(id, 10, 64)
	isLegal := checkToken(userId, token)

	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return isLegal
		},
	}).Upgrade(writer, request, nil)

	if err != nil {
		log.Println(err.Error())
		return
	}

	node := &Node{
		Conn:      conn,
		DataQueue: make(chan []byte, 50),
		GroupSets: set.New(set.ThreadSafe),
	}

	comIds := contactService.SearchCommunityIds(userId)
	for _, v := range comIds {
		node.GroupSets.Add(v)
	}

	//写锁
	rwLocker.Lock()
	clientMap[userId] = node
	rwLocker.Unlock()

	go sendproc(node)

	go recvproc(node)
}

// 在 Web 开发中，使用 token 进行身份验证是一种常见的安全策略。这是因为仅仅依赖用户 ID 来进行身份验证是不安全的，因为用户 ID 通常是公开的，任何人都可以知道，如果只需要用户 ID 就能进行操作，那么任何人都可以伪装成任何用户。
//
// Token 是服务器为每一个已认证的用户生成的一个唯一的、难以猜测的字符串。当用户进行登录操作时，服务器会生成一个 token 并返回给用户，用户在后续的请求中需要带上这个 token 来证明自己的身份。服务器收到请求后，会检查 token 是否有效（例如检查 token 是否存在，是否过期，是否和用户 ID 匹配等），只有当 token 有效时，服务器才会处理用户的请求。
//
// 这种方式的优点是，即使其他人知道了用户的 ID，但是如果没有有效的 token，他们仍然无法伪装成用户。此外，服务器可以随时使一个 token 失效（例如用户登出时），从而结束用户的会话。
// 校验token是否合法
func checkToken(userId int64, token string) bool {
	user := UserService.Find(userId)
	return user.Token == token
}

// 发送逻辑
func sendproc(node *Node) {
	for {
		select {
		case data := <-node.DataQueue:
			err := node.Conn.WriteMessage(websocket.TextMessage, data)
			if err != nil {
				log.Println(err.Error())
				return
			}
		}
	}
}

// 接收逻辑
func recvproc(node *Node) {
	for {
		_, data, err := node.Conn.ReadMessage()
		if err != nil {
			log.Println(err.Error())
			return
		}

		dispatch(data)
		//todo对data进一步处理
		fmt.Printf("recv<=%s", data)
	}
}
