package emu

import (
	"reflect"
	"testing"

	"github.com/miguelff/8080/encoding"
)

func ram(bytes string) memory {
	return encoding.HexToBin(bytes)
}
func TestParity(t *testing.T) {
	for _, tC := range []struct {
		b    byte
		want flags
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
		want flags
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
		want flags
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
		desc    string
		init    *Computer
		want    *Computer
		wantErr error
	}{
		{
			"ADC A: with carry",
			&Computer{
				cpu: cpu{
					A:     0x02,
					Flags: cf,
				},
				mem: ram("8F"),
			},
			&Computer{
				cpu: cpu{
					A:     0x05,
					PC:    0x01,
					Flags: pf,
				},
				mem: ram("8F"),
			},
			nil,
		},
		{
			"ADC A: no carry",
			&Computer{
				cpu: cpu{
					A:     0x02,
					Flags: none,
				},
				mem: ram("8F"),
			},
			&Computer{
				cpu: cpu{
					A:     0x04,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("8F"),
			},
			nil,
		},
		{
			"ADC B",
			&Computer{
				cpu: cpu{
					A:     0x02,
					B:     0x01,
					Flags: cf,
				},
				mem: ram("88"),
			},
			&Computer{
				cpu: cpu{
					A:     0x04,
					B:     0x01,
					PC:    1,
					Flags: none,
				},
				mem: ram("88"),
			},
			nil,
		},
		{
			"ADC C",
			&Computer{
				cpu: cpu{
					A:     0x01,
					C:     0xFD,
					Flags: cf,
				},
				mem: ram("89"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					C:     0xFD,
					PC:    0x01,
					Flags: sf | pf,
				},
				mem: ram("89"),
			},
			nil,
		},
		{
			"ADC D",
			&Computer{
				cpu: cpu{
					A:     0x00,
					D:     0xFD,
					Flags: cf,
				},
				mem: ram("8A"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFE,
					D:     0xFD,
					PC:    0x01,
					Flags: sf,
				},
				mem: ram("8A"),
			},
			nil,
		},
		{
			"ADC E",
			&Computer{
				cpu: cpu{
					A:     0x00,
					E:     0xFD,
					Flags: cf,
				},
				mem: ram("8B"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFE,
					E:     0xFD,
					PC:    0x01,
					Flags: sf,
				},
				mem: ram("8B"),
			},
			nil,
		},
		{
			"ADC H",
			&Computer{
				cpu: cpu{
					A:     0x00,
					H:     0xFD,
					Flags: cf,
				},
				mem: ram("8C"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFE,
					H:     0xFD,
					PC:    0x01,
					Flags: sf,
				},
				mem: ram("8C"),
			},
			nil,
		},
		{
			"ADC L",
			&Computer{
				cpu: cpu{
					A:     0x00,
					L:     0xFD,
					Flags: cf,
				},
				mem: ram("8D"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFE,
					L:     0xFD,
					PC:    0x01,
					Flags: sf,
				},
				mem: ram("8D"),
			},
			nil,
		},
		{
			"ADC M",
			&Computer{
				cpu: cpu{
					A:     0x01,
					H:     0x00,
					L:     0x02,
					Flags: cf,
				},
				mem: ram("8E 00 FD"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					H:     0x00,
					L:     0x02,
					PC:    0x01,
					Flags: sf | pf,
				},
				mem: ram("8E 00 FD"),
			},
			nil,
		},
		{
			"ADD A",
			&Computer{
				cpu: cpu{
					A: 0x02,
				},
				mem: ram("87"),
			},
			&Computer{
				cpu: cpu{
					A:     0x04,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("87"),
			},
			nil,
		},
		{
			"ADD B: adding two values that sum 0x0 sets the zero8 and parity8 flags",
			&Computer{
				mem: ram("80"),
			},
			&Computer{
				cpu: cpu{
					PC:    1,
					Flags: zf | pf,
				},
				mem: ram("80"),
			},
			nil,
		},
		{
			"ADD B: adding 0x09 + 0x07 sets the auxiliary carry flag",
			&Computer{
				cpu: cpu{
					A: 0x09,
					B: 0x07,
				},
				mem: ram("80"),
			},
			&Computer{
				cpu: cpu{
					A:     0x10,
					B:     0x07,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("80"),
			},
			nil,
		},
		{
			"ADD B: adding 0xFE and 0x02 sets the carry and auxiliary carry flags",
			&Computer{
				cpu: cpu{
					A: 0x03,
					B: 0xFE,
				},
				mem: ram("80"),
			},
			&Computer{
				cpu: cpu{
					A:     0x01,
					B:     0xFE,
					PC:    0x01,
					Flags: cf | acf,
				},
				mem: ram("80"),
			},
			nil,
		},
		{
			"ADD C: adding 0xFE and 0x01 set the parity8 and sign8 flags",
			&Computer{
				cpu: cpu{
					A: 0x01,
					C: 0xFE,
				},
				mem: ram("81"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					C:     0xFE,
					PC:    0x01,
					Flags: sf | pf,
				},
				mem: ram("81"),
			},
			nil,
		},
		{
			"ADD D: adding 0xFf and 0x00 set the sign8 flag",
			&Computer{
				cpu: cpu{
					A: 0x00,
					D: 0xFE,
				},
				mem: ram("82"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFE,
					D:     0xFE,
					PC:    0x01,
					Flags: sf,
				},
				mem: ram("82"),
			},
			nil,
		},
		{
			"ADD E",
			&Computer{
				cpu: cpu{
					A: 0x00,
					E: 0xFE,
				},
				mem: ram("83"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFE,
					E:     0xFE,
					PC:    0x01,
					Flags: sf,
				},
				mem: ram("83"),
			},
			nil,
		},
		{
			"ADD H",
			&Computer{
				cpu: cpu{
					A: 0x00,
					H: 0xFE,
				},
				mem: ram("84"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFE,
					H:     0xFE,
					PC:    0x01,
					Flags: sf,
				},
				mem: ram("84"),
			},
			nil,
		},
		{
			"ADD L",
			&Computer{
				cpu: cpu{
					A: 0x00,
					L: 0xFE,
				},
				mem: ram("85"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFE,
					L:     0xFE,
					PC:    0x01,
					Flags: sf,
				},
				mem: ram("85"),
			},
			nil,
		},
		{
			"ADD M",
			&Computer{
				cpu: cpu{
					A: 0x01,
					H: 0x00,
					L: 0x02,
				},
				mem: ram("86 00 FE"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					H:     0x00,
					L:     0x02,
					PC:    0x01,
					Flags: sf | pf,
				},
				mem: ram("86 00 FE"),
			},
			nil,
		},
		{
			"ANA A",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					Flags: cf,
				},
				mem: ram("A7"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("A7"),
			},
			nil,
		},
		{
			"ANA B",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					B:     0x0A,
					Flags: cf,
				},
				mem: ram("A0"),
			},
			&Computer{
				cpu: cpu{
					A:     0x0A,
					B:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				mem: ram("A0"),
			},
			nil,
		},
		{
			"ANA C",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					C:     0x0A,
					Flags: cf,
				},
				mem: ram("A1"),
			},
			&Computer{
				cpu: cpu{
					A:     0x0A,
					C:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				mem: ram("A1"),
			},
			nil,
		},
		{
			"ANA D",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					D:     0x0A,
					Flags: cf,
				},
				mem: ram("A2"),
			},
			&Computer{
				cpu: cpu{
					A:     0x0A,
					D:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				mem: ram("A2"),
			},
			nil,
		},
		{
			"ANA E",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					E:     0x0A,
					Flags: cf,
				},
				mem: ram("A3"),
			},
			&Computer{
				cpu: cpu{
					A:     0x0A,
					E:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				mem: ram("A3"),
			},
			nil,
		},
		{
			"ANA H",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					H:     0x0A,
					Flags: cf,
				},
				mem: ram("A4"),
			},
			&Computer{
				cpu: cpu{
					A:     0x0A,
					H:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				mem: ram("A4"),
			},
			nil,
		},
		{
			"ANA L",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					L:     0x0A,
					Flags: cf,
				},
				mem: ram("A5"),
			},
			&Computer{
				cpu: cpu{
					A:     0x0A,
					L:     0x0A,
					PC:    0x01,
					Flags: pf,
				},
				mem: ram("A5"),
			},
			nil,
		},
		{
			"CALL adr",
			&Computer{
				cpu: cpu{
					PC: 0x01,
					SP: 0x07,
				},
				mem: ram("00 CD 0A 00 00 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					PC: 0x0A,
					SP: 0x05,
				},
				mem: ram("00 CD 0A 00 00 00 01 00"),
			},
			nil,
		},
		{
			"CMP A",
			&Computer{
				cpu: cpu{
					A: 0x30,
				},
				mem: ram("BF"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("BF"),
			},
			nil,
		},
		{
			"CMP B: generates carry",
			&Computer{
				cpu: cpu{
					A: 0x30,
					B: 0x31,
				},
				mem: ram("B8"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					B:     0x31,
					PC:    0x01,
					Flags: cf | pf | sf,
				},
				mem: ram("B8"),
			},
			nil,
		}, {
			"CMP B",
			&Computer{
				cpu: cpu{
					A: 0x30,
					B: 0x01,
				},
				mem: ram("B8"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					B:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("B8"),
			},
			nil,
		},
		{
			"CMP C",
			&Computer{
				cpu: cpu{
					A: 0x30,
					C: 0x01,
				},
				mem: ram("B9"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					C:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("B9"),
			},
			nil,
		},
		{
			"CMP D",
			&Computer{
				cpu: cpu{
					A: 0x30,
					D: 0x01,
				},
				mem: ram("BA"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					D:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("BA"),
			},
			nil,
		},
		{
			"CMP E",
			&Computer{
				cpu: cpu{
					A: 0x30,
					E: 0x01,
				},
				mem: ram("BB"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					E:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("BB"),
			},
			nil,
		},
		{
			"CMP H",
			&Computer{
				cpu: cpu{
					A: 0x30,
					H: 0x01,
				},
				mem: ram("BC"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					H:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("BC"),
			},
			nil,
		},
		{
			"CMP L",
			&Computer{
				cpu: cpu{
					A: 0x30,
					L: 0x01,
				},
				mem: ram("BD"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					L:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("BD"),
			},
			nil,
		},
		{
			"DAD B",
			&Computer{
				cpu: cpu{
					B: 0x01,
					C: 0x01,
					H: 0x01,
					L: 0x01,
				},
				mem: ram("09"),
			},
			&Computer{
				cpu: cpu{
					B:     0x01,
					C:     0x01,
					H:     0x02,
					L:     0x02,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("09"),
			},
			nil,
		},
		{
			"DAD B: generates carry",
			&Computer{
				cpu: cpu{
					B: 0xFF,
					C: 0xFE,
					H: 0x00,
					L: 0x03,
				},
				mem: ram("09"),
			},
			&Computer{
				cpu: cpu{
					B:     0xFF,
					C:     0xFE,
					H:     0x00,
					L:     0x01,
					PC:    0x01,
					Flags: cf,
				},
				mem: ram("09"),
			},
			nil,
		},
		{
			"DAD D",
			&Computer{
				cpu: cpu{
					D: 0x01,
					E: 0xFF,
					H: 0x01,
					L: 0x01,
				},
				mem: ram("19"),
			},
			&Computer{
				cpu: cpu{
					D:     0x01,
					E:     0xFF,
					H:     0x03,
					L:     0x00,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("19"),
			},
			nil,
		},
		{
			"DAD H",
			&Computer{
				cpu: cpu{
					H: 0x01,
					L: 0x01,
				},
				mem: ram("29"),
			},
			&Computer{
				cpu: cpu{
					H:     0x02,
					L:     0x02,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("29"),
			},
			nil,
		},
		{
			"DAD SP",
			&Computer{
				cpu: cpu{
					H:  0x01,
					L:  0x03,
					SP: 0x0FFF,
				},
				mem: ram("39"),
			},
			&Computer{
				cpu: cpu{
					H:     0x11,
					L:     0x02,
					SP:    0x0FFF,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("39"),
			},
			nil,
		},
		{
			"DCR A",
			&Computer{
				cpu: cpu{
					A: 0x02,
				},
				mem: ram("3D"),
			},
			&Computer{
				cpu: cpu{
					A:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("3D"),
			},
			nil,
		},
		{
			"DCR B",
			&Computer{
				cpu: cpu{
					B: 0x02,
				},
				mem: ram("05"),
			},
			&Computer{
				cpu: cpu{
					B:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("05"),
			},
			nil,
		},
		{
			"DCR B: carry not set when there's borrow",
			&Computer{
				cpu: cpu{
					B: 0x00,
				},
				mem: ram("05"),
			},
			&Computer{
				cpu: cpu{
					B:     0xff,
					PC:    0x01,
					Flags: sf | pf,
				},
				mem: ram("05"),
			},
			nil,
		},
		{
			"DCR B: carry not modified when there was existing carry",
			&Computer{
				cpu: cpu{
					B:     0x00,
					Flags: cf,
				},
				mem: ram("05"),
			},
			&Computer{
				cpu: cpu{
					B:     0xff,
					PC:    0x01,
					Flags: sf | pf | cf,
				},
				mem: ram("05"),
			},
			nil,
		},
		{
			"DCR B: Generates auxiliary carry when there's carry in the lower nibble",
			&Computer{
				cpu: cpu{
					B: 0x1D,
				},
				mem: ram("05"),
			},
			&Computer{
				cpu: cpu{
					B:     0x1C,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("05"),
			},
			nil,
		},
		{
			"DCR C",
			&Computer{
				cpu: cpu{
					C: 0x02,
				},
				mem: ram("0D"),
			},
			&Computer{
				cpu: cpu{
					C:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("0D"),
			},
			nil,
		},
		{
			"DCR D",
			&Computer{
				cpu: cpu{
					D: 0x02,
				},
				mem: ram("15"),
			},
			&Computer{
				cpu: cpu{
					D:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("15"),
			},
			nil,
		},
		{
			"DCR E",
			&Computer{
				cpu: cpu{
					E: 0x02,
				},
				mem: ram("1D"),
			},
			&Computer{
				cpu: cpu{
					E:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("1D"),
			},
			nil,
		},
		{
			"DCR H",
			&Computer{
				cpu: cpu{
					H: 0x02,
				},
				mem: ram("20"),
			},
			&Computer{
				cpu: cpu{
					H:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("20"),
			},
			nil,
		},
		{
			"DCR L",
			&Computer{
				cpu: cpu{
					L: 0x02,
				},
				mem: ram("2D"),
			},
			&Computer{
				cpu: cpu{
					L:     0x01,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("2D"),
			},
			nil,
		},
		{
			"INR A",
			&Computer{
				cpu: cpu{
					A: 0xFF,
				},
				mem: ram("3C"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("3C"),
			},
			nil,
		},
		{
			"INR B",
			&Computer{
				cpu: cpu{
					B: 0xFF,
				},
				mem: ram("04"),
			},
			&Computer{
				cpu: cpu{
					B:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("04"),
			},
			nil,
		},
		{
			"INR B: generates auxiliary carry",
			&Computer{
				cpu: cpu{
					B: 0x0F,
				},
				mem: ram("04"),
			},
			&Computer{
				cpu: cpu{
					B:     0x10,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("04"),
			},
			nil,
		},
		{
			"INR C",
			&Computer{
				cpu: cpu{
					C: 0x0f,
				},
				mem: ram("0C"),
			},
			&Computer{
				cpu: cpu{
					C:     0x10,
					PC:    0x01,
					Flags: acf,
				},
				mem: ram("0C"),
			},
			nil,
		},
		{
			"INR D",
			&Computer{
				cpu: cpu{
					D: 0x03,
				},
				mem: ram("14"),
			},
			&Computer{
				cpu: cpu{
					D:     0x04,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("14"),
			},
			nil,
		},
		{
			"INR E",
			&Computer{
				cpu: cpu{
					E: 0x03,
				},
				mem: ram("1C"),
			},
			&Computer{
				cpu: cpu{
					E:     0x04,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("1C"),
			},
			nil,
		},
		{
			"INR H",
			&Computer{
				cpu: cpu{
					H: 0x03,
				},
				mem: ram("24"),
			},
			&Computer{
				cpu: cpu{
					H:     0x04,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("24"),
			},
			nil,
		},
		{
			"INR L",
			&Computer{
				cpu: cpu{
					L: 0x03,
				},
				mem: ram("2C"),
			},
			&Computer{
				cpu: cpu{
					L:     0x04,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("2C"),
			},
			nil,
		},
		{
			"INX B",
			&Computer{
				cpu: cpu{
					C: 0xFF,
				},
				mem: ram("03"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					C:  0x00,
					PC: 0x01,
				},
				mem: ram("03"),
			},
			nil,
		},
		{
			"INX D",
			&Computer{
				cpu: cpu{
					E: 0xFF,
				},
				mem: ram("13"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					E:  0x00,
					PC: 0x01,
				},
				mem: ram("13"),
			},
			nil,
		},
		{
			"INX H",
			&Computer{
				cpu: cpu{
					L: 0xFF,
				},
				mem: ram("23"),
			},
			&Computer{
				cpu: cpu{
					H:  0x01,
					L:  0x00,
					PC: 0x01,
				},
				mem: ram("23"),
			},
			nil,
		},
		{
			"INX SP",
			&Computer{
				cpu: cpu{
					SP: 0x0F,
				},
				mem: ram("33"),
			},
			&Computer{
				cpu: cpu{
					PC: 0x01,
					SP: 0x10,
				},
				mem: ram("33"),
			},
			nil,
		},
		{
			"JNZ adr: zero flag set",
			&Computer{
				cpu: cpu{
					Flags: zf,
				},
				mem: ram("C2 0A 00"),
			},
			&Computer{
				cpu: cpu{
					PC:    0x03,
					Flags: zf,
				},
				mem: ram("C2 0A 00"),
			},
			nil,
		},
		{
			"JNZ adr: zero flag not set",
			&Computer{
				mem: ram("C2 0A 00"),
			},
			&Computer{
				cpu: cpu{
					PC: 0x0A,
				},
				mem: ram("C2 0A 00"),
			},
			nil,
		},
		{
			"JMP adr",
			&Computer{
				mem: ram("C3 0A 00"),
			},
			&Computer{
				cpu: cpu{
					PC: 0x0A,
				},
				mem: ram("C3 0A 00"),
			},
			nil,
		},
		{
			"LDAX B",
			&Computer{
				cpu: cpu{
					B: 0x00,
					C: 0x02,
				},
				mem: ram("0a 00 ff"),
			},
			&Computer{
				cpu: cpu{
					A:  0xFF,
					B:  0x00,
					C:  0x02,
					PC: 0x01,
				},
				mem: ram("0a 00 ff"),
			},
			nil,
		},
		{
			"LDAX D",
			&Computer{
				cpu: cpu{
					D: 0x00,
					E: 0x02,
				},
				mem: ram("1a 00 ff"),
			},
			&Computer{
				cpu: cpu{
					A:  0xFF,
					D:  0x00,
					E:  0x02,
					PC: 0x01,
				},
				mem: ram("1a 00 ff"),
			},
			nil,
		},
		{
			"LXI B, D16",
			&Computer{
				mem: ram("01 0B 01"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					C:  0x0B,
					PC: 0x03,
				},
				mem: ram("01 0B 01"),
			},
			nil,
		},
		{
			"LXI D, D16",
			&Computer{
				mem: ram("11 0B 01"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					E:  0x0B,
					PC: 0x03,
				},
				mem: ram("11 0B 01"),
			},
			nil,
		},
		{
			"LXI H, D16",
			&Computer{
				mem: ram("21 0B 01"),
			},
			&Computer{
				cpu: cpu{
					H:  0x01,
					L:  0x0B,
					PC: 0x03,
				},
				mem: ram("21 0B 01"),
			},
			nil,
		},
		{
			"LXI SP, D16",
			&Computer{
				mem: ram("31 0B 01"),
			},
			&Computer{
				cpu: cpu{
					PC: 0x03,
					SP: 0x010B,
				},
				mem: ram("31 0B 01"),
			},
			nil,
		},
		{
			"MOV A, A",
			&Computer{
				cpu: cpu{
					A: 0x01,
				},
				mem: ram("7F"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					PC: 0x01,
				},
				mem: ram("7F"),
			},
			nil,
		},
		{
			"MOV A, B",
			&Computer{
				cpu: cpu{
					B: 0x01,
				},
				mem: ram("78"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					B:  0x01,
					PC: 0x01,
				},
				mem: ram("78"),
			},
			nil,
		},
		{
			"MOV A, C",
			&Computer{
				cpu: cpu{
					C: 0x01,
				},
				mem: ram("79"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					C:  0x01,
					PC: 0x01,
				},
				mem: ram("79"),
			},
			nil,
		},
		{
			"MOV A, D",
			&Computer{
				cpu: cpu{
					D: 0x01,
				},
				mem: ram("7A"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				mem: ram("7A"),
			},
			nil,
		},
		{
			"MOV A, E",
			&Computer{
				cpu: cpu{
					E: 0x01,
				},
				mem: ram("7B"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				mem: ram("7B"),
			},
			nil,
		},
		{
			"MOV A, H",
			&Computer{
				cpu: cpu{
					H: 0x01,
				},
				mem: ram("7C"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("7C"),
			},
			nil,
		},
		{
			"MOV A, L",
			&Computer{
				cpu: cpu{
					L: 0x01,
				},
				mem: ram("7D"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("7D"),
			},
			nil,
		},
		{
			"MOV A, M",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x02,
				},
				mem: ram("7E 00 FF"),
			},
			&Computer{
				cpu: cpu{
					A:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				mem: ram("7E 00 FF"),
			},
			nil,
		},
		{
			"MOV B, A",
			&Computer{
				cpu: cpu{
					A: 0x01,
				},
				mem: ram("47"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					B:  0x01,
					PC: 0x01,
				},
				mem: ram("47"),
			},
			nil,
		},
		{
			"MOV B, B",
			&Computer{
				cpu: cpu{
					B: 0x01,
				},
				mem: ram("40"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					PC: 0x01,
				},
				mem: ram("40"),
			},
			nil,
		},
		{
			"MOV B, C",
			&Computer{
				cpu: cpu{
					C: 0x01,
				},
				mem: ram("41"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					C:  0x01,
					PC: 0x01,
				},
				mem: ram("41"),
			},
			nil,
		},
		{
			"MOV B, D",
			&Computer{
				cpu: cpu{
					D: 0x01,
				},
				mem: ram("42"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				mem: ram("42"),
			},
			nil,
		},
		{
			"MOV B, E",
			&Computer{
				cpu: cpu{
					E: 0x01,
				},
				mem: ram("43"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				mem: ram("43"),
			},
			nil,
		},
		{
			"MOV B, H",
			&Computer{
				cpu: cpu{
					H: 0x01,
				},
				mem: ram("44"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("44"),
			},
			nil,
		},
		{
			"MOV B, L",
			&Computer{
				cpu: cpu{
					L: 0x01,
				},
				mem: ram("45"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("45"),
			},
			nil,
		},
		{
			"MOV B, M",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x02,
				},
				mem: ram("46 00 FF"),
			},
			&Computer{
				cpu: cpu{
					B:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				mem: ram("46 00 FF"),
			},
			nil,
		},
		{
			"MOV C, A",
			&Computer{
				cpu: cpu{
					A: 0x01,
				},
				mem: ram("4F"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					C:  0x01,
					PC: 0x01,
				},
				mem: ram("4F"),
			},
			nil,
		},
		{
			"MOV C, B",
			&Computer{
				cpu: cpu{
					B: 0x01,
				},
				mem: ram("48"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					C:  0x01,
					PC: 0x01,
				},
				mem: ram("48"),
			},
			nil,
		},
		{
			"MOV C, C",
			&Computer{
				cpu: cpu{
					C: 0x01,
				},
				mem: ram("49"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					PC: 0x01,
				},
				mem: ram("49"),
			},
			nil,
		},
		{
			"MOV C, D",
			&Computer{
				cpu: cpu{
					D: 0x01,
				},
				mem: ram("4A"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				mem: ram("4A"),
			},
			nil,
		},
		{
			"MOV C, E",
			&Computer{
				cpu: cpu{
					E: 0x01,
				},
				mem: ram("4B"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				mem: ram("4B"),
			},
			nil,
		},
		{
			"MOV C, H",
			&Computer{
				cpu: cpu{
					H: 0x01,
				},
				mem: ram("4C"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("4C"),
			},
			nil,
		},
		{
			"MOV C, L",
			&Computer{
				cpu: cpu{
					L: 0x01,
				},
				mem: ram("4D"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("4D"),
			},
			nil,
		},
		{
			"MOV C, M",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x02,
				},
				mem: ram("4E 00 FF"),
			},
			&Computer{
				cpu: cpu{
					C:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				mem: ram("4E 00 FF"),
			},
			nil,
		},
		{
			"MOV D, A",
			&Computer{
				cpu: cpu{
					A: 0x01,
				},
				mem: ram("57"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				mem: ram("57"),
			},
			nil,
		},
		{
			"MOV D, B",
			&Computer{
				cpu: cpu{
					B: 0x01,
				},
				mem: ram("50"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				mem: ram("50"),
			},
			nil,
		},
		{
			"MOV D, C",
			&Computer{
				cpu: cpu{
					C: 0x01,
				},
				mem: ram("51"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					D:  0x01,
					PC: 0x01,
				},
				mem: ram("51"),
			},
			nil,
		},
		{
			"MOV D, D",
			&Computer{
				cpu: cpu{
					D: 0x01,
				},
				mem: ram("52"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					PC: 0x01,
				},
				mem: ram("52"),
			},
			nil,
		},
		{
			"MOV D, E",
			&Computer{
				cpu: cpu{
					E: 0x01,
				},
				mem: ram("53"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				mem: ram("53"),
			},
			nil,
		},
		{
			"MOV D, H",
			&Computer{
				cpu: cpu{
					H: 0x01,
				},
				mem: ram("54"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("54"),
			},
			nil,
		},
		{
			"MOV D, L",
			&Computer{
				cpu: cpu{
					L: 0x01,
				},
				mem: ram("55"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("55"),
			},
			nil,
		},
		{
			"MOV D, M",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x02,
				},
				mem: ram("56 00 FF"),
			},
			&Computer{
				cpu: cpu{
					D:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				mem: ram("56 00 FF"),
			},
			nil,
		},
		{
			"MOV E, A",
			&Computer{
				cpu: cpu{
					A: 0x01,
				},
				mem: ram("5F"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				mem: ram("5F"),
			},
			nil,
		},
		{
			"MOV E, B",
			&Computer{
				cpu: cpu{
					B: 0x01,
				},
				mem: ram("58"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				mem: ram("58"),
			},
			nil,
		},
		{
			"MOV E, C",
			&Computer{
				cpu: cpu{
					C: 0x01,
				},
				mem: ram("59"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				mem: ram("59"),
			},
			nil,
		},
		{
			"MOV E, D",
			&Computer{
				cpu: cpu{
					D: 0x01,
				},
				mem: ram("5A"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					E:  0x01,
					PC: 0x01,
				},
				mem: ram("5A"),
			},
			nil,
		},
		{
			"MOV E, E",
			&Computer{
				cpu: cpu{
					E: 0x01,
				},
				mem: ram("5B"),
			},
			&Computer{
				cpu: cpu{
					PC: 0x01,
					E:  0x01,
				},
				mem: ram("5B"),
			},
			nil,
		},
		{
			"MOV E, H",
			&Computer{
				cpu: cpu{
					H: 0x01,
				},
				mem: ram("5C"),
			},
			&Computer{
				cpu: cpu{
					E:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("5C"),
			},
			nil,
		},
		{
			"MOV E, L",
			&Computer{
				cpu: cpu{
					L: 0x01,
				},
				mem: ram("5D"),
			},
			&Computer{
				cpu: cpu{
					E:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("5D"),
			},
			nil,
		},
		{
			"MOV E, M",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x02,
				},
				mem: ram("5E 00 FF"),
			},
			&Computer{
				cpu: cpu{
					E:  0xFF,
					H:  0x00,
					L:  0x02,
					PC: 0x01,
				},
				mem: ram("5E 00 FF"),
			},
			nil,
		},
		{
			"MOV H, A",
			&Computer{
				cpu: cpu{
					A: 0x01,
				},
				mem: ram("67"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("67"),
			},
			nil,
		},
		{
			"MOV H, B",
			&Computer{
				cpu: cpu{
					B: 0x01,
				},
				mem: ram("60"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("60"),
			},
			nil,
		},
		{
			"MOV H, C",
			&Computer{
				cpu: cpu{
					C: 0x01,
				},
				mem: ram("61"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("61"),
			},
			nil,
		},
		{
			"MOV H, D",
			&Computer{
				cpu: cpu{
					D: 0x01,
				},
				mem: ram("62"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("62"),
			},
			nil,
		},
		{
			"MOV H, E",
			&Computer{
				cpu: cpu{
					E: 0x01,
				},
				mem: ram("63"),
			},
			&Computer{
				cpu: cpu{
					E:  0x01,
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("63"),
			},
			nil,
		},
		{
			"MOV H, H",
			&Computer{
				cpu: cpu{
					H: 0x01,
				},
				mem: ram("64"),
			},
			&Computer{
				cpu: cpu{
					H:  0x01,
					PC: 0x01,
				},
				mem: ram("64"),
			},
			nil,
		},
		{
			"MOV H, L",
			&Computer{
				cpu: cpu{
					L: 0x01,
				},
				mem: ram("65"),
			},
			&Computer{
				cpu: cpu{
					H:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("65"),
			},
			nil,
		},
		{
			"MOV H, M",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x02,
				},
				mem: ram("66 00 FF"),
			},
			&Computer{
				cpu: cpu{
					H:  0xFF,
					L:  0x02,
					PC: 0x01,
				},
				mem: ram("66 00 FF"),
			},
			nil,
		},
		{
			"MOV L, A",
			&Computer{
				cpu: cpu{
					A: 0x01,
				},
				mem: ram("6F"),
			},
			&Computer{
				cpu: cpu{
					A:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("6F"),
			},
			nil,
		},
		{
			"MOV L, B",
			&Computer{
				cpu: cpu{
					B: 0x01,
				},
				mem: ram("68"),
			},
			&Computer{
				cpu: cpu{
					B:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("68"),
			},
			nil,
		},
		{
			"MOV L, C",
			&Computer{
				cpu: cpu{
					C: 0x01,
				},
				mem: ram("69"),
			},
			&Computer{
				cpu: cpu{
					C:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("69"),
			},
			nil,
		},
		{
			"MOV L, D",
			&Computer{
				cpu: cpu{
					D: 0x01,
				},
				mem: ram("6A"),
			},
			&Computer{
				cpu: cpu{
					D:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("6A"),
			},
			nil,
		},
		{
			"MOV L, E",
			&Computer{
				cpu: cpu{
					E: 0x01,
				},
				mem: ram("6B"),
			},
			&Computer{
				cpu: cpu{
					E:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("6B"),
			},
			nil,
		},
		{
			"MOV L, H",
			&Computer{
				cpu: cpu{
					H: 0x01,
				},
				mem: ram("6C"),
			},
			&Computer{
				cpu: cpu{
					H:  0x01,
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("6C"),
			},
			nil,
		},
		{
			"MOV L, L",
			&Computer{
				cpu: cpu{
					L: 0x01,
				},
				mem: ram("6D"),
			},
			&Computer{
				cpu: cpu{
					L:  0x01,
					PC: 0x01,
				},
				mem: ram("6D"),
			},
			nil,
		},
		{
			"MOV L, M",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x02,
				},
				mem: ram("6E 00 FF"),
			},
			&Computer{
				cpu: cpu{
					H:  0x00,
					L:  0xFF,
					PC: 0x01,
				},
				mem: ram("6E 00 FF"),
			},
			nil,
		},
		{
			"MOV M, A",
			&Computer{
				cpu: cpu{
					A: 0xff,
					H: 0x00,
					L: 0x03,
				},
				mem: ram("77 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					A:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				mem: ram("77 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, B",
			&Computer{
				cpu: cpu{
					B: 0xff,
					H: 0x00,
					L: 0x03,
				},
				mem: ram("70 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					B:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				mem: ram("70 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, C",
			&Computer{
				cpu: cpu{
					C: 0xff,
					H: 0x00,
					L: 0x03,
				},
				mem: ram("71 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					C:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				mem: ram("71 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, D",
			&Computer{
				cpu: cpu{
					D: 0xff,
					H: 0x00,
					L: 0x03,
				},
				mem: ram("72 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					D:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				mem: ram("72 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, E",
			&Computer{
				cpu: cpu{
					E: 0xff,
					H: 0x00,
					L: 0x03,
				},
				mem: ram("73 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					E:  0xff,
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				mem: ram("73 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, H",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x03,
				},
				mem: ram("74 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				mem: ram("74 00 00 00"),
			},
			nil,
		},
		{
			"MOV M, L",
			&Computer{
				cpu: cpu{
					H: 0x00,
					L: 0x03,
				},
				mem: ram("75 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					H:  0x00,
					L:  0x03,
					PC: 0x01,
				},
				mem: ram("75 00 00 03"),
			},
			nil,
		},
		{
			"MVI A, D8",
			&Computer{
				mem: ram("3E 0B"),
			},
			&Computer{
				cpu: cpu{
					A:  0x0B,
					PC: 0x02,
				},
				mem: ram("3E 0B"),
			},
			nil,
		},
		{
			"MVI B, D8",
			&Computer{
				mem: ram("06 0B"),
			},
			&Computer{
				cpu: cpu{
					B:  0x0B,
					PC: 0x02,
				},
				mem: ram("06 0B"),
			},
			nil,
		},
		{
			"MVI C, D8",
			&Computer{
				mem: ram("0E 0B"),
			},
			&Computer{
				cpu: cpu{
					C:  0x0B,
					PC: 0x02,
				},
				mem: ram("0E 0B"),
			},
			nil,
		},
		{
			"MVI D, D8",
			&Computer{
				mem: ram("16 0B"),
			},
			&Computer{
				cpu: cpu{
					D:  0x0B,
					PC: 0x02,
				},
				mem: ram("16 0B"),
			},
			nil,
		},
		{
			"MVI E, D8",
			&Computer{
				mem: ram("1E 0B"),
			},
			&Computer{
				cpu: cpu{
					E:  0x0B,
					PC: 0x02,
				},
				mem: ram("1E 0B"),
			},
			nil,
		},
		{
			"MVI H, D8",
			&Computer{
				mem: ram("26 0B"),
			},
			&Computer{
				cpu: cpu{
					H:  0x0B,
					PC: 0x02,
				},
				mem: ram("26 0B"),
			},
			nil,
		},
		{
			"MVI L, D8",
			&Computer{
				mem: ram("2E 0B"),
			},
			&Computer{
				cpu: cpu{
					L:  0x0B,
					PC: 0x02,
				},
				mem: ram("2E 0B"),
			},
			nil,
		},
		{
			"NOP",
			&Computer{
				mem: ram("00"),
			},
			&Computer{
				cpu: cpu{
					PC: 0x01,
				},
				mem: ram("00"),
			},
			nil,
		},
		{
			"ORA A",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					Flags: cf,
				},
				mem: ram("B7"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("B7"),
			},
			nil,
		},
		{
			"ORA B",
			&Computer{
				cpu: cpu{
					A:     0xFF,
					B:     0x0A,
					Flags: cf,
				},
				mem: ram("B0"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					B:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("B0"),
			},
			nil,
		},
		{
			"ORA C",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					C: 0x0A,
				},
				mem: ram("B1"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					C:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("B1"),
			},
			nil,
		},
		{
			"ORA D",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					D: 0x0A,
				},
				mem: ram("B2"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					D:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("B2"),
			},
			nil,
		},
		{
			"ORA E",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					E: 0x0A,
				},
				mem: ram("B3"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					E:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("B3"),
			},
			nil,
		},
		{
			"ORA H",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					H: 0x0A,
				},
				mem: ram("B4"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					H:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("B4"),
			},
			nil,
		},
		{
			"ORA L",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					L: 0x0A,
				},
				mem: ram("B5"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					L:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("B5"),
			},
			nil,
		},
		{
			"SBB A: with borrow",
			&Computer{
				cpu: cpu{
					A:     0x01,
					Flags: cf,
				},
				mem: ram("9F"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					PC:    0x01,
					Flags: sf | pf | cf,
				},
				mem: ram("9F"),
			},
			nil,
		},
		{
			"SBB A: no borrow",
			&Computer{
				cpu: cpu{
					A:     0x01,
					Flags: none,
				},
				mem: ram("9F"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("9F"),
			},
			nil,
		},
		{
			"SBB B",
			&Computer{
				cpu: cpu{
					A:     0x02,
					B:     0x01,
					Flags: cf,
				},
				mem: ram("98"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					B:     0x01,
					PC:    1,
					Flags: zf | pf,
				},
				mem: ram("98"),
			},
			nil,
		},
		{
			"SBB C",
			&Computer{
				cpu: cpu{
					A:     0x00,
					C:     0xFF,
					Flags: cf,
				},
				mem: ram("99"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					C:     0xFF,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("99"),
			},
			nil,
		},
		{
			"SBB D",
			&Computer{
				cpu: cpu{
					A:     0x00,
					D:     0xFF,
					Flags: cf,
				},
				mem: ram("9A"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					D:     0xFF,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("9A"),
			},
			nil,
		},
		{
			"SBB E",
			&Computer{
				cpu: cpu{
					A:     0x00,
					E:     0xFF,
					Flags: cf,
				},
				mem: ram("9B"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					E:     0xFF,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("9B"),
			},
			nil,
		},
		{
			"SBB H",
			&Computer{
				cpu: cpu{
					A:     0x00,
					H:     0xFF,
					Flags: cf,
				},
				mem: ram("9C"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					H:     0xFF,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("9C"),
			},
			nil,
		},
		{
			"SBB L",
			&Computer{
				cpu: cpu{
					A:     0x08,
					L:     0x02,
					Flags: cf,
				},
				mem: ram("9D"),
			},
			&Computer{
				cpu: cpu{
					A:     0x05,
					L:     0x02,
					PC:    0x01,
					Flags: pf,
				},
				mem: ram("9D"),
			},
			nil,
		},
		{
			"STAX B",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					B: 0x00,
					C: 0x10,
				},
				mem: ram("02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					A:  0xFF,
					B:  0x00,
					C:  0x10,
					PC: 0x01,
				},
				mem: ram("02 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF"),
			},
			nil,
		},
		{
			"STAX D",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					D: 0x00,
					E: 0x10,
				},
				mem: ram("12 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					A:  0xFF,
					D:  0x00,
					E:  0x10,
					PC: 0x01,
				},
				mem: ram("12 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 FF"),
			},
			nil,
		},
		{
			"SUB A",
			&Computer{
				cpu: cpu{
					A: 0x30,
				},
				mem: ram("97"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("97"),
			},
			nil,
		},
		{
			"SUB B: substracting 0",
			&Computer{
				cpu: cpu{
					A: 0x00,
					B: 0x00,
				},
				mem: ram("90"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					B:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("90"),
			},
			nil,
		},
		{
			"SUB B: 0x02 - 0x01",
			&Computer{
				cpu: cpu{
					A: 0x02,
					B: 0x01,
				},
				mem: ram("90"),
			},
			&Computer{
				cpu: cpu{
					A:     0x01,
					B:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("90"),
			},
			nil,
		},
		{
			"SUB B: 0x01 - 0x02",
			&Computer{
				cpu: cpu{
					A: 0x01,
					B: 0x02,
				},
				mem: ram("90"),
			},
			&Computer{
				cpu: cpu{
					A:     0xFF,
					B:     0x02,
					PC:    0x01,
					Flags: sf | cf | pf,
				},
				mem: ram("90"),
			},
			nil,
		},
		{
			"SUB C",
			&Computer{
				cpu: cpu{
					A: 0x30,
					C: 0x01,
				},
				mem: ram("91"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					C:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("91"),
			},
			nil,
		},
		{
			"SUB D",
			&Computer{
				cpu: cpu{
					A: 0x30,
					D: 0x01,
				},
				mem: ram("92"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					D:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("92"),
			},
			nil,
		},
		{
			"SUB E",
			&Computer{
				cpu: cpu{
					A: 0x30,
					E: 0x01,
				},
				mem: ram("93"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					E:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("93"),
			},
			nil,
		},
		{
			"SUB H",
			&Computer{
				cpu: cpu{
					A: 0x30,
					H: 0x01,
				},
				mem: ram("94"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					H:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("94"),
			},
			nil,
		},
		{
			"SUB L",
			&Computer{
				cpu: cpu{
					A: 0x30,
					L: 0x01,
				},
				mem: ram("95"),
			},
			&Computer{
				cpu: cpu{
					A:     0x2F,
					L:     0x01,
					PC:    0x01,
					Flags: none,
				},
				mem: ram("95"),
			},
			nil,
		},
		{
			"XRA A",
			&Computer{
				cpu: cpu{
					A: 0xFF,
				},
				mem: ram("A8"),
			},
			&Computer{
				cpu: cpu{
					A:     0x00,
					PC:    0x01,
					Flags: zf | pf,
				},
				mem: ram("A8"),
			},
			nil,
		},
		{
			"XRA B",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					B: 0x0A,
				},
				mem: ram("A9"),
			},
			&Computer{
				cpu: cpu{
					A:     0xF5,
					B:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("A9"),
			},
			nil,
		},
		{
			"XRA C",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					C: 0x0A,
				},
				mem: ram("AA"),
			},
			&Computer{
				cpu: cpu{
					A:     0xF5,
					C:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("AA"),
			},
			nil,
		},
		{
			"XRA D",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					D: 0x0A,
				},
				mem: ram("AB"),
			},
			&Computer{
				cpu: cpu{
					A:     0xF5,
					D:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("AB"),
			},
			nil,
		},
		{
			"XRA E",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					E: 0x0A,
				},
				mem: ram("AC"),
			},
			&Computer{
				cpu: cpu{
					A:     0xF5,
					E:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("AC"),
			},
			nil,
		},
		{
			"XRA H",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					H: 0x0A,
				},
				mem: ram("AD"),
			},
			&Computer{
				cpu: cpu{
					A:     0xF5,
					H:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("AD"),
			},
			nil,
		},
		{
			"XRA L",
			&Computer{
				cpu: cpu{
					A: 0xFF,
					L: 0x0A,
				},
				mem: ram("AF"),
			},
			&Computer{
				cpu: cpu{
					A:     0xF5,
					L:     0x0A,
					PC:    0x01,
					Flags: pf | sf,
				},
				mem: ram("AF"),
			},
			nil,
		},
	} {
		t.Run(tC.desc, func(t *testing.T) {
			gotErr := tC.init.Step()
			if gotErr != tC.wantErr {
				t.Fatalf("got err=%v, want=%v", gotErr, tC.wantErr)
			}
			if !reflect.DeepEqual(*tC.init, *tC.want) {
				t.Fatalf("got: %+v\n want: %+v", *tC.init, *tC.want)
			}
		})
	}
}
