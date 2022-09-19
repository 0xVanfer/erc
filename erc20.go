package erc

import (
	"errors"
	"math"

	"github.com/0xVanfer/abigen/erc20"
	"github.com/0xVanfer/coingecko"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Basic info of ERC20 token that are stable.
type ERC20Info struct {
	Network  string
	Address  string
	Symbol   string
	Decimals int
	Contract *erc20.Erc20
}

func (t *ERC20Info) Init(address string, network string, client bind.ContractBackend) error {
	err := addressRegularCheck(address)
	if err != nil {
		return err
	}
	token, err := erc20.NewErc20(types.ToAddress(address), client)
	if err != nil {
		return err
	}
	decimals, err := token.Decimals(nil)
	if err != nil {
		return err
	}
	symbol, err := token.Symbol(nil)
	if err != nil {
		return err
	}

	t = &ERC20Info{
		Network:  network,
		Address:  address,
		Symbol:   symbol,
		Decimals: int(decimals.Int64()),
		Contract: token,
	}
	return nil
}

// Return token's total supply amount, already divided by decimals.
func (t *ERC20Info) TotalSupply() (float64, error) {
	supply, err := t.Contract.TotalSupply(nil)
	if err != nil {
		return 0, err
	}
	totalSupply := types.ToFloat64(supply) * math.Pow10(-t.Decimals)
	if totalSupply == 0 {
		return 0, errors.New(t.Symbol + " total supply is zero")
	}
	return totalSupply, nil
}

// Return token's price.
func (t *ERC20Info) PriceUSD(gecko *coingecko.Gecko) (float64, error) {
	price, err := gecko.GetPriceBySymbol(t.Symbol, t.Network, "usd")
	if err != nil {
		return 0, err
	}
	if price == 0 {
		return 0, errors.New(t.Symbol + "price is zero")
	}
	return price, nil
}

// Return token's total supply in usd.
func (t *ERC20Info) TotalSupplyUSD(gecko *coingecko.Gecko) (float64, error) {
	supply, err := t.TotalSupply()
	if err != nil {
		return 0, err
	}
	price, err := t.PriceUSD(gecko)
	if err != nil {
		return 0, err
	}
	supplyUSD := supply * price
	return supplyUSD, nil
}

// Return the balance, already divided by decimals.
func (t *ERC20Info) BalanceOf(address string) (float64, error) {
	balanceBig, err := t.Contract.BalanceOf(nil, types.ToAddress(address))
	if err != nil {
		return 0, err
	}
	balance := types.ToFloat64(balanceBig) * math.Pow10(-t.Decimals)
	return balance, nil
}
