package dto

import (
	"regexp"
	"strings"
)

type Message struct {
	ID        string `json:"id"`
	ChannelID string `json:"channel_id"`
	Content   string `json:"content"`
	Author    *User
}

type PostMessage struct {
	Content string `json:"content,omitempty"`
	// MsgID   string `json:"msg_id,omitempty"`
}

// 消息体头部 at 某人
var atRE = regexp.MustCompile(`<@!\d+>`)

// ETLMessage 去除消息体头部 at 某人，并 trim
func ETLMessage(msg string) string {
	etlData := string(atRE.ReplaceAll([]byte(msg), []byte("")))
	etlData = strings.TrimSpace(etlData)
	return etlData
}

// 获取消息体头部 at 某人
func MentionUser(userID string) string {
	return "<@!" + userID + "> "
}

// 创建回复某认的消息体
func (pm *PostMessage) RestMessageWithAtUsr(userID string, content string) *PostMessage {
	pm.Content = MentionUser(userID) + content
	return pm
}

// 修改为回复某人的消息体
func (pm *PostMessage) AddAtUsr(userID string) *PostMessage {
	pm.Content = MentionUser(userID) + pm.Content
	return pm
}
