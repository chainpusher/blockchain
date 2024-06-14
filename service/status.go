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
	Code  Code
	Cause error
}

func (s *Status) GetCode() Code {
	return s.Code
}

func (s *Status) GetCause() error {
	return s.Cause
}

func (s *Status) Error() string {
	return s.Cause.Error()
}

func NewError(code Code, rpcError error) error {
	_ = OK // avoid unused error
	return &Status{
		Code:  code,
		Cause: rpcError,
	}
}

func NewEthereumError(err error) error {
	_, ok := err.(interface{ Error() string })
	if ok && err.Error() == "not found" {
		return NewError(NotFound, err)
	}

	return NewError(RpcError, err)
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

func IsNotFound(err error) bool {
	status, ok := FromError(err)
	if !ok {
		return false
	}
	return status.Code == NotFound
}
