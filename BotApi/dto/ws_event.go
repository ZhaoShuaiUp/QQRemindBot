package dto

// Intent 消息类型
type Intent int

const (
	EVENT_CODE_AT_MESSAGE_CREATE     EventType = "AT_MESSAGE_CREATE"
	EVENT_CODE_DIRECT_MESSAGE_CREATE EventType = "DIRECT_MESSAGE_CREATE"
)
