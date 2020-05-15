package gerte

import "testing"

func TestGertCFromToBytes(t *testing.T) {
	address := GERTc{
		GERTe: GertAddress{
			Upper: 0,
			Lower: 0,
		},
		GERTi: GertAddress{
			Upper: 123,
			Lower: 456,
		},
	}

	addr := address.ToBytes()
	address2 := GertCFromBytes(addr)

	if address != address2 {
		t.Error("addresses don't match")
	}
}
