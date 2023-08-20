package BotApi

import (
	"RemindBot/BotApi/openapi"
	"RemindBot/BotApi/token"
)

// 创建 OpenAPI 实例
func NewOpenAPI(token *token.Token) *openapi.OpenAPI {
	return openapi.CreateOpenAPI(token, false)
}

// 创建沙箱 OpenAPI 实例
func NewSandBoxOpenAPI(token *token.Token) *openapi.OpenAPI {
	return openapi.CreateOpenAPI(token, true)
}
