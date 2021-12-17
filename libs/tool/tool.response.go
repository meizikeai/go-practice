package tool

type response struct {
	Code    int64       `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

var empty = []string{}
var responseMessage = map[int64]string{
	200: "操作成功",
	400: "非法请求",
	403: "操作被禁止",
	404: "未找到",
	405: "请求的方法不支持",
	415: "请求错误",
	422: "参数错误",
	500: "未知错误，请检查您的网络",
}

func NewResponse(code int64, data interface{}) *response {
	message := ""

	if _, ok := responseMessage[code]; ok {
		message = responseMessage[code]
	}

	if data == nil {
		data = empty
	}

	return &response{
		Code:    code,
		Message: message,
		Data:    data,
	}
}
