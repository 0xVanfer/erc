package erc

import (
	"errors"

	"github.com/0xVanfer/abigen/erc20"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
	"github.com/shopspring/decimal"
)

// Create a new ERC20 token.
func NewErc20(address string, network string, client bind.ContractBackend) (*ERC20Info, error) {
	// Address must be length 42 and not 0x0.
	err := addressRegularCheck(address)
	if err != nil {
		return nil, err
	}
	var new ERC20Info
	new.Network = &network
	// Use checksumed address.
	checksumed := checksumEthereumAddress(address)
	new.Address = &checksumed
	// Contract.
	new.Contract, err = erc20.NewErc20(types.ToAddress(address), client)
	if err != nil {
		return nil, err
	}
	// Decimals.
	decimals, err := new.Contract.Decimals(nil)
	if err != nil {
		return nil, err
	}
	decimalint := int(decimals.Int64())
	new.Decimals = &decimalint
	// Symbol.
	var symbol string
	// Maker symbol must be handled seperatedly.
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

// Return token's total supply amount, already divided by decimals.
func (t *ERC20Info) TotalSupply() (decimal.Decimal, error) {
	// Token must initiated.
	if t.Contract == nil {
		return decimal.New(0, 0), errors.New("token must be initiated")
	}
	// Get total supply.
	supply, err := t.Contract.TotalSupply(nil)
	if err != nil {
		return decimal.New(0, 0), err
	}
	// Divide by decimals and turn into decimal.Decimal.
	totalSupply := types.ToDecimal(supply).Div(decimal.New(1, int32(*t.Decimals)))
	// Read total supply succeeded, but is zero.
	if totalSupply.IsZero() {
		return decimal.New(0, 0), errors.New(*t.Symbol + " total supply is zero")
	}
	return totalSupply, nil
}

// Return the balance, already divided by decimals.
func (t *ERC20Info) BalanceOf(address string) (decimal.Decimal, error) {
	// Token must initiated.
	if t.Contract == nil {
		return decimal.New(0, 0), errors.New("token must be initiated")
	}
	// Read the balance.
	balanceBig, err := t.Contract.BalanceOf(nil, types.ToAddress(address))
	if err != nil {
		return decimal.New(0, 0), err
	}
	// Divide by decimals and turn into decimal.Decimal.
	balance := types.ToDecimal(balanceBig).Div(decimal.New(1, int32(*t.Decimals)))
	return balance, nil
}

// Return token's history total supply amount at `blockNumber`, already divided by decimals.
func (t *ERC20Info) HistoryTotalSupply(blockNumber int64) (decimal.Decimal, error) {
	// Token must initiated.
	if t.Contract == nil {
		return decimal.New(0, 0), errors.New("token must be initiated")
	}
	// Get total supply.
	supply, err := t.Contract.TotalSupply(&bind.CallOpts{BlockNumber: types.ToBigInt(blockNumber)})
	if err != nil {
		return decimal.New(0, 0), err
	}
	// Divide by decimals and turn into decimal.Decimal.
	totalSupply := types.ToDecimal(supply).Div(decimal.New(1, int32(*t.Decimals)))
	// Read total supply succeeded, but is zero.
	if totalSupply.IsZero() {
		return decimal.New(0, 0), errors.New(*t.Symbol + " total supply is zero")
	}
	return totalSupply, nil
}

// Return the balance, already divided by decimals.
func (t *ERC20Info) HistoryBalanceOf(address string, blockNumber int64) (decimal.Decimal, error) {
	// Token must initiated.
	if t.Contract == nil {
		return decimal.New(0, 0), errors.New("token must be initiated")
	}
	// Read the balance.
	balanceBig, err := t.Contract.BalanceOf(&bind.CallOpts{BlockNumber: types.ToBigInt(blockNumber)}, types.ToAddress(address))
	if err != nil {
		return decimal.New(0, 0), err
	}
	// Divide by decimals and turn into decimal.Decimal.
	balance := types.ToDecimal(balanceBig).Div(decimal.New(1, int32(*t.Decimals)))
	return balance, nil
}
