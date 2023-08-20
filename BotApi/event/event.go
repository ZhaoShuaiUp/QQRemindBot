package event

import (
	"RemindBot/BotApi/dto"
	"encoding/json"
	"fmt"

	"github.com/tidwall/gjson"
)

// 事件处理函数类型
type eventParseFunc func(event *dto.WSPayload, message []byte) error

var eventParseFuncMap = map[dto.OPCode]map[dto.EventType]eventParseFunc{
	dto.OPCODE_DISPATCH: {
		dto.EVENT_CODE_AT_MESSAGE_CREATE: atMessageHandler,
	},
}

// 解析事件数据
func ParseData(message []byte, target interface{}) error {
	data := gjson.Get(string(message), "d")
	return json.Unmarshal([]byte(data.String()), target)
}

// 解析并处理
func ParseAndHandle(payload *dto.WSPayload) error {
	// 指定类型的 handler
	if h, ok := eventParseFuncMap[payload.OPCode][payload.Type]; ok {
		return h(payload, payload.RawMessage)
	}
	// 未指定则不处理
	return nil
}

// 默认 at 事件处理
func atMessageHandler(payload *dto.WSPayload, message []byte) error {
	data := &dto.WSATMessageData{}
	if err := ParseData(message, data); err != nil {
		return err
	}
	fmt.Printf("[atMessage] message %v\n", data)
	// 调用用户注册的处理函数
	if DefaultHandlers.ATMessage != nil {
		return DefaultHandlers.ATMessage(payload, data)
	}
	return nil
}
