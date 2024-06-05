package util

type WebsocketResponse struct {
	Data  interface{} `json:"data"`
	Error string      `json:"error"`
}

func GetWebsocketSuccessResponse(data interface{}) *WebsocketResponse {
	return &WebsocketResponse{
		Data:  data,
		Error: "",
	}
}

func GetWebsocketFailResponse(err string) *WebsocketResponse {
	return &WebsocketResponse{
		Data:  nil,
		Error: err,
	}
}
