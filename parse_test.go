package gerte

import "testing"

func TestAddressFromString(t *testing.T) {
	addr, err := AddressFromString("0123.0456")
	if err != nil {
		t.Errorf("error on parse address string: %+v", err)
	}
	addrT := GertAddress{
		Upper: 123,
		Lower: 456,
	}
	if addr.Upper != addrT.Upper || addr.Lower != addrT.Lower {
		t.Error("addresses don't match")
	}
}
