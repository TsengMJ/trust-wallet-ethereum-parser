package util

type Response struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func GetSuccessResponse(data interface{}) *Response {
	return &Response{
		Data:  data,
		Error: "",
	}
}

func GetFailResponse(err string) *Response {
	return &Response{
		Data:  nil,
		Error: err,
	}
}
