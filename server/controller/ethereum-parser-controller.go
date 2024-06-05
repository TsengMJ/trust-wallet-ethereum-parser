package controller

import (
	"encoding/json"
	"net/http"

	"ethereum-parser/logger"
	ethereumParser "ethereum-parser/pkg/ethereum-parser"
	"ethereum-parser/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func HandleWebSocket(c *gin.Context) {
	log := logger.Logger

	var parser ethereumParser.EthereumParser = ethereumParser.NewBasicEthereumParser()
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Error("Failed to upgrade connection, " + err.Error())
		return
	}
	defer conn.Close()

	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Error("Failed to read message, " + err.Error())
				conn.WriteJSON(util.GetWebsocketFailResponse("Failed to read message"))
			}
			continue
		}

		var request map[string]interface{}
		if err := json.Unmarshal(message, &request); err != nil {
			log.Error("Failed to unmarshal message, " + err.Error())

			err = conn.WriteJSON(util.GetWebsocketFailResponse("Fai"))
			if err != nil {
				log.Error("Failed to write message, " + err.Error())
				conn.WriteJSON(util.GetWebsocketFailResponse("Failed to unmarshal message"))
			}

			continue
		}

		switch request["action"] {
		case "GetCurrentBlock":
			block, err := parser.GetCurrentBlock()
			if err != nil {
				log.Error("Failed to get current block, " + err.Error())
				conn.WriteJSON(util.GetWebsocketFailResponse("Failed to get current block"))
				continue
			}

			response := map[string]interface{}{
				"action": "GetCurrentBlock",
				"block":  block,
			}
			if err := conn.WriteJSON(util.GetWebsocketSuccessResponse(response)); err != nil {
				log.Error("Failed to write message, " + err.Error())
			}

		case "Subscribe":
			address, ok := request["address"].(string)
			if !ok {
				log.Error("Invalid address format, " + err.Error())
				conn.WriteJSON(util.GetWebsocketFailResponse("Invalid address format"))
				continue
			}

			subscribed, err := parser.Subscribe(address)
			if err != nil {
				log.Error("Failed to subscribe, " + err.Error())
				conn.WriteJSON(util.GetWebsocketFailResponse("Failed to subscribe"))
				continue
			}

			response := map[string]interface{}{
				"action":     "Subscribe",
				"subscribed": subscribed,
			}
			if err := conn.WriteJSON(util.GetWebsocketSuccessResponse(response)); err != nil {
				log.Error("Failed to write message, " + err.Error())
			}

		case "GetTransactions":
			transactions, err := parser.GetTransactions()
			if err != nil {
				log.Error("Failed to get transactions, " + err.Error())
				conn.WriteJSON(util.GetWebsocketFailResponse("Failed to get transactions"))
				continue
			}

			response := map[string]interface{}{
				"action":       "GetTransactions",
				"transactions": transactions,
			}
			if err := conn.WriteJSON(util.GetWebsocketSuccessResponse(response)); err != nil {
				log.Error("Failed to write message, " + err.Error())
			}

		default:
			if err := conn.WriteJSON(util.GetWebsocketFailResponse("Invalid Action")); err != nil {
				log.Error("Failed to write message, " + err.Error())
			}
		}
	}

}
