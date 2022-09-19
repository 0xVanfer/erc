package erc

import "errors"

func addressRegularCheck(address string) error {
	if address == "" {
		return errors.New("address should not be empty")
	}
	if (address == "0x0") || (address == "0x0000000000000000000000000000000000000000") {
		return errors.New("address should not be zero address")
	}
	if len(address) != 42 {
		return errors.New("address length must be 42")
	}
	return nil
}
