package encoding

import (
	"encoding/hex"
	"strings"
)

// HexToBin converts a string containing hexadecimal digits and whitespace to its binary representation. Whitespace is
// ignored. It assumes hexStr represents a correct hexadecimal string.
func HexToBin(hexStr string) []byte {
	b := []byte(strings.ReplaceAll(hexStr, " ", ""))
	in := make([]byte, hex.DecodedLen(len(b)))
	_, _ = hex.Decode(in, b)
	return in
}
