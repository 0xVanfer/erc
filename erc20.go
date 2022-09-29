package erc

import (
	"errors"
	"math"
	"strings"

	"github.com/0xVanfer/abigen/erc20"
	"github.com/0xVanfer/chainId"
	"github.com/0xVanfer/coingecko"
	"github.com/0xVanfer/types"
	"github.com/ethereum/go-ethereum/accounts/abi/bind"
)

// Basic info of ERC20 token that are stable.
type ERC20Info struct {
	Network  *string
	Address  *string
	Symbol   *string
	Decimals *int
	Contract *erc20.Erc20
}

// Initiate the erc20 token.
func (t *ERC20Info) Init(address string, network string, client bind.ContractBackend) error {
	err := addressRegularCheck(address)
	if err != nil {
		return err
	}
	t.Network = &network
	t.Address = &address
	t.Contract, err = erc20.NewErc20(types.ToAddress(address), client)
	if err != nil {
		return err
	}
	decimals, err := t.Contract.Decimals(nil)
	if err != nil {
		return err
	}
	decimalint := int(decimals.Int64())
	t.Decimals = &decimalint
	var symbol string
	// fuck maker
	if strings.EqualFold(address, "0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2") && strings.EqualFold(network, chainId.EthereumChainName) {
		symbol = "MKR"
	} else {
		symbol, err = t.Contract.Symbol(nil)
		if err != nil {
			return err
		}
	}
	t.Symbol = &symbol
	return nil
}

// Return token's total supply amount, already divided by decimals.
func (t *ERC20Info) TotalSupply() (float64, error) {
	if t.Contract == nil {
		return 0, errors.New("token must be initiated")
	}
	supply, err := t.Contract.TotalSupply(nil)
	if err != nil {
		return 0, err
	}
	totalSupply := types.ToFloat64(supply) * math.Pow10(-*t.Decimals)
	if totalSupply == 0 {
		return 0, errors.New(*t.Symbol + " total supply is zero")
	}
	return totalSupply, nil
}

// Return token's price.
func (t *ERC20Info) PriceUSD(gecko *coingecko.Gecko) (float64, error) {
	if t.Contract == nil {
		return 0, errors.New("token must be initiated")
	}
	price, err := gecko.GetPriceBySymbol(*t.Symbol, *t.Network, "usd")
	if err != nil {
		return 0, err
	}
	if price == 0 {
		return 0, errors.New(*t.Symbol + "price is zero")
	}
	return price, nil
}

// Return token's total supply in usd.
func (t *ERC20Info) TotalSupplyUSD(gecko *coingecko.Gecko) (float64, error) {
	if t.Contract == nil {
		return 0, errors.New("token must be initiated")
	}
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
	if t.Contract == nil {
		return 0, errors.New("token must be initiated")
	}
	balanceBig, err := t.Contract.BalanceOf(nil, types.ToAddress(address))
	if err != nil {
		return 0, err
	}
	balance := types.ToFloat64(balanceBig) * math.Pow10(-*t.Decimals)
	return balance, nil
}
