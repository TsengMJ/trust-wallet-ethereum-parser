package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"ethereum-parser/logger"

	ethereumParser "ethereum-parser/pkg/ethereum-parser"

	evm "ethereum-parser/pkg/ethereum-rpc-client"
	pubsub "ethereum-parser/pkg/pub-sub"
	"ethereum-parser/util"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

func ProduceMessage(c *gin.Context) {
	publisher := pubsub.DefaultPublisher

	// mock message
	message := &pubsub.Message{
		Data: []byte("Hello, World!"),
	}

	publisher.Publish(c, message)

	c.JSON(http.StatusOK, gin.H{
		"message": "Message published",
	})
}

func HandleWebSocket(c *gin.Context) {
	log := logger.Logger
	publisher := pubsub.DefaultPublisher

	parser := ethereumParser.NewBasicEthereumParser()

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

	subscriber := pubsub.NewSubscriber()
	publisher.Subscribe(c, subscriber)
	defer publisher.Unsubscribe(c, subscriber)

	go func() {
		for {
			select {
			case msg := <-subscriber.Handler:
				if msg.Data == nil {
					continue
				}

				if len(parser.Subscriptions) == 0 {
					continue
				}

				var targetTxs []evm.Transaction

				for _, tx := range msg.Data.(*evm.Block).Transactions {
					fmt.Println(tx.From, tx.To, parser.Subscriptions[strings.ToLower(tx.From)], parser.Subscriptions[strings.ToLower(tx.To)])
					if parser.Subscriptions[strings.ToLower(tx.From)] || parser.Subscriptions[strings.ToLower(tx.To)] {
						targetTxs = append(targetTxs, tx)
					}
				}

				if len(targetTxs) == 0 {
					continue
				}

				response := map[string]interface{}{
					"action": "Subscribe",
					"txs":    targetTxs,
				}

				err := conn.WriteJSON(util.GetWebsocketSuccessResponse(response))
				if err != nil {
					log.Error("Failed to write message, " + err.Error())
				}

			case <-subscriber.Quit:
				return
			case <-c.Done():
				return
			case <-c.Writer.CloseNotify():
				break
			}
		}
	}()

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
		case "UnSubscribe":
			address, ok := request["address"].(string)
			if !ok {
				log.Error("Invalid address format, " + err.Error())
				conn.WriteJSON(util.GetWebsocketFailResponse("Invalid address format"))
				continue
			}

			unsubscribed, err := parser.UnSubscribe(address)
			if err != nil {
				log.Error("Failed to unsubscribe, " + err.Error())
				conn.WriteJSON(util.GetWebsocketFailResponse("Failed to unsubscribe"))
				continue
			}

			response := map[string]interface{}{
				"action":       "UnSubscribe",
				"unsubscribed": unsubscribed,
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
