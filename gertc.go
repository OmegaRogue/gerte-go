package gerte

import "fmt"

// GERTc is a 6 byte GERTc Address
type GERTc struct {
	GERTe GertAddress
	GERTi GertAddress
}

// GertCFromBytes parses bytes to a GERTc Address
func GertCFromBytes(data []byte) GERTc {
	return GERTc{
		GERTe: GertAddress{
			Upper: (int(data[0]) << 4) | (int(data[1]) >> 4),
			Lower: ((int(data[1]) & 0x0F) << 8) | int(data[2]),
		},
		GERTi: GertAddress{
			Upper: (int(data[3]) << 4) | (int(data[4]) >> 4),
			Lower: ((int(data[4]) & 0x0F) << 8) | int(data[5]),
		},
	}

}

// ToBytes converts a GERTc to bytes for sending
func (addr GERTc) ToBytes() []byte {
	return append(addr.GERTe.ToBytes(), addr.GERTi.ToBytes()...)
}

// String prints a GERTc Address as a string
func (addr GERTc) String() string {
	return fmt.Sprintf("%04v.%04v:%04v.%04v", addr.GERTe.Upper, addr.GERTe.Lower, addr.GERTi.Upper, addr.GERTi.Lower)
}

// GoString prints a GERTc Address as a string surrounded with brackets
func (addr GERTc) GoString() string {
	return fmt.Sprintf("[%v]", addr)
}
