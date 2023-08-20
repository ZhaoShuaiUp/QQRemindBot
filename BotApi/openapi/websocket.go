package openapi

import (
	"RemindBot/BotApi/dto"
	"context"
)

// 获取 websocket 接入信息
func (api *OpenAPI) GetGateway(ctx context.Context) (*dto.WebsocketAP, error) {
	resp, err := api.request(ctx).
		SetResult(dto.WebsocketAP{}).
		Get(api.getURL(URI_GET_GATEWAY))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.WebsocketAP), nil
}
