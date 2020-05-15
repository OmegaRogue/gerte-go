package gerte

import "testing"

func TestVersionFromToBytes(t *testing.T) {
	version := Version{
		Major: 1,
		Minor: 1,
		Patch: 0,
	}
	ver := version.ToBytes()
	version2 := VersionFromBytes(ver)

	if version != version2 {
		t.Error("versions don't match")
	}
}
