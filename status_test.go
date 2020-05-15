package gerte

import "testing"

func TestStatusFromBytes(t *testing.T) {
	st := []byte{byte(StateConnected)}
	state := Status{
		Status: StateConnected,
		Size:   4,
		Version: Version{
			Major: 1,
			Minor: 1,
			Patch: 0,
		},
	}
	verByte := Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}.ToBytes()
	st = append(st, verByte[0], verByte[1])
	state2, err := StatusFromBytes(st)
	if err != nil {
		t.Errorf("error on unmarshal command: %+v", err)
	}
	if state != state2 {
		t.Errorf("commands don't match:\n%+v\n%+v", state, state2)
	}
}
