package helper

//EmptyObj object is used when data doesnt want to be null on json
type EmptyObj struct{}

type Response struct {
	Status  bool        `json:"status"`
	Message string      `json:"message"`
	Data    interface{} `json:"data"`
}

func BuildSuccessResponse(status bool, message string, data interface{}) Response {
	res := Response{
		Status:  status,
		Message: message,
		Data:    data,
	}
	return res
}

func BuildErrorResponse(message string, data interface{}) Response {
	res := Response{
		Message: message,
		Data:    data,
	}
	return res
}
