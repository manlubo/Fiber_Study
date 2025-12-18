package response

type Response struct {
	Success bool   `json:"success"`
	Code    string `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data"`
}

func OK(message string, data any) Response {
	return Response{
		Success: true,
		Code:    "OK",
		Message: message,
		Data:    data,
	}
}

func Error(code string, message string, data any) Response {
	return Response{
		Success: false,
		Code:    code,
		Message: message,
		Data:    data,
	}
}
