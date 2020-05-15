package gerte

import "fmt"

type (
	// GertStatus indicates the Status from a Status Command
	GertStatus byte

	// Status is the Parsed Status returned from Api.Parse for a Status Command
	Status struct {
		Status  GertStatus
		Size    byte
		Error   GertError
		Version Version
	}
)

const (
	// StateFailure is the initial gateway state.
	// Should be changed upon negotiation.
	// Also a response to failed commands with an error
	StateFailure GertStatus = iota
	// StateConnected indicates that the Gateway is connected to a relay, no other action has been taken
	StateConnected
	// StateAssigned indicates that the Gateway has successfully claimed an address.
	// Used as a response to the REGISTER command.
	StateAssigned
	// StateClosed indicates that the Gateway has closed the connection.
	// Used as a response to the CLOSE command.
	StateClosed
	// StateSent indicates that the Data has been successfully sent.
	// Used as a response to the DATA command.
	// This is not a guarantee for data that has to be sent to another peer, although it's unlikely to be incorrect.
	StateSent
)

// String prints a GertStatus to a Human-readable string
func (state GertStatus) String() string {
	switch state {
	case StateFailure:
		return "FAILURE"
	case StateConnected:
		return "CONNECTED"
	case StateAssigned:
		return "ASSIGNED"
	case StateClosed:
		return "CLOSED"
	case StateSent:
		return "SENT"
	}
	return "nil"
}

// GoString prints a GertStatus to a Human-readable string with surrounding Brackets
func (state GertStatus) GoString() string {
	return fmt.Sprintf("[%v]", state)
}

// StatusFromBytes parses bytes to a GERT Status
func StatusFromBytes(data []byte) (Status, error) {

	switch data[0] {
	case byte(StateFailure):
		return Status{
			Status: StateFailure,
			Size:   2,
			Error:  GertError(data[1]),
		}, nil
	case byte(StateConnected):
		if len(data) < 3 {
			return Status{}, fmt.Errorf("data too short: %v<3", len(data))
		}
		return Status{
			Status: StateConnected,
			Size:   4,
			Version: Version{
				Major: data[1],
				Minor: data[2],
			},
		}, nil
	case byte(StateAssigned):
		return Status{
			Status: StateAssigned,
			Size:   1,
		}, nil
	case byte(StateClosed):
		return Status{
			Status: StateClosed,
			Size:   1,
		}, nil
	case byte(StateSent):
		return Status{
			Status: StateSent,
			Size:   1,
		}, nil
	}
	return Status{}, fmt.Errorf("state didn't match any known state: %v", data[0])
}

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

// String prints a GERT Status to a Human-readable string
func (status Status) String() string {
	switch status.Status {
	case StateFailure:
		return fmt.Sprintf("%v %v", status.Status, status.Error)
	case StateConnected:
		return fmt.Sprintf("%v %v", status.Status, status.Version)
	case StateAssigned:
		return fmt.Sprintf("%v", status.Status)
	case StateClosed:
		return fmt.Sprintf("%v", status.Status)
	case StateSent:
		return fmt.Sprintf("%v", status.Status)
	}
	return "nil"
}

// GoString prints a GERT Status to a Human-readable string surrounded with brackets
func (status Status) GoString() string {
	switch status.Status {
	case StateFailure:
		return fmt.Sprintf("%#v%#v", status.Status, status.Error)
	case StateConnected:
		return fmt.Sprintf("%#v%#v", status.Status, status.Version)
	case StateAssigned:
		return fmt.Sprintf("%#v", status.Status)
	case StateClosed:
		return fmt.Sprintf("%#v", status.Status)
	case StateSent:
		return fmt.Sprintf("%#v", status.Status)
	}
	return "[nil]"
}
