package gerte

import "fmt"

// Version is the Version used by the Connected Status
type Version struct {
	Major byte
	Minor byte
	Patch byte
}

// ToBytes converts a GERT Version to bytes for sending
func (ver Version) ToBytes() []byte {
	return []byte{ver.Major, ver.Minor}
}

// VersionFromBytes parses bytes to a GERT Version
func VersionFromBytes(b []byte) Version {
	return Version{
		Major: b[0],
		Minor: b[1],
	}
}

// String prints a GERT Version to a Human-readable string
func (ver Version) String() string {
	return fmt.Sprintf("%v.%v.%v", ver.Major, ver.Minor, ver.Patch)
}

// GoString prints a GERT Version to a Human-readable string surrounded with brackets
func (ver Version) GoString() string {
	return fmt.Sprintf("[%v]", ver)
}
