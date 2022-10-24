package erc

import (
	"strings"

	"golang.org/x/crypto/sha3"
)

// Return checksummed ethereum address.
func checksumEthereumAddress(addr string) string {
	hex := strings.ToLower(addr)[2:]
	d := sha3.NewLegacyKeccak256()
	d.Write([]byte(hex))
	hash := d.Sum(nil)
	checksumed := "0x"
	for i, b := range hex {
		c := string(b)
		if b < '0' || b > '9' {
			if hash[i/2]&byte(128-i%2*120) != 0 {
				c = string(b - 32)
			}
		}
		checksumed += c
	}
	return checksumed
}
