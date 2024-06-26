package service

import (
	"encoding/json"
	"math/big"
	"testing"

	"github.com/fbsobreira/gotron-sdk/pkg/client"
	"github.com/sirupsen/logrus"
	"github.com/stretchr/testify/assert"
)

var c *client.GrpcClient

func TestMain(m *testing.M) {
	var err error
	c, err = NewTronClient()
	if err != nil {
		logrus.Errorf("Error creating tron client: %v", err)
	}
	m.Run()
}

func TestTronService_GetLatestBlock(t *testing.T) {

	service := NewTronBlockChainService(c, []BlockListener{})

	block, err := service.GetLatestBlock()
	if err != nil {
		t.Errorf("GetLatestBlock() error = %v", err)
		return
	}

	assert.NotNil(t, block)
}

func TestTronService_GetBlock(t *testing.T) {
	service := NewTronBlockChainService(c, []BlockListener{})

	// 62606270
	// 62544893
	block, err := service.GetBlock(big.NewInt(62606533))
	if err != nil {
		t.Errorf("GetBlock() error = %v", err)
		return
	}

	serialized, err := json.Marshal(block)
	assert.Nil(t, err)
	assert.NotNil(t, serialized)

	assert.NotNil(t, block)
	assert.Greater(t, len(block.Transactions), 0, "Block should have transactions")
}
