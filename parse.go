package gerte

import (
	"fmt"
	"strconv"
	"strings"
)

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
