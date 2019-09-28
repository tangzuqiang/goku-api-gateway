package common

import (
	"net/http"
	"strconv"

	goku_plugin "github.com/eolinker/goku-plugin"
)

var _ goku_plugin.ContextProxy = (*Context)(nil)

//Context context
type Context struct {
	w http.ResponseWriter
	*CookiesHandler
	*PriorityHeader
	*StatusHandler
	*StoreHandler
	RequestOrg           *RequestReader
	ProxyRequest         *Request
	ProxyResponseHandler *ResponseReader
	Body                 []byte
	strategyID           string
	strategyName         string
	apiID                int
	requestID            string
	finalTargetServer    string
	retryTargetServers   string
}

//FinalTargetServer 获取最终目标转发服务器
func (ctx *Context) FinalTargetServer() string {
	return ctx.finalTargetServer
}

//SetFinalTargetServer 设置最终目标服务器
func (ctx *Context) SetFinalTargetServer(finalTargetServer string) {
	ctx.finalTargetServer = finalTargetServer
}

//RetryTargetServers 重试目标服务器
func (ctx *Context) RetryTargetServers() string {
	return ctx.retryTargetServers
}

//SetRetryTargetServers 设置重试目标服务器
func (ctx *Context) SetRetryTargetServers(retryTargetServers string) {
	ctx.retryTargetServers = retryTargetServers
}

//Finish 请求结束
func (ctx *Context) Finish() (n int, statusCode int) {

	header := ctx.PriorityHeader.header

	statusCode = ctx.StatusHandler.code
	if statusCode == 0 {
		statusCode = 504
	}

	bodyAllowed := true
	switch {
	case statusCode >= 100 && statusCode <= 199:
		bodyAllowed = false
		break
	case statusCode == 204:
		bodyAllowed = false
		break
	case statusCode == 304:
		bodyAllowed = false
		break
	}

	if ctx.PriorityHeader.appendHeader != nil {
		for k, vs := range ctx.PriorityHeader.appendHeader.header {
			for _, v := range vs {
				header.Add(k, v)
			}
		}
	}

	if ctx.PriorityHeader.setHeader != nil {
		for k, vs := range ctx.PriorityHeader.setHeader.header {
			header.Del(k)
			for _, v := range vs {
				header.Add(k, v)
			}
		}
	}

	for k, vs := range ctx.PriorityHeader.header {
		if k == "Content-Length" && bodyAllowed {
			vs = []string{strconv.Itoa(len(string(ctx.Body)))}
		}
		for _, v := range vs {
			ctx.w.Header().Add(k, v)
		}
	}

	ctx.w.WriteHeader(statusCode)

	if !bodyAllowed {
		return 0, statusCode
	}
	n, _ = ctx.w.Write(ctx.Body)
	return n, statusCode
}

//RequestId 获取请求ID
func (ctx *Context) RequestId() string {
	return ctx.requestID
}

//NewContext 创建context
func NewContext(r *http.Request, requestID string, w http.ResponseWriter) *Context {
	requestreader := NewRequestReader(r)
	return &Context{
		CookiesHandler:       newCookieHandle(r.Header),
		PriorityHeader:       NewPriorityHeader(),
		StatusHandler:        NewStatusHandler(),
		StoreHandler:         NewStoreHandler(),
		RequestOrg:           requestreader,
		ProxyRequest:         NewRequest(requestreader),
		ProxyResponseHandler: nil,
		requestID:            requestID,
		w:                    w,
	}
}

//SetProxyResponse 设置转发响应
func (ctx *Context) SetProxyResponse(response *http.Response) {

	ctx.ProxyResponseHandler = newResponseReader(response)
	if ctx.ProxyResponseHandler != nil {
		ctx.Body = ctx.ProxyResponseHandler.body
		ctx.SetStatus(ctx.ProxyResponseHandler.StatusCode(), ctx.ProxyResponseHandler.Status())
		ctx.header = ctx.ProxyResponseHandler.header
	}

}
func (ctx *Context) Write(w http.ResponseWriter) {
	if ctx.StatusCode() == 0 {
		ctx.SetStatus(200, "200 ok")
	}
	if ctx.Body != nil {
		w.Write(ctx.Body)
	}

	w.WriteHeader(ctx.StatusCode())

}

//GetBody 获取body内容
func (ctx *Context) GetBody() []byte {
	return ctx.Body
}

//SetBody 设置body内容
func (ctx *Context) SetBody(data []byte) {
	ctx.Body = data
}

//ProxyResponse 转发响应
func (ctx *Context) ProxyResponse() goku_plugin.ResponseReader {
	return ctx.ProxyResponseHandler
}

//StrategyId 获取策略ID
func (ctx *Context) StrategyId() string {
	return ctx.strategyID
}

//SetStrategyId 设置策略ID
func (ctx *Context) SetStrategyId(strategyID string) {
	ctx.strategyID = strategyID
}

//StrategyName 获取策略名称
func (ctx *Context) StrategyName() string {
	return ctx.strategyName
}

//SetStrategyName 设置策略名称
func (ctx *Context) SetStrategyName(strategyName string) {
	ctx.strategyName = strategyName
}

//ApiID 获取接口ID
func (ctx *Context) ApiID() int {
	return ctx.apiID
}

//SetAPIID 设置接口ID
func (ctx *Context) SetAPIID(apiID int) {
	ctx.apiID = apiID
}

//Request 获取请求原始数据
func (ctx *Context) Request() goku_plugin.RequestReader {
	return ctx.RequestOrg
}

//Proxy 获取代理请求
func (ctx *Context) Proxy() goku_plugin.Request {
	return ctx.ProxyRequest
}
