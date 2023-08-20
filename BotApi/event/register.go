package event

import "RemindBot/BotApi/dto"

var DefaultHandlers struct {
	ATMessage ATMessageEventHandler
}

// ATMessageEventHandler at 机器人消息事件 handler
type ATMessageEventHandler func(event *dto.WSPayload, data *dto.WSATMessageData) error

// 事件处理函数注册
func RegisterHandlers(handlers ...interface{}) dto.Intent {
	var i dto.Intent = 0
	for _, h := range handlers {
		switch handle := h.(type) {
		case ATMessageEventHandler:
			DefaultHandlers.ATMessage = handle
			i = i | 1<<30
		}
	}
	return i
}
