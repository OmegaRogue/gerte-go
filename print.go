package gerte

import (
	"fmt"
	"strings"
)

// TODO Make Doc
func (cmd Command) PrintCommand() (string, error) {
	output := ""
	switch cmd.Command {
	case CommandState:
		output += "[STATE]"
		switch cmd.Status.Status {
		case StateFailure:
			output += "[FAILURE]"
			switch cmd.Status.Error {
			case ErrorVersion:
				output += "[VERSION]"
				break
			case ErrorBadKey:
				output += "[BAD_KEY]"
				break
			case ErrorAlreadyRegistered:
				output += "[ALREADY_REGISTERED]"
				break
			case ErrorNotRegistered:
				output += "[NOT_REGISTERED]"
				break
			case ErrorNoRoute:
				output += "[NO_ROUTE]"
				break
			case ErrorAddressTaken:
				output += "[ADDRESS_TAKEN]"
				break
			}
			break
		case StateConnected:
			output += "[CONNECTED]"
			output += "[" + cmd.Status.Version.PrintVersion() + "]"
			break
		case StateAssigned:
			output += "[ASSIGNED]"
			break
		case StateClosed:
			output += "[CLOSED]"
			break
		case StateSent:
			output += "[SENT]"
			break
		}
		break
	case CommandClose:
		output += "[CLOSE]"
		break
	case CommandData:
		output += "[DATA]"
		output += cmd.Packet.PrintPacket()
		break
	default:
		return "", fmt.Errorf("no valid cmd: %v", cmd.Command)
	}
	return output, nil
}

// TODO Make Doc
func (ver Version) PrintVersion() string {
	return fmt.Sprintf("%v.%v.%v", ver.Major, ver.Minor, ver.Patch)
}

// PrintAddress prints a GertAddress as a string
func (addr GertAddress) PrintAddress() string {
	return fmt.Sprintf("%04v.%04v", addr.Upper, addr.Lower)
}

// PrintGERTc prints a GERTc Address as a string
func (addr GERTc) PrintGERTc() string {
	return fmt.Sprintf("%04v.%04v:%04v.%04v", addr.GERTe.Upper, addr.GERTe.Lower, addr.GERTi.Upper, addr.GERTi.Lower)
}
func (error GertError) PrintError() (string, error) {
	switch error {
	case ErrorVersion:
		return "[VERSION]", nil
	case ErrorBadKey:
		return "[BAD_KEY]", nil
	case ErrorAlreadyRegistered:
		return "[ALREADY_REGISTERED]", nil
	case ErrorNotRegistered:
		return "[NOT_REGISTERED]", nil
	case ErrorNoRoute:
		return "[NO_ROUTE]", nil
	case ErrorAddressTaken:
		return "[ADDRESS_TAKEN]", nil
	}
	return "", fmt.Errorf("invalid Failure")
}

// TODO Make Doc
func (status Status) PrintStatus() (string, error) {
	switch status.Status {
	case StateFailure:
		p, err := status.Error.PrintError()
		if err != nil {
			return "", fmt.Errorf("error while parsing Failure: %+v", err)
		}
		return "[FAILURE]" + p, nil
	case StateConnected:
		return "[CONNECTED][" + status.Version.PrintVersion() + "]", nil
	case StateAssigned:
		return "[ASSIGNED]", nil
	case StateClosed:
		return "[CLOSED]", nil
	case StateSent:
		return "[SENT]", nil
	}
	return "", fmt.Errorf("invalid status")
}

// TODO Make Doc
func (pkt Packet) PrintPacket() string {

	fmt.Println(string(pkt.Data))
	return fmt.Sprintf("[%v][%v][%v][%v]", pkt.Source.PrintGERTc(), pkt.Target.PrintGERTc(), len(pkt.Data), strings.Trim(string(pkt.Data), " "))
}

// TODO Make Doc
func PrettyPrint(data []byte) (string, error) {
	output := ""
	switch data[0] {
	case byte(CommandState):
		state, err := StatusFromBytes(data[1:])
		if err != nil {
			return "", fmt.Errorf("error while parsing status data: %+v", err)
		}
		p, err := state.PrintStatus()
		if err != nil {
			return "", fmt.Errorf("error while printing status: %+v", err)
		}
		output += "[STATE]" + p
		break
	case byte(CommandRegister):
		output += "[REGISTER]"
		addr := AddressFromBytes(data[1:4])
		output += "[" + addr.PrintAddress() + "]"
		key := string(data[4:24])
		output += "[" + key + "]"
		break
	case byte(CommandData):
		source := GertCFromBytes(data[1:7])
		target := AddressFromBytes(data[7:10])
		length := data[10]
		dat := data[11 : 11+length]
		output += "[DATA]"
		output += fmt.Sprintf("[%v][%v][%v][%v]", source.PrintGERTc(), target.PrintAddress(), length, string(dat))
		break
	case byte(CommandClose):
		output += "[CLOSE]"
		break
	default:
		return "", fmt.Errorf("no valid command: %v", data[0])
	}
	return output, nil
}
