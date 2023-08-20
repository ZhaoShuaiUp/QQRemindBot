package openapi

import (
	"RemindBot/BotApi/token"
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/go-resty/resty/v2"
)

const (
	traceIDKey = "X-Tps-trace-ID"
)

type OpenAPI struct {
	token   *token.Token
	timeout time.Duration

	isSandBox   bool   // 请求是否为沙盒模式
	lastTraceID string // 最后一次请求的 TraceID

	restyClient *resty.Client // resty 客户端
}

// 创建 OpenAPI 实例
func CreateOpenAPI(token *token.Token, isSandBox bool) *OpenAPI {
	api := &OpenAPI{
		token:       token,
		timeout:     3 * time.Second,
		isSandBox:   false,
		restyClient: resty.New(),
	}
	api.initRestyClient()
	return api
}

// 初始化 resty client
func (api *OpenAPI) initRestyClient() {
	if api.restyClient == nil {
		api.restyClient = resty.New()
	}
	api.restyClient.SetTimeout(api.timeout).
		SetAuthToken(api.token.GetAuthorization()).
		SetAuthScheme(api.token.Type).
		//SetHeader("User-Agent", "QQBotSDK/1.0.0").
		OnAfterResponse(
			func(client *resty.Client, response *resty.Response) error {
				fmt.Println(respInfo(response))
				api.lastTraceID = response.Header().Get(traceIDKey)
				return nil
			},
		)
}

// 发送一个请求
func (api *OpenAPI) request(ctx context.Context) *resty.Request {
	return api.restyClient.R().SetContext(ctx)
}

// 格式化请求/响应参数
func respInfo(resp *resty.Response) string {
	bodyJson, _ := json.Marshal(resp.Request.Body)
	return fmt.Sprintf(
		"[OPENAPI] method: %s, url: %s, traceID: %v, status: %v, req: %v, resp: %v",
		resp.Request.Method,
		resp.Request.URL,
		resp.Header().Get(traceIDKey),
		resp.Status(),
		string(bodyJson),
		string(resp.Body()),
	)
}

// 获取 TraceID
func (api *OpenAPI) GetTraceID() string {
	return api.lastTraceID
}
