package openapi

import (
	"RemindBot/BotApi/dto"
	"context"
	"fmt"
)

// 发送消息
func (api *OpenAPI) PostMessage(ctx context.Context, channelID string, message *dto.PostMessage) (*dto.Message, error) {
	fmt.Printf("[PostMessage] channelID: %s, message: %v\n", channelID, message)
	resp, err := api.request(ctx).
		SetResult(dto.Message{}).
		SetPathParam("channel_id", channelID).
		SetBody(message).
		Post(api.getURL(URI_POST_MESSAGE))
	if err != nil {
		return nil, err
	}

	return resp.Result().(*dto.Message), nil
}
