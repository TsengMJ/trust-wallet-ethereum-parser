package pubsub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	evm "ethereum-parser/pkg/ethereum-rpc-client"
	pubsub "ethereum-parser/pkg/pub-sub"
)

func TestBlockPublisher_Subscribe(t *testing.T) {
	publisher := pubsub.NewBlockPublisher()
	subscriber := pubsub.NewBlockSubscriber()

	err := publisher.Subscribe(subscriber)

	assert.NoError(t, err, "Error should be nil")
}

func TestBlockPublisher_Unsubscribe(t *testing.T) {
	publisher := pubsub.NewBlockPublisher()
	subscriber := pubsub.NewBlockSubscriber()

	err := publisher.Unsubscribe(subscriber)

	assert.NoError(t, err, "Error should be nil")
}

func TestBlockPublisher_AddBlock(t *testing.T) {
	publisher := pubsub.NewBlockPublisher()
	block := &evm.Block{
		Number: "0x1",
		Hash:   "0x123",
	}

	err := publisher.AddBlock(block)

	assert.NoError(t, err, "Error should be nil")
}

func TestBlockPublisher_GetLatestBlock(t *testing.T) {
	publisher := pubsub.NewBlockPublisher()
	block := &evm.Block{
		Number: "0x1",
		Hash:   "0x123",
	}

	publisher.AddBlock(block)

	latestBlock, err := publisher.GetLatestBlock()

	assert.NoError(t, err, "Error should be nil")
	assert.Equal(t, block, latestBlock, "Latest block should match the added block")
}
