package rpc

type ErrNoERC20Methods struct {
}

func (e *ErrNoERC20Methods) Error() string {
	return "not implementing one of ERC20 contract's method 'name()', or 'symbol()' or 'decimals()'"
}
