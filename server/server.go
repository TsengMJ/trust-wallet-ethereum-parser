package server

import (
	"ethereum-parser/config"
	"ethereum-parser/server/controller"
	"strconv"

	"github.com/gin-gonic/gin"
)

func StartServer() {
	r := gin.Default()

	r.GET("/ws", controller.HandleWebSocket)
	r.GET("/current-block", controller.GetCurrentBlock)
	r.GET("/transaction/:address", controller.GetCurrentBlockTransactionsByAddress)

	port := config.Config.Server.Port
	r.Run(":" + strconv.Itoa(port))
}
