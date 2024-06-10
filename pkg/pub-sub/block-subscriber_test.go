package pubsub_test

import (
	"testing"

	"github.com/stretchr/testify/assert"

	evm "ethereum-parser/pkg/ethereum-rpc-client"
	pubsub "ethereum-parser/pkg/pub-sub"
)

func TestNewBlockSubscriber(t *testing.T) {
	subscriber := pubsub.NewBlockSubscriber()

	assert.NotNil(t, subscriber.Handler, "Handler channel should not be nil")
	assert.NotNil(t, subscriber.Quit, "Quit channel should not be nil")
}

func TestBlockSubscriber_Publish(t *testing.T) {
	// Create a new block subscriber
	subscriber := pubsub.NewBlockSubscriber()

	// Create a block to publish
	block := &evm.Block{
		Number: "0x1",
		Hash:   "0x123",
	}

	// Publish the block
	subscriber.Publish(block)

	// Assert that the block was received by the subscriber's handler channel
	receivedBlock := <-subscriber.Handler
	assert.Equal(t, block, receivedBlock, "Received block does not match published block")
}
