package cron

import (
	"ethereum-parser/config"
	evm "ethereum-parser/pkg/ethereum-rpc-client"
	pubsub "ethereum-parser/pkg/pub-sub"
	"fmt"
	"time"

	"github.com/robfig/cron"
)

type ListenEthereumBlockCron struct {
	RpcUrl    string
	Period    string
	Publisher *pubsub.BlockPublisher
}

func NewListenEthereumBlockCron(publisher *pubsub.BlockPublisher) *ListenEthereumBlockCron {
	cronConfig := config.GetConfig().Cron

	return &ListenEthereumBlockCron{
		RpcUrl:    cronConfig.Url,
		Period:    cronConfig.Period,
		Publisher: publisher,
	}
}

func (c *ListenEthereumBlockCron) Start() {
	// Do something

	cronInstance := cron.New()

	initBlockNumber, err := evm.GetBlockNumber()
	if err != nil {
		fmt.Printf("Error getting block number: %v\n", err)
		return
	}

	lastUpdateBlock := initBlockNumber - 1

	period := c.Period
	err = cronInstance.AddFunc(period, func() {
		currentBlock, err := evm.GetBlockNumber()
		if err != nil {
			fmt.Printf("Error getting block number: %v\n", err)
			return
		}

		for i := lastUpdateBlock + 1; i <= currentBlock; i++ {
			block, err := evm.GetBlockByNumber(i)
			if err != nil {
				fmt.Printf("Error getting block by number: %v\n", err)
				return
			}

			c.Publisher.AddBlock(block)
			c.Publisher.Publish(block)

			lastUpdateBlock = i

			time.Sleep(1 * time.Second)
		}
	})

	cronInstance.Start()

	select {}

}
