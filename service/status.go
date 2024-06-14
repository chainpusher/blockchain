package service

import "errors"

type Code uint32

const (
	OK Code = iota
	NotFound
	RpcError
	Other
)

type Status struct {
	Code     Code
	RpcError error
}

func (s *Status) GetCode() Code {
	return s.Code
}

func (s *Status) GetRpcError() error {
	return s.RpcError
}

func (s *Status) Error() string {
	return s.RpcError.Error()
}

func NewError(code Code, rpcError error) error {
	return &Status{
		Code:     code,
		RpcError: rpcError,
	}
}

func FromError(err error) (s *Status, ok bool) {
	if err == nil {
		return nil, true
	}

	// if err is Status
	var status *Status
	if errors.As(err, &status) {
		return status, true
	}
	return nil, false
}

func NewStatus(code Code, err error) *Status {
	return &Status{
		Code:     code,
		RpcError: err,
	}
}
