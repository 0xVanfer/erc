package erc

import (
	"strings"

	"github.com/0xVanfer/chainId"
)

// Maker symbol is bytes instead of string.
//
// Return whether the token is MKR.
func fuckMaker(address string, network string) bool {
	// Avalanche: use bridged MKR, called MKR.e.
	// Polygon: use child proxy.
	switch network {
	// ethereum
	case chainId.EthereumChainName:
		// Whether it is MKR.
		return strings.EqualFold(address, "0x9f8f72aa9304c8b593d555f12ef6589cc3a579a2")
	default:
		return false
	}
}
