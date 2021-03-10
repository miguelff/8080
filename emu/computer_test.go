package emu

import (
	"bytes"
	"testing"

	"github.com/miguelff/8080/encoding"
)

func ram(bytes string) []byte {
	return encoding.HexToBin(bytes)
}
func TestParity(t *testing.T) {
	for _, tC := range []struct {
		b    byte
		want Flags
	}{
		{
			0b00001111,
			pf,
		},
		{
			0b00001110,
			none,
		},
	} {
		if got := parity8(tC.b); got != tC.want {
			t.Errorf("got %b, want %b", got, tC.want)
		}
	}
}
func TestSign(t *testing.T) {
	for _, tC := range []struct {
		b    byte
		want Flags
	}{
		{
			0b10001111,
			sf,
		},
		{
			0b00001110,
			none,
		},
	} {
		if got := sign8(tC.b); got != tC.want {
			t.Errorf("got %b, want %b", got, tC.want)
		}
	}
}
func TestZero(t *testing.T) {
	for _, tC := range []struct {
		b    byte
		want Flags
	}{
		{
			0b10001111,
			none,
		},
		{
			0b00000000,
			zf,
		},
	} {
		if got := zero8(tC.b); got != tC.want {
			t.Errorf("got %b, want %b", got, tC.want)
		}
	}
}
func TestComputer_Step(t *testing.T) {
	for _, tC := range []struct {
		desc string
		init *Computer
		want *Computer
	}{
		{
			"ADC A: with carry",
			newComputer(
				CPU{
					A:     0x02,
					Flags: cf,
				},
				ram("8F"),
			),
			newComputer(
				CPU{
					A:     0x05,
					PC:    0x01,
					Flags: pf,
				},
				ram("8F"),
			),
		},
		{
			"ADC A: no carry",
			newComputer(
				CPU{
					A:     0x02,
					Flags: none,
				},
				ram("8F"),
			),
			newComputer(
				CPU{
					A:     0x04,
					PC:    0x01,
					Flags: none,
				},
				ram("8F"),
			),
		},
		{
			"ADC B",
			newComputer(
				CPU{
					A:     0x02,
					B:     0x01,
					Flags: cf,
				},
				ram("88"),
			),
			newComputer(
				CPU{
					A:     0x04,
					B:     0x01,
					PC:    1,
					Flags: none,
				},
				ram("88"),
			),
		},
		{
			"ADC C",
			newComputer(
				CPU{
					A:     0x01,
					C:     0xFD,
					Flags: cf,
				},
				ram("89"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					C:     0xFD,
					PC:    0x01,
					Flags: sf | pf,
				},
				ram("89"),
			),
		},
		{
			"ADC D",
			newComputer(
				CPU{
					A:     0x00,
					D:     0xFD,
					Flags: cf,
				},
				ram("8A"),
			),
			newComputer(
				CPU{
					A:     0xFE,
					D:     0xFD,
					PC:    0x01,
					Flags: sf,
				},
				ram("8A"),
			),
		},
		{
			"ADC E",
			newComputer(
				CPU{
					A:     0x00,
					E:     0xFD,
					Flags: cf,
				},
				ram("8B"),
			),
			newComputer(
				CPU{
					A:     0xFE,
					E:     0xFD,
					PC:    0x01,
					Flags: sf,
				},
				ram("8B"),
			),
		},
		{
			"ADC H",
			newComputer(
				CPU{
					A:     0x00,
					H:     0xFD,
					Flags: cf,
				},
				ram("8C"),
			),
			newComputer(
				CPU{
					A:     0xFE,
					H:     0xFD,
					PC:    0x01,
					Flags: sf,
				},
				ram("8C"),
			),
		},
		{
			"ADC L",
			newComputer(
				CPU{
					A:     0x00,
					L:     0xFD,
					Flags: cf,
				},
				ram("8D"),
			),
			newComputer(
				CPU{
					A:     0xFE,
					L:     0xFD,
					PC:    0x01,
					Flags: sf,
				},
				ram("8D"),
			),
		},
		{
			"ADC M",
			newComputer(
				CPU{
					A:     0x01,
					H:     0x00,
					L:     0x02,
					Flags: cf,
				},
				ram("8E 00 FD"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					H:     0x00,
					L:     0x02,
					PC:    0x01,
					Flags: sf | pf,
				},
				ram("8E 00 FD"),
			),
		},
		{
			"ADD A",
			newComputer(
				CPU{
					A: 0x02,
				},
				ram("87"),
			),
			newComputer(
				CPU{
					A:     0x04,
					PC:    0x01,
					Flags: none,
				},
				ram("87"),
			),
		},
		{
			"ADD B: adding two values that sum 0x0 sets the zero8 and parity8 Flags",
			newComputer(
				CPU{},
				ram("80"),
			),
			newComputer(
				CPU{
					PC:    1,
					Flags: zf | pf,
				},
				ram("80"),
			),
		},
		{
			"ADD B: adding 0x09 + 0x07 sets the auxiliary carry flag",
			newComputer(
				CPU{
					A: 0x09,
					B: 0x07,
				},
				ram("80"),
			),
			newComputer(
				CPU{
					A:     0x10,
					B:     0x07,
					PC:    0x01,
					Flags: acf,
				},
				ram("80"),
			),
		},
		{
			"ADD B: adding 0xFE and 0x02 sets the carry and auxiliary carry Flags",
			newComputer(
				CPU{
					A: 0x03,
					B: 0xFE,
				},
				ram("80"),
			),
			newComputer(
				CPU{
					A:     0x01,
					B:     0xFE,
					PC:    0x01,
					Flags: cf | acf,
				},
				ram("80"),
			),
		},
		{
			"ADD C: adding 0xFE and 0x01 set the parity8 and sign8 Flags",
			newComputer(
				CPU{
					A: 0x01,
					C: 0xFE,
				},
				ram("81"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					C:     0xFE,
					PC:    0x01,
					Flags: sf | pf,
				},
				ram("81"),
			),
		},
		{
			"ADD D: adding 0xFf and 0x00 set the sign8 flag",
			newComputer(
				CPU{
					A: 0x00,
					D: 0xFE,
				},
				ram("82"),
			),
			newComputer(
				CPU{
					A:     0xFE,
					D:     0xFE,
					PC:    0x01,
					Flags: sf,
				},
				ram("82"),
			),
		},
		{
			"ADD E",
			newComputer(
				CPU{
					A: 0x00,
					E: 0xFE,
				},
				ram("83"),
			),
			newComputer(
				CPU{
					A:     0xFE,
					E:     0xFE,
					PC:    0x01,
					Flags: sf,
				},
				ram("83"),
			),
		},
		{
			"ADD H",
			newComputer(
				CPU{
					A: 0x00,
					H: 0xFE,
				},
				ram("84"),
			),
			newComputer(
				CPU{
					A:     0xFE,
					H:     0xFE,
					PC:    0x01,
					Flags: sf,
				},
				ram("84"),
			),
		},
		{
			"ADD L",
			newComputer(
				CPU{
					A: 0x00,
					L: 0xFE,
				},
				ram("85"),
			),
			newComputer(
				CPU{
					A:     0xFE,
					L:     0xFE,
					PC:    0x01,
					Flags: sf,
				},
				ram("85"),
			),
		},
		{
			"ADD M",
			newComputer(
				CPU{
					A: 0x01,
					H: 0x00,
					L: 0x02,
				},
				ram("86 00 FE"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					H:     0x00,
					L:     0x02,
					PC:    0x01,
					Flags: sf | pf,
				},
				ram("86 00 FE"),
			),
		},
		{
			"ANA A",
			newComputer(
				CPU{
					A:     0xFF,
					Flags: cf,
				},
				ram("A7"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("A7"),
			),
		},
		{
			"ANA B",
			newComputer(
				CPU{
					A:     0xFF,
					B:     0x0A,
					Flags: cf,
				},
				ram("A0"),
			),
			newComputer(
				CPU{
					A:     0x0A,
					B:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				ram("A0"),
			),
		},
		{
			"ANA C",
			newComputer(
				CPU{
					A:     0xFF,
					C:     0x0A,
					Flags: cf,
				},
				ram("A1"),
			),
			newComputer(
				CPU{
					A:     0x0A,
					C:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				ram("A1"),
			),
		},
		{
			"ANA D",
			newComputer(
				CPU{
					A:     0xFF,
					D:     0x0A,
					Flags: cf,
				},
				ram("A2"),
			),
			newComputer(
				CPU{
					A:     0x0A,
					D:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				ram("A2"),
			),
		},
		{
			"ANA E",
			newComputer(
				CPU{
					A:     0xFF,
					E:     0x0A,
					Flags: cf,
				},
				ram("A3"),
			),
			newComputer(
				CPU{
					A:     0x0A,
					E:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				ram("A3"),
			),
		},
		{
			"ANA H",
			newComputer(
				CPU{
					A:     0xFF,
					H:     0x0A,
					Flags: cf,
				},
				ram("A4"),
			),
			newComputer(
				CPU{
					A:     0x0A,
					H:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				ram("A4"),
			),
		},
		{
			"ANA L",
			newComputer(
				CPU{
					A:     0xFF,
					L:     0x0A,
					Flags: cf,
				},
				ram("A5"),
			),
			newComputer(
				CPU{
					A:     0x0A,
					L:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				ram("A5"),
			),
		},
		{
			"CALL adr",
			newComputer(
				CPU{
					PC: 0x01,
					SP: 0x07,
				},
				ram("00 CD 0A 00 00 00 00 00"),
			),
			newComputer(
				CPU{
					PC: 0x0A,
					SP: 0x05,
				},
				ram("00 CD 0A 00 00 00 01 00"),
			),
		},
		{
			"CMP A",
			newComputer(
				CPU{
					A: 0x30,
				},
				ram("BF"),
			),
			newComputer(
				CPU{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("BF"),
			),
		},
		{
			"CMP B: generates carry",
			newComputer(
				CPU{
					A: 0x30,
					B: 0x31,
				},
				ram("B8"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					B:     0x31,
					PC:    0x01,
					Flags: cf | pf | sf,
				},
				ram("B8"),
			),
		}, {
			"CMP B",
			newComputer(
				CPU{
					A: 0x30,
					B: 0x01,
				},
				ram("B8"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					B:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("B8"),
			),
		},
		{
			"CMP C",
			newComputer(
				CPU{
					A: 0x30,
					C: 0x01,
				},
				ram("B9"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					C:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("B9"),
			),
		},
		{
			"CMP D",
			newComputer(
				CPU{
					A: 0x30,
					D: 0x01,
				},
				ram("BA"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					D:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("BA"),
			),
		},
		{
			"CMP E",
			newComputer(
				CPU{
					A: 0x30,
					E: 0x01,
				},
				ram("BB"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					E:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("BB"),
			),
		},
		{
			"CMP H",
			newComputer(
				CPU{
					A: 0x30,
					H: 0x01,
				},
				ram("BC"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					H:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("BC"),
			),
		},
		{
			"CMP L",
			newComputer(
				CPU{
					A: 0x30,
					L: 0x01,
				},
				ram("BD"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					L:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("BD"),
			),
		},
		{
			"DAD B",
			newComputer(
				CPU{
					B: 0x01,
					C: 0x01,
					H: 0x01,
					L: 0x01,
				},
				ram("09"),
			),
			newComputer(
				CPU{
					B:     0x01,
					C:     0x01,
					H:     0x02,
					L:     0x02,
					PC:    0x01,
					Flags: none,
				},
				ram("09"),
			),
		},
		{
			"DAD B: generates carry",
			newComputer(
				CPU{
					B: 0xFF,
					C: 0xFE,
					H: 0x00,
					L: 0x03,
				},
				ram("09"),
			),
			newComputer(
				CPU{
					B:     0xFF,
					C:     0xFE,
					H:     0x00,
					L:     0x01,
					PC:    0x01,
					Flags: cf,
				},
				ram("09"),
			),
		},
		{
			"DAD D",
			newComputer(
				CPU{
					D: 0x01,
					E: 0xFF,
					H: 0x01,
					L: 0x01,
				},
				ram("19"),
			),
			newComputer(
				CPU{
					D:     0x01,
					E:     0xFF,
					H:     0x03,
					L:     0x00,
					PC:    0x01,
					Flags: none,
				},
				ram("19"),
			),
		},
		{
			"DAD H",
			newComputer(
				CPU{
					H: 0x01,
					L: 0x01,
				},
				ram("29"),
			),
			newComputer(
				CPU{
					H:     0x02,
					L:     0x02,
					PC:    0x01,
					Flags: none,
				},
				ram("29"),
			),
		},
		{
			"DAD SP",
			newComputer(
				CPU{
					H:  0x01,
					L:  0x03,
					SP: 0x0FFF,
				},
				ram("39"),
			),
			newComputer(
				CPU{
					H:     0x11,
					L:     0x02,
					SP:    0x0FFF,
					PC:    0x01,
					Flags: none,
				},
				ram("39"),
			),
		},
		{
			"DCR A",
			newComputer(
				CPU{
					A: 0x02,
				},
				ram("3D"),
			),
			newComputer(
				CPU{
					A:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				ram("3D"),
			),
		},
		{
			"DCR B",
			newComputer(
				CPU{
					B: 0x02,
				},
				ram("05"),
			),
			newComputer(
				CPU{
					B:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				ram("05"),
			),
		},
		{
			"DCR B: carry not set when there's borrow",
			newComputer(
				CPU{
					B: 0x00,
				},
				ram("05"),
			),
			newComputer(
				CPU{
					B:     0xff,
					PC:    0x01,
					Flags: sf | pf,
				},
				ram("05"),
			),
		},
		{
			"DCR B: carry not modified when there was existing carry",
			newComputer(
				CPU{
					B:     0x00,
					Flags: cf,
				},
				ram("05"),
			),
			newComputer(
				CPU{
					B:     0xff,
					PC:    0x01,
					Flags: sf | pf | cf,
				},
				ram("05"),
			),
		},
		{
			"DCR B: Generates auxiliary carry when there's carry in the lower nibble",
			newComputer(
				CPU{
					B: 0x1D,
				},
				ram("05"),
			),
			newComputer(
				CPU{
					B:     0x1C,
					PC:    0x01,
					Flags: acf,
				},
				ram("05"),
			),
		},
		{
			"DCR C",
			newComputer(
				CPU{
					C: 0x02,
				},
				ram("0D"),
			),
			newComputer(
				CPU{
					C:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				ram("0D"),
			),
		},
		{
			"DCR D",
			newComputer(
				CPU{
					D: 0x02,
				},
				ram("15"),
			),
			newComputer(
				CPU{
					D:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				ram("15"),
			),
		},
		{
			"DCR E",
			newComputer(
				CPU{
					E: 0x02,
				},
				ram("1D"),
			),
			newComputer(
				CPU{
					E:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				ram("1D"),
			),
		},
		{
			"DCR H",
			newComputer(
				CPU{
					H: 0x02,
				},
				ram("20"),
			),
			newComputer(
				CPU{
					H:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				ram("20"),
			),
		},
		{
			"DCR L",
			newComputer(
				CPU{
					L: 0x02,
				},
				ram("2D"),
			),
			newComputer(
				CPU{
					L:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				ram("2D"),
			),
		},
		{
			"INR A",
			newComputer(
				CPU{
					A: 0xFF,
				},
				ram("3C"),
			),
			newComputer(
				CPU{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("3C"),
			),
		},
		{
			"INR B",
			newComputer(
				CPU{
					B: 0xFF,
				},
				ram("04"),
			),
			newComputer(
				CPU{
					B:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("04"),
			),
		},
		{
			"INR B: generates auxiliary carry",
			newComputer(
				CPU{
					B: 0x0F,
				},
				ram("04"),
			),
			newComputer(
				CPU{
					B:     0x10,
					PC:    0x01,
					Flags: acf,
				},
				ram("04"),
			),
		},
		{
			"INR C",
			newComputer(
				CPU{
					C: 0x0f,
				},
				ram("0C"),
			),
			newComputer(
				CPU{
					C:     0x10,
					PC:    0x01,
					Flags: acf,
				},
				ram("0C"),
			),
		},
		{
			"INR D",
			newComputer(
				CPU{
					D: 0x03,
				},
				ram("14"),
			),
			newComputer(
				CPU{
					D:     0x04,
					PC:    0x01,
					Flags: none,
				},
				ram("14"),
			),
		},
		{
			"INR E",
			newComputer(
				CPU{
					E: 0x03,
				},
				ram("1C"),
			),
			newComputer(
				CPU{
					E:     0x04,
					PC:    0x01,
					Flags: none,
				},
				ram("1C"),
			),
		},
		{
			"INR H",
			newComputer(
				CPU{
					H: 0x03,
				},
				ram("24"),
			),
			newComputer(
				CPU{
					H:     0x04,
					PC:    0x01,
					Flags: none,
				},
				ram("24"),
			),
		},
		{
			"INR L",
			newComputer(
				CPU{
					L: 0x03,
				},
				ram("2C"),
			),
			newComputer(
				CPU{
					L:     0x04,
					PC:    0x01,
					Flags: none,
				},
				ram("2C"),
			),
		},
		{
			"INX B",
			newComputer(
				CPU{
					C: 0xFF,
				},
				ram("03"),
			),
			newComputer(
				CPU{
					B:  0x01,
					C:  0x00,
					PC: 0x01,
				},
				ram("03"),
			),
		},
		{
			"INX D",
			newComputer(
				CPU{
					E: 0xFF,
				},
				ram("13"),
			),
			newComputer(
				CPU{
					D:  0x01,
					E:  0x00,
					PC: 0x01,
				},
				ram("13"),
			),
		},
		{
			"INX H",
			newComputer(
				CPU{
					L: 0xFF,
				},
				ram("23"),
			),
			newComputer(
				CPU{
					H:  0x01,
					L:  0x00,
					PC: 0x01,
				},
				ram("23"),
			),
		},
		{
			"INX SP",
			newComputer(
				CPU{
					SP: 0x0F,
				},
				ram("33"),
			),
			newComputer(
				CPU{
					PC: 0x01,
					SP: 0x10,
				},
				ram("33"),
			),
		},
		{
			"JNZ adr: zero flag set",
			newComputer(
				CPU{
					Flags: zf,
				},
				ram("C2 0A 00"),
			),
			newComputer(
				CPU{
					PC:    0x03,
					Flags: zf,
				},
				ram("C2 0A 00"),
			),
		},
		{
			"JNZ adr: zero flag not set",
			newComputer(
				CPU{},
				ram("C2 0A 00"),
			),
			newComputer(
				CPU{
					PC: 0x0A,
				},
				ram("C2 0A 00"),
			),
		},
		{
			"JMP adr",
			newComputer(
				CPU{},
				ram("C3 0A 00"),
			),
			newComputer(
				CPU{
					PC: 0x0A,
				},
				ram("C3 0A 00"),
			),
		},
		{
			"LDAX B",
			newComputer(
				CPU{
					B: 0x00,
					C: 0x02,
				},
				ram("0a 00 ff"),
			),
			newComputer(
				CPU{
					A:  0xFF,
					B:  0x00,
					C:  0x02,
					PC: 0x01,
				},
				ram("0a 00 ff"),
			),
		},
		{
			"LDAX D",
			newComputer(
				CPU{
					D: 0x00,
					E: 0x02,
				},
				ram("1a 00 ff"),
			),
			newComputer(
				CPU{
					A:  0xFF,
					D:  0x00,
					E:  0x02,
					PC: 0x01,
				},
				ram("1a 00 ff"),
			),
		},
		{
			"LXI B, D16",
			newComputer(
				CPU{},
				ram("01 0B 01"),
			),
			newComputer(
				CPU{
					B:  0x01,
					C:  0x0B,
					PC: 0x03,
				},
				ram("01 0B 01"),
			),
		},
		{
			"LXI D, D16",
			newComputer(
				CPU{},
				ram("11 0B 01"),
			),
			newComputer(
				CPU{
					D:  0x01,
					E:  0x0B,
					PC: 0x03,
				},
				ram("11 0B 01"),
			),
		},
		{
			"LXI H, D16",
			newComputer(
				CPU{},
				ram("21 0B 01"),
			),
			newComputer(
				CPU{
					H:  0x01,
					L:  0x0B,
					PC: 0x03,
				},
				ram("21 0B 01"),
			),
		},
		{
			"LXI SP, D16",
			newComputer(
				CPU{},
				ram("31 0B 01"),
			),
			newComputer(
				CPU{
					PC: 0x03,
					SP: 0x010B,
				},
				ram("31 0B 01"),
			),
		},
		{
			"MOV A, A",
			newComputer(
				CPU{
					A: 0x01,
				},
				ram("7F"),
			),
			newComputer(
				CPU{
					A:  0x01,
					PC: 0x01,
				},
				ram("7F"),
			),
		},
		{
			"MOV A, B",
			newComputer(
				CPU{
					B: 0x01,
				},
				ram("78"),
			),
			newComputer(
				CPU{
					A:  0x01,
					B:  0x01,
					PC: 0x01,
				},
				ram("78"),
			),
		},
		{
			"MOV A, C",
			newComputer(
				CPU{
					C: 0x01,
				},
				ram("79"),
			),
			newComputer(
				CPU{
					A:  0x01,
					C:  0x01,
					PC: 0x01,
				},
				ram("79"),
			),
		},
		{
			"MOV A, D",
			newComputer(
				CPU{
					D: 0x01,
				},
				ram("7A"),
			),
			newComputer(
				CPU{
					A:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				ram("7A"),
			),
		},
		{
			"MOV A, E",
			newComputer(
				CPU{
					E: 0x01,
				},
				ram("7B"),
			),
			newComputer(
				CPU{
					A:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				ram("7B"),
			),
		},
		{
			"MOV A, H",
			newComputer(
				CPU{
					H: 0x01,
				},
				ram("7C"),
			),
			newComputer(
				CPU{
					A:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("7C"),
			),
		},
		{
			"MOV A, L",
			newComputer(
				CPU{
					L: 0x01,
				},
				ram("7D"),
			),
			newComputer(
				CPU{
					A:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("7D"),
			),
		},
		{
			"MOV A, M",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x02,
				},
				ram("7E 00 FF"),
			),
			newComputer(
				CPU{
					A:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				ram("7E 00 FF"),
			),
		},
		{
			"MOV B, A",
			newComputer(
				CPU{
					A: 0x01,
				},
				ram("47"),
			),
			newComputer(
				CPU{
					A:  0x01,
					B:  0x01,
					PC: 0x01,
				},
				ram("47"),
			),
		},
		{
			"MOV B, B",
			newComputer(
				CPU{
					B: 0x01,
				},
				ram("40"),
			),
			newComputer(
				CPU{
					B:  0x01,
					PC: 0x01,
				},
				ram("40"),
			),
		},
		{
			"MOV B, C",
			newComputer(
				CPU{
					C: 0x01,
				},
				ram("41"),
			),
			newComputer(
				CPU{
					B:  0x01,
					C:  0x01,
					PC: 0x01,
				},
				ram("41"),
			),
		},
		{
			"MOV B, D",
			newComputer(
				CPU{
					D: 0x01,
				},
				ram("42"),
			),
			newComputer(
				CPU{
					B:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				ram("42"),
			),
		},
		{
			"MOV B, E",
			newComputer(
				CPU{
					E: 0x01,
				},
				ram("43"),
			),
			newComputer(
				CPU{
					B:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				ram("43"),
			),
		},
		{
			"MOV B, H",
			newComputer(
				CPU{
					H: 0x01,
				},
				ram("44"),
			),
			newComputer(
				CPU{
					B:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("44"),
			),
		},
		{
			"MOV B, L",
			newComputer(
				CPU{
					L: 0x01,
				},
				ram("45"),
			),
			newComputer(
				CPU{
					B:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("45"),
			),
		},
		{
			"MOV B, M",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x02,
				},
				ram("46 00 FF"),
			),
			newComputer(
				CPU{
					B:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				ram("46 00 FF"),
			),
		},
		{
			"MOV C, A",
			newComputer(
				CPU{
					A: 0x01,
				},
				ram("4F"),
			),
			newComputer(
				CPU{
					A:  0x01,
					C:  0x01,
					PC: 0x01,
				},
				ram("4F"),
			),
		},
		{
			"MOV C, B",
			newComputer(
				CPU{
					B: 0x01,
				},
				ram("48"),
			),
			newComputer(
				CPU{
					B:  0x01,
					C:  0x01,
					PC: 0x01,
				},
				ram("48"),
			),
		},
		{
			"MOV C, C",
			newComputer(
				CPU{
					C: 0x01,
				},
				ram("49"),
			),
			newComputer(
				CPU{
					C:  0x01,
					PC: 0x01,
				},
				ram("49"),
			),
		},
		{
			"MOV C, D",
			newComputer(
				CPU{
					D: 0x01,
				},
				ram("4A"),
			),
			newComputer(
				CPU{
					C:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				ram("4A"),
			),
		},
		{
			"MOV C, E",
			newComputer(
				CPU{
					E: 0x01,
				},
				ram("4B"),
			),
			newComputer(
				CPU{
					C:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				ram("4B"),
			),
		},
		{
			"MOV C, H",
			newComputer(
				CPU{
					H: 0x01,
				},
				ram("4C"),
			),
			newComputer(
				CPU{
					C:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("4C"),
			),
		},
		{
			"MOV C, L",
			newComputer(
				CPU{
					L: 0x01,
				},
				ram("4D"),
			),
			newComputer(
				CPU{
					C:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("4D"),
			),
		},
		{
			"MOV C, M",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x02,
				},
				ram("4E 00 FF"),
			),
			newComputer(
				CPU{
					C:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				ram("4E 00 FF"),
			),
		},
		{
			"MOV D, A",
			newComputer(
				CPU{
					A: 0x01,
				},
				ram("57"),
			),
			newComputer(
				CPU{
					A:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				ram("57"),
			),
		},
		{
			"MOV D, B",
			newComputer(
				CPU{
					B: 0x01,
				},
				ram("50"),
			),
			newComputer(
				CPU{
					B:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				ram("50"),
			),
		},
		{
			"MOV D, C",
			newComputer(
				CPU{
					C: 0x01,
				},
				ram("51"),
			),
			newComputer(
				CPU{
					C:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				ram("51"),
			),
		},
		{
			"MOV D, D",
			newComputer(
				CPU{
					D: 0x01,
				},
				ram("52"),
			),
			newComputer(
				CPU{
					D:  0x01,
					PC: 0x01,
				},
				ram("52"),
			),
		},
		{
			"MOV D, E",
			newComputer(
				CPU{
					E: 0x01,
				},
				ram("53"),
			),
			newComputer(
				CPU{
					D:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				ram("53"),
			),
		},
		{
			"MOV D, H",
			newComputer(
				CPU{
					H: 0x01,
				},
				ram("54"),
			),
			newComputer(
				CPU{
					D:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("54"),
			),
		},
		{
			"MOV D, L",
			newComputer(
				CPU{
					L: 0x01,
				},
				ram("55"),
			),
			newComputer(
				CPU{
					D:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("55"),
			),
		},
		{
			"MOV D, M",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x02,
				},
				ram("56 00 FF"),
			),
			newComputer(
				CPU{
					D:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				ram("56 00 FF"),
			),
		},
		{
			"MOV E, A",
			newComputer(
				CPU{
					A: 0x01,
				},
				ram("5F"),
			),
			newComputer(
				CPU{
					A:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				ram("5F"),
			),
		},
		{
			"MOV E, B",
			newComputer(
				CPU{
					B: 0x01,
				},
				ram("58"),
			),
			newComputer(
				CPU{
					B:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				ram("58"),
			),
		},
		{
			"MOV E, C",
			newComputer(
				CPU{
					C: 0x01,
				},
				ram("59"),
			),
			newComputer(
				CPU{
					C:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				ram("59"),
			),
		},
		{
			"MOV E, D",
			newComputer(
				CPU{
					D: 0x01,
				},
				ram("5A"),
			),
			newComputer(
				CPU{
					D:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				ram("5A"),
			),
		},
		{
			"MOV E, E",
			newComputer(
				CPU{
					E: 0x01,
				},
				ram("5B"),
			),
			newComputer(
				CPU{
					PC: 0x01,
					E:  0x01,
				},
				ram("5B"),
			),
		},
		{
			"MOV E, H",
			newComputer(
				CPU{
					H: 0x01,
				},
				ram("5C"),
			),
			newComputer(
				CPU{
					E:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("5C"),
			),
		},
		{
			"MOV E, L",
			newComputer(
				CPU{
					L: 0x01,
				},
				ram("5D"),
			),
			newComputer(
				CPU{
					E:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("5D"),
			),
		},
		{
			"MOV E, M",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x02,
				},
				ram("5E 00 FF"),
			),
			newComputer(
				CPU{
					E:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				ram("5E 00 FF"),
			),
		},
		{
			"MOV H, A",
			newComputer(
				CPU{
					A: 0x01,
				},
				ram("67"),
			),
			newComputer(
				CPU{
					A:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("67"),
			),
		},
		{
			"MOV H, B",
			newComputer(
				CPU{
					B: 0x01,
				},
				ram("60"),
			),
			newComputer(
				CPU{
					B:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("60"),
			),
		},
		{
			"MOV H, C",
			newComputer(
				CPU{
					C: 0x01,
				},
				ram("61"),
			),
			newComputer(
				CPU{
					C:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("61"),
			),
		},
		{
			"MOV H, D",
			newComputer(
				CPU{
					D: 0x01,
				},
				ram("62"),
			),
			newComputer(
				CPU{
					D:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("62"),
			),
		},
		{
			"MOV H, E",
			newComputer(
				CPU{
					E: 0x01,
				},
				ram("63"),
			),
			newComputer(
				CPU{
					E:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				ram("63"),
			),
		},
		{
			"MOV H, H",
			newComputer(
				CPU{
					H: 0x01,
				},
				ram("64"),
			),
			newComputer(
				CPU{
					H:  0x01,
					PC: 0x01,
				},
				ram("64"),
			),
		},
		{
			"MOV H, L",
			newComputer(
				CPU{
					L: 0x01,
				},
				ram("65"),
			),
			newComputer(
				CPU{
					H:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("65"),
			),
		},
		{
			"MOV H, M",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x02,
				},
				ram("66 00 FF"),
			),
			newComputer(
				CPU{
					H:  0xFF,
					L:  0x02,
					PC: 0x01,
				},
				ram("66 00 FF"),
			),
		},
		{
			"MOV L, A",
			newComputer(
				CPU{
					A: 0x01,
				},
				ram("6F"),
			),
			newComputer(
				CPU{
					A:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("6F"),
			),
		},
		{
			"MOV L, B",
			newComputer(
				CPU{
					B: 0x01,
				},
				ram("68"),
			),
			newComputer(
				CPU{
					B:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("68"),
			),
		},
		{
			"MOV L, C",
			newComputer(
				CPU{
					C: 0x01,
				},
				ram("69"),
			),
			newComputer(
				CPU{
					C:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("69"),
			),
		},
		{
			"MOV L, D",
			newComputer(
				CPU{
					D: 0x01,
				},
				ram("6A"),
			),
			newComputer(
				CPU{
					D:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("6A"),
			),
		},
		{
			"MOV L, E",
			newComputer(
				CPU{
					E: 0x01,
				},
				ram("6B"),
			),
			newComputer(
				CPU{
					E:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("6B"),
			),
		},
		{
			"MOV L, H",
			newComputer(
				CPU{
					H: 0x01,
				},
				ram("6C"),
			),
			newComputer(
				CPU{
					H:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				ram("6C"),
			),
		},
		{
			"MOV L, L",
			newComputer(
				CPU{
					L: 0x01,
				},
				ram("6D"),
			),
			newComputer(
				CPU{
					L:  0x01,
					PC: 0x01,
				},
				ram("6D"),
			),
		},
		{
			"MOV L, M",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x02,
				},
				ram("6E 00 FF"),
			),
			newComputer(
				CPU{
					H:  0x00,
					L:  0xFF,
					PC: 0x01,
				},
				ram("6E 00 FF"),
			),
		},
		{
			"MOV M, A",
			newComputer(
				CPU{
					A: 0xff,
					H: 0x00,
					L: 0x03,
				},
				ram("77 00 00 00"),
			),
			newComputer(
				CPU{
					A:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				ram("77 00 00 FF"),
			),
		},
		{
			"MOV M, B",
			newComputer(
				CPU{
					B: 0xff,
					H: 0x00,
					L: 0x03,
				},
				ram("70 00 00 00"),
			),
			newComputer(
				CPU{
					B:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				ram("70 00 00 FF"),
			),
		},
		{
			"MOV M, C",
			newComputer(
				CPU{
					C: 0xff,
					H: 0x00,
					L: 0x03,
				},
				ram("71 00 00 00"),
			),
			newComputer(
				CPU{
					C:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				ram("71 00 00 FF"),
			),
		},
		{
			"MOV M, D",
			newComputer(
				CPU{
					D: 0xff,
					H: 0x00,
					L: 0x03,
				},
				ram("72 00 00 00"),
			),
			newComputer(
				CPU{
					D:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				ram("72 00 00 FF"),
			),
		},
		{
			"MOV M, E",
			newComputer(
				CPU{
					E: 0xff,
					H: 0x00,
					L: 0x03,
				},
				ram("73 00 00 00"),
			),
			newComputer(
				CPU{
					E:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				ram("73 00 00 FF"),
			),
		},
		{
			"MOV M, H",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x03,
				},
				ram("74 00 00 00"),
			),
			newComputer(
				CPU{
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				ram("74 00 00 00"),
			),
		},
		{
			"MOV M, L",
			newComputer(
				CPU{
					H: 0x00,
					L: 0x03,
				},
				ram("75 00 00 00"),
			),
			newComputer(
				CPU{
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				ram("75 00 00 03"),
			),
		},
		{
			"MVI A, D8",
			newComputer(
				CPU{},
				ram("3E 0B"),
			),
			newComputer(
				CPU{
					A:  0x0B,
					PC: 0x02,
				},
				ram("3E 0B"),
			),
		},
		{
			"MVI B, D8",
			newComputer(
				CPU{},
				ram("06 0B"),
			),
			newComputer(
				CPU{
					B:  0x0B,
					PC: 0x02,
				},
				ram("06 0B"),
			),
		},
		{
			"MVI C, D8",
			newComputer(
				CPU{},
				ram("0E 0B"),
			),
			newComputer(
				CPU{
					C:  0x0B,
					PC: 0x02,
				},
				ram("0E 0B"),
			),
		},
		{
			"MVI D, D8",
			newComputer(
				CPU{},
				ram("16 0B"),
			),
			newComputer(
				CPU{
					D:  0x0B,
					PC: 0x02,
				},
				ram("16 0B"),
			),
		},
		{
			"MVI E, D8",
			newComputer(
				CPU{},
				ram("1E 0B"),
			),
			newComputer(
				CPU{
					E:  0x0B,
					PC: 0x02,
				},
				ram("1E 0B"),
			),
		},
		{
			"MVI H, D8",
			newComputer(
				CPU{},
				ram("26 0B"),
			),
			newComputer(
				CPU{
					H:  0x0B,
					PC: 0x02,
				},
				ram("26 0B"),
			),
		},
		{
			"MVI L, D8",
			newComputer(
				CPU{},
				ram("2E 0B"),
			),
			newComputer(
				CPU{
					L:  0x0B,
					PC: 0x02,
				},
				ram("2E 0B"),
			),
		},
		{
			"NOP",
			newComputer(
				CPU{},
				ram("00"),
			),
			newComputer(
				CPU{
					PC: 0x01,
				},
				ram("00"),
			),
		},
		{
			"ORA A",
			newComputer(
				CPU{
					A:     0xFF,
					Flags: cf,
				},
				ram("B7"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("B7"),
			),
		},
		{
			"ORA B",
			newComputer(
				CPU{
					A:     0xFF,
					B:     0x0A,
					Flags: cf,
				},
				ram("B0"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					B:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("B0"),
			),
		},
		{
			"ORA C",
			newComputer(
				CPU{
					A: 0xFF,
					C: 0x0A,
				},
				ram("B1"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					C:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("B1"),
			),
		},
		{
			"ORA D",
			newComputer(
				CPU{
					A: 0xFF,
					D: 0x0A,
				},
				ram("B2"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					D:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("B2"),
			),
		},
		{
			"ORA E",
			newComputer(
				CPU{
					A: 0xFF,
					E: 0x0A,
				},
				ram("B3"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					E:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("B3"),
			),
		},
		{
			"ORA H",
			newComputer(
				CPU{
					A: 0xFF,
					H: 0x0A,
				},
				ram("B4"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					H:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("B4"),
			),
		},
		{
			"ORA L",
			newComputer(
				CPU{
					A: 0xFF,
					L: 0x0A,
				},
				ram("B5"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					L:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("B5"),
			),
		},
		{
			"RET",
			newComputer(
				CPU{
					SP: 0x01,
				},
				ram("C9 04 00 00 00"),
			),
			newComputer(
				CPU{
					SP: 0x03,
					PC: 0x04,
				},
				ram("C9 04 00 00 00"),
			),
		},
		{
			"SBB A: with borrow",
			newComputer(
				CPU{
					A:     0x01,
					Flags: cf,
				},
				ram("9F"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					PC:    0x01,
					Flags: sf | pf | cf,
				},
				ram("9F"),
			),
		},
		{
			"SBB A: no borrow",
			newComputer(
				CPU{
					A:     0x01,
					Flags: none,
				},
				ram("9F"),
			),
			newComputer(
				CPU{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("9F"),
			),
		},
		{
			"SBB B",
			newComputer(
				CPU{
					A:     0x02,
					B:     0x01,
					Flags: cf,
				},
				ram("98"),
			),
			newComputer(
				CPU{
					A:     0x00,
					B:     0x01,
					PC:    1,
					Flags: zf | pf,
				},
				ram("98"),
			),
		},
		{
			"SBB C",
			newComputer(
				CPU{
					A:     0x00,
					C:     0xFF,
					Flags: cf,
				},
				ram("99"),
			),
			newComputer(
				CPU{
					A:     0x00,
					C:     0xFF,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("99"),
			),
		},
		{
			"SBB D",
			newComputer(
				CPU{
					A:     0x00,
					D:     0xFF,
					Flags: cf,
				},
				ram("9A"),
			),
			newComputer(
				CPU{
					A:     0x00,
					D:     0xFF,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("9A"),
			),
		},
		{
			"SBB E",
			newComputer(
				CPU{
					A:     0x00,
					E:     0xFF,
					Flags: cf,
				},
				ram("9B"),
			),
			newComputer(
				CPU{
					A:     0x00,
					E:     0xFF,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("9B"),
			),
		},
		{
			"SBB H",
			newComputer(
				CPU{
					A:     0x00,
					H:     0xFF,
					Flags: cf,
				},
				ram("9C"),
			),
			newComputer(
				CPU{
					A:     0x00,
					H:     0xFF,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("9C"),
			),
		},
		{
			"SBB L",
			newComputer(
				CPU{
					A:     0x08,
					L:     0x02,
					Flags: cf,
				},
				ram("9D"),
			),
			newComputer(
				CPU{
					A:     0x05,
					L:     0x02,
					PC:    0x01,
					Flags: pf,
				},
				ram("9D"),
			),
		},
		{
			"STAX B",
			newComputer(
				CPU{
					A: 0xFF,
					B: 0x00,
					C: 0x10,
				},
				ram("02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"),
			),
			newComputer(
				CPU{
					A:  0xFF,
					B:  0x00,
					C:  0x10,
					PC: 0x01,
				},
				ram("02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF"),
			),
		},
		{
			"STAX D",
			newComputer(
				CPU{
					A: 0xFF,
					D: 0x00,
					E: 0x10,
				},
				ram("12 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"),
			),
			newComputer(
				CPU{
					A:  0xFF,
					D:  0x00,
					E:  0x10,
					PC: 0x01,
				},
				ram("12 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF"),
			),
		},
		{
			"SUB A",
			newComputer(
				CPU{
					A: 0x30,
				},
				ram("97"),
			),
			newComputer(
				CPU{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("97"),
			),
		},
		{
			"SUB B: substracting 0",
			newComputer(
				CPU{
					A: 0x00,
					B: 0x00,
				},
				ram("90"),
			),
			newComputer(
				CPU{
					A:     0x00,
					B:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("90"),
			),
		},
		{
			"SUB B: 0x02 - 0x01",
			newComputer(
				CPU{
					A: 0x02,
					B: 0x01,
				},
				ram("90"),
			),
			newComputer(
				CPU{
					A:     0x01,
					B:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("90"),
			),
		},
		{
			"SUB B: 0x01 - 0x02",
			newComputer(
				CPU{
					A: 0x01,
					B: 0x02,
				},
				ram("90"),
			),
			newComputer(
				CPU{
					A:     0xFF,
					B:     0x02,
					PC:    0x01,
					Flags: sf | cf | pf,
				},
				ram("90"),
			),
		},
		{
			"SUB C",
			newComputer(
				CPU{
					A: 0x30,
					C: 0x01,
				},
				ram("91"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					C:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("91"),
			),
		},
		{
			"SUB D",
			newComputer(
				CPU{
					A: 0x30,
					D: 0x01,
				},
				ram("92"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					D:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("92"),
			),
		},
		{
			"SUB E",
			newComputer(
				CPU{
					A: 0x30,
					E: 0x01,
				},
				ram("93"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					E:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("93"),
			),
		},
		{
			"SUB H",
			newComputer(
				CPU{
					A: 0x30,
					H: 0x01,
				},
				ram("94"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					H:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("94"),
			),
		},
		{
			"SUB L",
			newComputer(
				CPU{
					A: 0x30,
					L: 0x01,
				},
				ram("95"),
			),
			newComputer(
				CPU{
					A:     0x2F,
					L:     0x01,
					PC:    0x01,
					Flags: none,
				},
				ram("95"),
			),
		},
		{
			"XRA A",
			newComputer(
				CPU{
					A: 0xFF,
				},
				ram("A8"),
			),
			newComputer(
				CPU{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				ram("A8"),
			),
		},
		{
			"XRA B",
			newComputer(
				CPU{
					A: 0xFF,
					B: 0x0A,
				},
				ram("A9"),
			),
			newComputer(
				CPU{
					A:     0xF5,
					B:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("A9"),
			),
		},
		{
			"XRA C",
			newComputer(
				CPU{
					A: 0xFF,
					C: 0x0A,
				},
				ram("AA"),
			),
			newComputer(
				CPU{
					A:     0xF5,
					C:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("AA"),
			),
		},
		{
			"XRA D",
			newComputer(
				CPU{
					A: 0xFF,
					D: 0x0A,
				},
				ram("AB"),
			),
			newComputer(
				CPU{
					A:     0xF5,
					D:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("AB"),
			),
		},
		{
			"XRA E",
			newComputer(
				CPU{
					A: 0xFF,
					E: 0x0A,
				},
				ram("AC"),
			),
			newComputer(
				CPU{
					A:     0xF5,
					E:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("AC"),
			),
		},
		{
			"XRA H",
			newComputer(
				CPU{
					A: 0xFF,
					H: 0x0A,
				},
				ram("AD"),
			),
			newComputer(
				CPU{
					A:     0xF5,
					H:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("AD"),
			),
		},
		{
			"XRA L",
			newComputer(
				CPU{
					A: 0xFF,
					L: 0x0A,
				},
				ram("AF"),
			),
			newComputer(
				CPU{
					A:     0xF5,
					L:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				ram("AF"),
			),
		},
	} {
		t.Run(tC.desc, func(t *testing.T) {
			err := tC.init.Step()
			if err != nil {
				t.Fatalf("Unexpected error %v", err)
			}
			if !(tC.init.CPU == tC.want.CPU && bytes.Equal(tC.init.Mem, tC.want.Mem)) {
				t.Fatalf("got: \n%v\n want: \n%v", tC.init, tC.want)
			}
		})
	}
}
