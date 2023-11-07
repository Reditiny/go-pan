package response

import "gopan/enums"

// Response 统一响应
type Response struct {
	Status string      `json:"status"`
	Code   int         `json:"code"`
	Info   string      `json:"info"`
	Data   interface{} `json:"data"`
}

// Success 成功
func Success(data interface{}) Response {
	return Response{"success", 200, "请求成功", data}
}

// Fail 失败
func Fail(r enums.ResponseCode) Response {
	return Response{"fail", r.Code, r.Info, nil}
}
