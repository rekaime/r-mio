package response

type Response struct {
	Code    int         `json:"code"`
    Message string      `json:"message"`
    Data    interface{} `json:"data,omitempty"`
}

func Success(data interface{}) *Response {
	return &Response{
		Code:    int(SuccessCode),
		Message: ErrorCodeInfo[SuccessCode],
		Data:    data,
	}
}

func Error(code ErrorCode) *Response {
    return &Response{
        Code:    int(code),
        Message: ErrorCodeInfo[code],
        Data:    nil,
    }
}