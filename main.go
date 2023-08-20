package main

import (
	"RemindBot/BotApi"
	"RemindBot/BotApi/dto"
	"RemindBot/BotApi/event"
	"RemindBot/BotApi/openapi"
	"RemindBot/BotApi/token"
	"RemindBot/BotApi/websocket"
	"RemindBot/DB"
	"RemindBot/Remind"
	"context"
	"fmt"
	"github.com/robfig/cron/v3"
	"os"
	"time"
)

var api *openapi.OpenAPI
var ctx context.Context
var sendTemporary *Remind.SendTemporary
var sendOnce Remind.SendService
var sendDay Remind.SendService

func aTMessageEventHandler(event *dto.WSPayload, data *dto.WSATMessageData) error {
	// 去除用户 AT 信息
	fmt.Printf("[Listening Receive] message: %v\n", data)
	remindMessage := dto.ETLMessage(data.Content)
	// 抽取出remind信息, 插入数据库
	remind := Remind.ExtractMessage(remindMessage)
	if remind != nil {
		remind.Uid = data.Author.ID
		remind.Cid = data.ChannelID
		if remind.NeedSendTemporary() {
			sendTemporary.AddRemind(*remind)
		}
	}
	//remind.Insert(DB.GetDBInstance())

	var resp string
	if remind == nil {
		resp = "请您确认提醒格式是否正确！"
	} else {
		resp = "已为您设置好提醒服务"
	}

	// 发送回复
	_, err := api.PostMessage(ctx, data.ChannelID, (&dto.PostMessage{
		Content: resp,
	}).AddAtUsr(data.Author.ID))
	if err != nil {
		fmt.Printf("[Send Message] failed, err: %v\n", err)
		return err
	}
	return err
}

func main() {
	// 实例化 token 并读取配置文件
	token, err := token.CreateDefaultToken().ReadFromConfig("config.yaml")
	if err != nil {
		fmt.Println(err)
		return
	}

	// 实例化 OpenAPI
	api = BotApi.NewSandBoxOpenAPI(token)

	go startSendTemporaryService()
	//go startSendOnceService()
	//go startSendDayService()

	// 创建 websocket 客户端并启动
	ctx := context.Background()

	// 获取 websocket 接入信息
	wsAp, err := api.GetGateway(ctx)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	// 注册事件处理函数
	var atMessage event.ATMessageEventHandler = aTMessageEventHandler
	intent := event.RegisterHandlers(atMessage)

	// 创建 websocket 客户端并启动
	wsClient := websocket.CreateClient(wsAp, token, &intent)
	wsClient.Start()
}

func startSendTemporaryService() {
	sendTemporary = Remind.NewSendTemporary(api)
	sendTemporary.SendRemind()
}

func startSendOnceService() {
	remindDate := time.Now().Format(Remind.LayoutDate)
	remindTime := time.Now().Format(Remind.LayoutTime)
	sendOnce = Remind.NewSendRemindInDB(api)
	sendOnce.AddRemindList(Remind.SelectTodayOnce(DB.GetDBInstance(), remindDate, remindTime))
	sendOnce.SendRemind()

	c := cron.New()
	c.AddFunc("@daily", func() {
		remindDate = time.Now().Format(Remind.LayoutDate)
		sendOnce.AddRemindList(Remind.SelectTodayOnce(DB.GetDBInstance(), remindDate, ""))
		sendOnce.SendRemind()
	})
}

func startSendDayService() {
	remindDate := time.Now().Format(Remind.LayoutDate)
	remindTime := time.Now().Format(Remind.LayoutTime)
	sendDay = Remind.NewSendRemindInDB(api)
	sendDay.AddRemindList(Remind.SelectTodayDay(DB.GetDBInstance(), remindDate, remindTime))
	sendDay.SendRemind()

	c := cron.New()
	c.AddFunc("@daily", func() {
		remindDate = time.Now().Format(Remind.LayoutDate)
		sendDay.AddRemindList(Remind.SelectTodayDay(DB.GetDBInstance(), remindDate, ""))
		sendDay.SendRemind()
	})
}
