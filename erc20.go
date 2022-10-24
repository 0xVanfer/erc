package erc

import (
	"errors"
	"strings"

	"github.com/0xVanfer/abigen/erc20"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
)

// Basic info of ERC20 token.
type ERC20Info struct {
	Network  *string
	Address  *string
	Symbol   *string
	Decimals *int
	Contract *erc20.Erc20
}

// Create a new ERC20 token.
func NewErc20(address string, network string, client bind.ContractBackend) (*ERC20Info, error) {
	err := addressRegularCheck(address)
	if err != nil {
		return nil, err
	}
	var new ERC20Info
	new.Network = &network
	checksumed := checksumEthereumAddress(address)
	new.Address = &checksumed
	new.Contract, err = erc20.NewErc20(types.ToAddress(address), client)
	if err != nil {
		return nil, err
	}
	decimals, err := new.Contract.Decimals(nil)
	if err != nil {
		return nil, err
	}
	decimalint := int(decimals.Int64())
	new.Decimals = &decimalint
	var symbol string
	// fuck maker
	if fuckMaker(address, network) {
		symbol = "MKR"
	} else {
		symbol, err = new.Contract.Symbol(nil)
		if err != nil {
			return nil, err
		}
	}
	new.Symbol = &symbol
	return &new, nil
}

// Maker symbol is bytes instead of string.
//
// Return whether the token is MKR.
func fuckMaker(address string, network string) bool {
	switch network {
	// ethereum
	case chainId.EthereumChainName:
		// whether it is MKR
		return strings.EqualFold(address, "0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2")
	default:
		return false
	}
}

// Return token's total supply amount, already divided by decimals.
func (t *ERC20Info) TotalSupply() (decimal.Decimal, error) {
	if t.Contract == nil {
		return decimal.New(0, 0), errors.New("token must be initiated")
	}
	supply, err := t.Contract.TotalSupply(nil)
	if err != nil {
		return decimal.New(0, 0), err
	}
	totalSupply := types.ToDecimal(supply).Div(decimal.New(1, int32(*t.Decimals)))
	if totalSupply.IsZero() {
		return decimal.New(0, 0), errors.New(*t.Symbol + " total supply is zero")
	}
	return totalSupply, nil
}

// Return the balance, already divided by decimals.
func (t *ERC20Info) BalanceOf(address string) (decimal.Decimal, error) {
	if t.Contract == nil {
		return decimal.New(0, 0), errors.New("token must be initiated")
	}
	balanceBig, err := t.Contract.BalanceOf(nil, types.ToAddress(address))
	if err != nil {
		return decimal.New(0, 0), err
	}
	balance := types.ToDecimal(balanceBig).Div(decimal.New(1, int32(*t.Decimals)))
	return balance, nil
}
