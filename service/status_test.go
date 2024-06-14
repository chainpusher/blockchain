package service_test

import (
	"errors"
	"github.com/chainpusher/blockchain/service"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestStatus_Error(t *testing.T) {
	getError := func() error {
		return &service.Status{
			Code:     service.Other,
			RpcError: nil,
		}
	}
	err := getError()

	status, ok := service.FromError(err)

	assert.Truef(t, ok, "Error should be a status")
	assert.Equal(t, service.Other, status.Code)
	assert.Nil(t, status.RpcError)
	assert.Falsef(t, service.IsNotFound(err), "Error should not be NotFound")
}

func TestStatus_ErrorWithRpcError(t *testing.T) {
	err := errors.New("not found")

	status := service.NewEthereumError(err)
	assert.Truef(t, service.IsNotFound(status), "Error should be NotFound")
}
