package gerte

import "fmt"

// GertError is the Error Code in a "Failed" Status
type GertError byte

const (
	// ErrorVersion indicates that an incompatible version was used during negotiation.
	ErrorVersion GertError = iota
	// ErrorBadKey indicates that the Key did not match that used for the requested address.
	// Requested address may not exist
	ErrorBadKey
	// ErrorAlreadyRegistered indicates that Registration has already been performed successfully
	ErrorAlreadyRegistered
	// ErrorNotRegistered indicates that the Gateway hasn't been registered yet.
	// The gateway cannot send data before claiming an address
	ErrorNotRegistered
	// ErrorNoRoute indicates that Data failed to send because the remote gateway couldn't be found
	ErrorNoRoute
	// ErrorAddressTaken indicates that the Address request has already been claimed
	ErrorAddressTaken
)

// PrintError prints a GERT Error to a Human-readable string
func (error GertError) String() string {
	switch error {
	case ErrorVersion:
		return "VERSION"
	case ErrorBadKey:
		return "BAD_KEY"
	case ErrorAlreadyRegistered:
		return "ALREADY_REGISTERED"
	case ErrorNotRegistered:
		return "NOT_REGISTERED"
	case ErrorNoRoute:
		return "NO_ROUTE"
	case ErrorAddressTaken:
		return "ADDRESS_TAKEN"
	}
	return "nil"
}

// PrintError prints a GERT Error to a Human-readable string
func (error GertError) GoString() string {
	return fmt.Sprintf("[%v]", error)
}
