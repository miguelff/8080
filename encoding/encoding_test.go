package encoding

import (
	"fmt"
	"strings"
	"testing"
)

func TestHexToBin(t *testing.T) {
	str := "00 D0 C1"
	bin := HexToBin(str)

	want := strings.ReplaceAll(str, " ", "")
	got := fmt.Sprintf("%X", bin)
	if got != want {
		t.Errorf("got %s \n\n want %s", got, want)
	}
}
