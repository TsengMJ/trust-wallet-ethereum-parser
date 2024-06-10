package controller

import (
	"encoding/json"
	"errors"
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

func HandleWebSocket(c *gin.Context) {
	log := logger.Logger
	publisher := pubsub.DefaultPublisher

	parser := ethereumParser.NewBasicEthereumParser()

	conn, err := getWebsocketConnection(c)
	if err != nil {
		log.Error("Failed to get websocket connection, " + err.Error())
		c.JSON(http.StatusInternalServerError, util.GetFailResponse("Failed to get websocket connection"))
		return
	}
	defer conn.Close()

	subscriber := pubsub.NewBlockSubscriber()
	err = publisher.Subscribe(subscriber)
	if err != nil {
		log.Error("Failed to subscribe, " + err.Error())
		conn.WriteJSON(util.GetFailResponse("Failed to subscribe"))
		return
	}
	defer publisher.Unsubscribe(subscriber)

	go notifySubscribers(conn, subscriber, parser)

	// Handle incoming actions (GetCurrentBlock, Subscribe, UnSubscribe)
	for {
		request, err := getRequest(conn)
		if err != nil {
			if err := conn.Close(); err != nil {
				log.Error("Failed to close connection, " + err.Error())
				break
			}

			log.Error("Failed to get websocket request, " + err.Error())
			conn.WriteJSON(util.GetFailResponse("Failed to get websocket request, " + err.Error()))
			continue
		}

		switch request["action"] {
		case "GetCurrentBlock":
			err := handleGetCurrentBlock(conn, parser)
			if err != nil {
				log.Error("Failed to handle GetCurrentBlock, " + err.Error())
			}
		case "Subscribe":
			address, ok := request["address"].(string)
			if !ok {
				log.Error("Invalid address format, " + err.Error())
				conn.WriteJSON(util.GetFailResponse("Invalid address format"))
				continue
			}

			err := handleSubscribe(conn, parser, address)
			if err != nil {
				log.Error("Failed to handle Subscribe, " + err.Error())
			}
		case "UnSubscribe":
			address, ok := request["address"].(string)
			if !ok {
				log.Error("Invalid address format, " + err.Error())
				conn.WriteJSON(util.GetFailResponse("Invalid address format"))
				continue
			}

			err := handleUnSubscribe(conn, parser, address)
			if err != nil {
				log.Error("Failed to handle UnSubscribe, " + err.Error())
			}

		default:
			if err := conn.WriteJSON(util.GetFailResponse("Invalid Action")); err != nil {
				log.Error("Failed to write message, " + err.Error())
			}
		}
	}
}

func getWebsocketConnection(c *gin.Context) (*websocket.Conn, error) {
	upgrader := websocket.Upgrader{
		ReadBufferSize:  1024,
		WriteBufferSize: 1024,
		CheckOrigin: func(r *http.Request) bool {
			return true
		},
	}

	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		return nil, err
	}

	return conn, nil
}

func notifySubscribers(conn *websocket.Conn, subscriber *pubsub.BlockSubscriber, parser *ethereumParser.BasicEthereumParser) {
	for {
		select {
		case block := <-subscriber.Handler:
			if block == nil {
				continue
			}

			if len(parser.Subscriptions) == 0 {
				continue
			}

			var targetTxs []evm.Transaction

			for _, tx := range block.Transactions {
				if parser.Subscriptions[strings.ToLower(tx.From)] || parser.Subscriptions[strings.ToLower(tx.To)] {
					targetTxs = append(targetTxs, tx)
				}
			}

			if len(targetTxs) == 0 {
				continue
			}

			response := map[string]interface{}{
				"action": "Transactions",
				"txs":    targetTxs,
			}

			err := conn.WriteJSON(util.GetSuccessResponse(response))
			if err != nil {
				logger.Logger.Error("Failed to write message, " + err.Error())
			}

		case <-subscriber.Quit:
			return
		}
	}
}

func getRequest(conn *websocket.Conn) (map[string]interface{}, error) {
	_, message, err := conn.ReadMessage()
	if err != nil {
		return nil, errors.New("Connection failed to read message, " + err.Error())
	}

	var request map[string]interface{}
	if err := json.Unmarshal(message, &request); err != nil {
		return nil, errors.New("Failed to unmarshal message, " + err.Error())
	}

	return request, nil
}

func handleGetCurrentBlock(conn *websocket.Conn, parser *ethereumParser.BasicEthereumParser) error {
	block, err := parser.GetCurrentBlock()
	if err != nil {
		logger.Logger.Error("Failed to get current block, " + err.Error())
		conn.WriteJSON(util.GetFailResponse("Failed to get current block"))
		return err
	}

	response := map[string]interface{}{
		"action": "GetCurrentBlock",
		"block":  block,
	}
	if err := conn.WriteJSON(util.GetSuccessResponse(response)); err != nil {
		logger.Logger.Error("Failed to write message, " + err.Error())
	}

	return nil
}

func handleSubscribe(conn *websocket.Conn, parser *ethereumParser.BasicEthereumParser, address string) error {
	subscribed, err := parser.Subscribe(address)
	if err != nil {
		logger.Logger.Error("Failed to subscribe, " + err.Error())
		conn.WriteJSON(util.GetFailResponse("Failed to subscribe"))
		return err
	}

	response := map[string]interface{}{
		"action":     "Subscribe",
		"subscribed": subscribed,
	}
	if err := conn.WriteJSON(util.GetSuccessResponse(response)); err != nil {
		logger.Logger.Error("Failed to write message, " + err.Error())
	}

	return nil
}

func handleUnSubscribe(conn *websocket.Conn, parser *ethereumParser.BasicEthereumParser, address string) error {
	unsubscribed, err := parser.UnSubscribe(address)
	if err != nil {
		logger.Logger.Error("Failed to unsubscribe, " + err.Error())
		conn.WriteJSON(util.GetFailResponse("Failed to unsubscribe"))
		return err
	}

	response := map[string]interface{}{
		"action":       "UnSubscribe",
		"unsubscribed": unsubscribed,
	}

	if err := conn.WriteJSON(util.GetSuccessResponse(response)); err != nil {
		logger.Logger.Error("Failed to write message, " + err.Error())
	}

	return nil
}
