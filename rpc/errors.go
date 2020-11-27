package rpc

import "fmt"

type ErrResponse struct {
	Code    int         `json:"code"`
	Message string      `json:"message"`
	Data    interface{} `json:"data,omitempty"`
}

func (e *ErrResponse) Error() string {
	return fmt.Sprintf("rpc error (code %d): %s", e.Code, e.Message)
}

type ErrNoERC20Methods struct {
}

func (e *ErrNoERC20Methods) Error() string {
	return "not implementing one of ERC20 contract's method 'name()', or 'symbol()' or 'decimals()'"
}
