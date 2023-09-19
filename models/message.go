package models

import (
	"encoding/json"
	"fmt"
	"github.com/fasthttp/websocket"
	"github.com/gofiber/fiber/v2"
	"github.com/valyala/fasthttp"
	"gopkg.in/fatih/set.v0"
	"gorm.io/gorm"
	"log"
	"strconv"
	"sync"
)

type Message struct {
	gorm.Model
	FormId   uint   `json:"formId"`   //发送者
	TargetId uint   `json:"targetId"` //接收者
	Type     string `json:"type"`     //消息类型 群聊 私聊 广播
	Media    string `json:"media"`    //消息类型 文字 图片 音频
	Content  string `json:"content"`  //消息内容
	Pic      string `json:"pic"`      //图片
	Url      string `json:"url"`      //链接
	Desc     string `json:"desc"`     //描述
	Amount   int    `json:"amount"`   //其他数字统计
}

// TableName 设置表名
func (table *Message) TableName() string {
	return "message"
}

type Node struct {
	Conn      *websocket.Conn `json:"conn"`
	DataQueue chan []byte     `json:"dataQueue"`
	GroupSets set.Interface   `json:"groupSets"`
	CloseChan chan bool       `json:"closeChan"` // 增加了安全关机功能
}

// Chat
var clientMap map[int64]*Node

// Chat
var rwLocker sync.RWMutex

func Chat(c *fiber.Ctx) error {
	userId, err := strconv.ParseInt(c.Query("id"), 10, 64) // 您应该实现一个函数来从请求中获取userId
	if err != nil {
		return err
	}

	isValid := true // 您应该实现一个函数来检查有效性

	upgrader := websocket.FastHTTPUpgrader{
		CheckOrigin: func(ctx *fasthttp.RequestCtx) bool {
			return isValid
		},
	}

	err = upgrader.Upgrade(c.Context(), func(conn *websocket.Conn) {
		log.Println("WebSocket连接!")
		defer func() {
			rwLocker.Lock()           // 加写锁
			delete(clientMap, userId) // 删除节点
			rwLocker.Unlock()         // 解写锁
			conn.Close()              // 关闭连接
		}()

		node := &Node{ // 创建节点
			Conn:      conn,                    // 连接
			DataQueue: make(chan []byte, 50),   // 数据队列
			GroupSets: set.New(set.ThreadSafe), // 群组集合
			CloseChan: make(chan bool),         // 关闭通道
		}

		rwLocker.Lock()          // 加写锁
		clientMap[userId] = node // 添加节点
		rwLocker.Unlock()        // 解写锁

		go sendProc(node) // 发送协程
		go recvProc(node) // 接收协程

		sendMsg(userId, []byte("欢迎来到聊天系统")) // 发送欢迎消息

		<-node.CloseChan // 阻塞
	})

	return err
}

func sendProc(node *Node) {
	for { // 不停地从通道中读取数据
		select { // 从多个通道中选择
		case data := <-node.DataQueue: // 从数据队列中读取数据
			err := node.Conn.WriteMessage(websocket.TextMessage, data) // 写入 WebSocket
			if err != nil {
				log.Println("写错误信息:", err) // 打印错误信息
				node.CloseChan <- true     // 发送关闭通道
				return
			}
		case <-node.CloseChan: // 从关闭通道中读取数据
			return
		}
	}
}

func recvProc(node *Node) {
	for { // 不停地从 WebSocket 中读取数据
		_, data, err := node.Conn.ReadMessage() // 读取 WebSocket
		if err != nil {
			log.Println("读取错误信息:", err) // 打印错误信息
			node.CloseChan <- true      // 发送关闭通道
			return
		}
		dispatch(data) // 调度消息
	}
}

func dispatch(data []byte) {
	msg := Message{}                  // 创建消息结构体
	fmt.Println(string(data))         // 打印消息
	err := json.Unmarshal(data, &msg) // 解组消息
	if err != nil {
		log.Printf("解组消息出错: %v. Data: %s\n", err, string(data)) // 打印错误信息
		return
	}
	switch msg.Type { // 根据消息类型进行调度
	case "1":
		sendMsg(int64(msg.TargetId), data) // 私聊
	}
}

func sendMsg(userId int64, msg []byte) {
	rwLocker.RLock()              // 加读锁
	node, ok := clientMap[userId] // 从 clientMap 中获取节点
	rwLocker.RUnlock()            // 解读锁
	if ok {
		node.DataQueue <- msg // 写入数据队列
	}
}

func init() {
	clientMap = make(map[int64]*Node) // 初始化 clientMap
}
