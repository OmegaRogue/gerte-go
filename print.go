package gerte

import (
	"fmt"
)

// PrettyPrint prints a GERT Message into a human readable string
func PrettyPrint(data []byte) (string, error) {

	switch data[0] {
	case byte(CommandState):
		state, err := StatusFromBytes(data[1:])
		if err != nil {
			return "", fmt.Errorf("error while parsing status data: %+v", err)
		}

		return fmt.Sprintf("%#v%#v", CommandState, state), nil
	case byte(CommandRegister):
		addr := AddressFromBytes(data[1:4])
		key := string(data[4:24])
		return fmt.Sprintf("%#v%#v[%v]", CommandRegister, addr, key), nil
	case byte(CommandData):
		source := GertCFromBytes(data[1:7])
		target := AddressFromBytes(data[7:10])
		length := data[10]
		dat := data[11 : 11+length]
		return fmt.Sprintf("[DATA]%#v%#v[%v][%v]", source, target, length, string(dat)), nil
	case byte(CommandClose):
		return fmt.Sprintf("%#v", CommandClose), nil

	}
	return "[nil]", fmt.Errorf("no valid command: %v", data[0])
}
