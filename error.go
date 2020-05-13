package gerte

import "fmt"

func (status Status) parseError() error {
	switch status.Error {
	case ErrorVersion:
		return fmt.Errorf("incompatible version during negotiation: %v", status.Version)
	case ErrorBadKey:
		return fmt.Errorf("key did not match that used for the requested address. Requested address may not exist")
	case ErrorAlreadyRegistered:
		return fmt.Errorf("registration has already been performed successfully")
	case ErrorNotRegistered:
		return fmt.Errorf("gateway cannot send data before claiming an address")
	case ErrorNoRoute:
		return fmt.Errorf("data failed to send because remote gateway could not be found")
	case ErrorAddressTaken:
		return fmt.Errorf("address request has already been claimed")
	}
	return fmt.Errorf("no valid error")
}
