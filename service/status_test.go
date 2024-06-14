package service_test

import (
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
