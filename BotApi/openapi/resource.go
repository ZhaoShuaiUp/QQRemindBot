package openapi

import (
	"fmt"
)

const (
	domain        = "api.sgroup.qq.com"
	sandBoxDomain = "sandbox.api.sgroup.qq.com"
	scheme        = "https"
)

type uri string

// 接口地址常量
const (
	URI_POST_MESSAGE uri = "/channels/{channel_id}/messages"

	URI_GET_GATEWAY     uri = "/gateway"
	URI_GET_GATEWAY_BOT uri = "/gateway/bot"
)

func (api *OpenAPI) getURL(endpoint uri) string {
	d := domain
	if api.isSandBox {
		d = sandBoxDomain
	}
	return fmt.Sprintf("%s://%s%s", scheme, d, endpoint)
}
