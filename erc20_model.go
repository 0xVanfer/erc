package erc

import "github.com/0xVanfer/abigen/erc20"

// Basic info of ERC20 token.
type ERC20Info struct {
	Network  *string
	Address  *string
	Symbol   *string
	Decimals *int
	Contract *erc20.Erc20
}

// Intend to read by timestamp. But do not want to import blockscan.
type HistoryOpt struct {
	BlockNumber int64
}
