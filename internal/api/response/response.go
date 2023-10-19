package response

type Response struct {
	Status string      `json:"status"`
	Error  string      `json:"error,omitempty"`
	Data   interface{} `json:"data,omitempty"`
}

func Error(msg string) Response {
	return Response{
		Status: "Error",
		Error:  msg,
	}
}

func OK(d interface{}) Response {
	return Response{
		Status: "OK",
		Data:   d,
	}
}
