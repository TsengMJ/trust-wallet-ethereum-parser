package controller

import (
	evm "ethereum-parser/pkg/ethereum-rpc-client"
	pubsub "ethereum-parser/pkg/pub-sub"
	"ethereum-parser/util"
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func GetCurrentBlock(c *gin.Context) {
	publisher := pubsub.DefaultPublisher

	block, err := publisher.GetLatestBlock()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetFailResponse(err.Error()))
		return
	}

	response := map[string]interface{}{
		"block": block,
	}
	c.JSON(http.StatusOK, util.GetSuccessResponse(response))
}

func GetCurrentBlockTransactionsByAddress(c *gin.Context) {
	address := c.Param("address")

	if util.IsValidAddress(address) == false {
		c.JSON(http.StatusBadRequest, util.GetFailResponse("Invalid address"))
		return
	}

	publisher := pubsub.DefaultPublisher
	block, err := publisher.GetLatestBlock()
	if err != nil {
		c.JSON(http.StatusInternalServerError, util.GetFailResponse(err.Error()))
		return
	}

	blockTxs := &block.Transactions
	addressMapping := map[string]bool{strings.ToLower(address): true}
	filteredTxs := filterTransactionsByAddresses(blockTxs, &addressMapping)

	response := map[string]interface{}{
		"address":      address,
		"transactions": filteredTxs,
	}
	c.JSON(http.StatusOK, util.GetSuccessResponse(response))
}

func filterTransactionsByAddresses(transactions *[]evm.Transaction, address *map[string]bool) []evm.Transaction {
	var targetTxs []evm.Transaction

	for _, tx := range *transactions {
		if (*address)[strings.ToLower(tx.From)] || (*address)[strings.ToLower(tx.To)] {
			targetTxs = append(targetTxs, tx)
		}
	}

	return targetTxs
}
