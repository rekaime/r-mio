package response

type ErrorCode int

const (
	SuccessCode ErrorCode = iota
)

var ErrorCodeInfo = map[ErrorCode]string{
	SuccessCode: "success",
}