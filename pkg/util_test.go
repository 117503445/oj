package pkg

import (
	"testing"
)

func TestGetFileNameWithoutExt(t *testing.T) {
	s := GetFileNameWithoutExt("1.2.txt")
	if s != "1.2" {
		t.Fail()
	}
}
