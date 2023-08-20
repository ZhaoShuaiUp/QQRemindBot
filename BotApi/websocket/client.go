package websocket

import (
	"RemindBot/BotApi/dto"
	"RemindBot/BotApi/event"
	"RemindBot/BotApi/token"
	"encoding/json"
	"fmt"
	"os"
	"time"

	wss "github.com/gorilla/websocket"
)

// DefaultQueueSize 监听队列的缓冲长度
const DefaultQueueSize = 10000

type messageChan chan *dto.WSPayload
type closeErrorChan chan error

// 会话信息
type Session struct {
	ID      string
	URL     string
	Token   token.Token
	Intent  dto.Intent
	LastSeq uint32
}

type Client struct {
	version         int
	conn            *wss.Conn
	session         *Session
	heartBeatTicker *time.Ticker
	messageQueue    messageChan
	closeQueue      closeErrorChan
	user            struct {
		ID       string
		Username string
		Bot      bool
	}
}

// 创建一个 Client 实例
func CreateClient(apInfo *dto.WebsocketAP, token *token.Token, intents *dto.Intent) *Client {
	return &Client{
		messageQueue:    make(messageChan, DefaultQueueSize),
		closeQueue:      make(closeErrorChan, 10),
		heartBeatTicker: time.NewTicker(60 * time.Second),
		session: &Session{
			URL:     apInfo.URL,
			Token:   *token,
			Intent:  *intents,
			LastSeq: 0,
		},
	}
}

// 启动 Client 实例
func (c *Client) Start() {
	// 连接 websocket
	if err := c.Connect(); err != nil {
		fmt.Printf("Websocket start failed. session: %v", c.session)
		os.Exit(1)
		return
	}

	// 开始鉴权
	if err := c.Identify(); err != nil {
		fmt.Printf("Websocket Identify failed. session: %v", c.session)
		os.Exit(1)
		return
	}

	// 开始监听
	if err := c.Listening(); err != nil {
		fmt.Printf("Websocket Listening failed. session: %v", c.session)
		os.Exit(1)
		return
	}
}

// 连接到 websocket
func (c *Client) Connect() error {
	if c.session.URL == "" {
		return fmt.Errorf("websocket url is empty")
	}
	var err error
	c.conn, _, err = wss.DefaultDialer.Dial(c.session.URL, nil)
	if err != nil {
		return err
	}
	return nil
}

// 进行鉴权
func (c *Client) Identify() error {
	payload := &dto.WSPayload{
		Data: &dto.WSIdentityData{
			Token:   c.session.Token.GetAuthorization(),
			Intents: c.session.Intent,
			Shard:   []uint32{0, 1},
		},
	}
	payload.OPCode = dto.OPCODE_IDENTIFY
	return c.Write(payload)
}

// 开始监听 websocket
func (c *Client) Listening() error {
	defer c.Close()
	go c.readMessageToQueue()
	go c.listenMessageAndHandle()

	// 处理 client 上的事件
	for {
		select {
		case err := <-c.closeQueue:
			fmt.Printf("%v Listening stop. err is %v\n", c.session, err)
			return err
		case <-c.heartBeatTicker.C:
			fmt.Printf("%v listened heartBeat", c.session)
			heartBeatEvent := &dto.WSPayload{
				OPCode: dto.OPCODE_HEARTBEAT,
				Data:   c.session.LastSeq,
			}

			_ = c.Write(heartBeatEvent)
		}
	}
}

// 写入 websocket 数据
func (c *Client) Write(payload *dto.WSPayload) error {
	m, _ := json.Marshal(payload)
	fmt.Printf("[Write] %s\n", m)
	if err := c.conn.WriteMessage(wss.TextMessage, m); err != nil {
		fmt.Printf("%v WriteMessage failed, %v", c.session, err)
		c.closeQueue <- err
		return err
	}
	return nil
}

// 关闭 websocket
func (c *Client) Close() {
	if err := c.conn.Close(); err != nil {
		fmt.Printf("close websocket error: %v\n", err)
	}
	c.heartBeatTicker.Stop()
}

// 从 websocket 上获取消息并保存到管道中，生产者
func (c *Client) readMessageToQueue() {
	for {
		_, message, err := c.conn.ReadMessage()
		if err != nil {
			fmt.Printf("%v read message failed, %v, message %s", c.session, err, string(message))
			return
		}
		payload := &dto.WSPayload{}
		if err := json.Unmarshal(message, payload); err != nil {
			fmt.Printf("%v json failed, %v", c.session, err)
			continue
		}
		payload.RawMessage = message
		if c.isHandleBuildIn(payload) {
			continue
		}
		fmt.Printf("[Receive] %v receive %v message, %s\n", c.session, payload.OPCode, string(message))
		c.messageQueue <- payload
	}
}

// 从管道中读取消息，并处理，消费者
func (c *Client) listenMessageAndHandle() {
	for payload := range c.messageQueue {
		c.session.LastSeq = payload.Seq
		if payload.Type == "READY" {
			c.readyHandler(payload)
			continue
		}
		if err := event.ParseAndHandle(payload); err != nil {
			fmt.Printf("%v ParseAndHandle failed, %v", c.session, err)
		}
	}
}

// 处理非自定义事件
func (c *Client) isHandleBuildIn(payload *dto.WSPayload) bool {
	switch payload.OPCode {
	case dto.OPCODE_HELLO:
		c.startHeartBeatTicker(payload.RawMessage)
	case dto.OPCODE_HEARTBEAT_ACK:
	case dto.OPCODE_RECONNECT:

	case dto.OPCODE_INVALID_SESSION:

	default:
		return false
	}
	return true
}

// startHeartBeatTicker 启动定时心跳
func (c *Client) startHeartBeatTicker(message []byte) {
	helloData := &dto.WSHelloData{}
	if err := event.ParseData(message, helloData); err != nil {
		fmt.Printf("%v hello data parse failed, %v, message %v", c.session, err, message)
	}
	// 重新心跳定时器
	c.heartBeatTicker.Reset(time.Duration(helloData.HeartbeatInterval) * time.Millisecond)
}

// 针对 ready 报文进行处理
func (c *Client) readyHandler(payload *dto.WSPayload) {
	readyData := &dto.WSReadyData{}
	if err := event.ParseData(payload.RawMessage, readyData); err != nil {
		fmt.Printf("%v parseReadyData failed, %v, message %v", c.session, err, payload.RawMessage)
	}
	c.version = readyData.Version
	c.session.ID = readyData.SessionID
	c.user.Bot = readyData.User.Bot
	c.user.Username = readyData.User.Username
	c.user.ID = readyData.User.ID
}
