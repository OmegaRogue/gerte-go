package gerte

import (
	"fmt"
	"strconv"
	"strings"
)

// GertAddress is a 3 byte address used as a GERTe/i Address
type GertAddress struct {
	Upper int
	Lower int
}

// ToBytes converts a GERT Address to bytes for sending
func (addr GertAddress) ToBytes() []byte {
	var b strings.Builder
	b.WriteByte(byte(addr.Upper >> 4))
	b.WriteByte(byte(((addr.Upper & 0x0F) << 4) | (addr.Lower >> 8)))
	b.WriteByte(byte(addr.Lower & 0xFF))
	return []byte(b.String())
}

// AddressFromBytes parses bytes to a GERT Address
func AddressFromBytes(data []byte) GertAddress {
	return GertAddress{
		Upper: (int(data[0]) << 4) | (int(data[1]) >> 4),
		Lower: ((int(data[1]) & 0x0F) << 8) | int(data[2]),
	}
}

// String prints a GertAddress as a string
func (addr GertAddress) String() string {
	return fmt.Sprintf("%04v.%04v", addr.Upper, addr.Lower)
}

// GoString prints a GertAddress as a string surrounded with brackets
func (addr GertAddress) GoString() string {
	return fmt.Sprintf("[%v]", addr)
}

// AddressFromString converts a string with an address in the format "XXXX.YYYY" into the corresponding GertAddress.
// It returns the GertAddress and any encountered errors.
func AddressFromString(addr string) (GertAddress, error) {
	parts := strings.Split(addr, ".")

	upper, err := strconv.ParseInt(parts[0], 10, 0)
	if err != nil {
		return GertAddress{}, fmt.Errorf("error on parse upper String: %w", err)
	}
	lower, err := strconv.ParseInt(parts[1], 10, 0)
	if err != nil {
		return GertAddress{}, fmt.Errorf("error on parse lower String: %w", err)
	}
	return GertAddress{
		Upper: int(upper),
		Lower: int(lower),
	}, nil
}
