package emu

import (
	"reflect"
	"testing"

	"github.com/miguelff/8080/encoding"
)

func rom(rom string) memory {
	return encoding.HexToBin(rom)
}

func TestParity8(t *testing.T) {
	for _, tC := range []struct {
		b    byte
		want flags
	}{
		{
			0b00001111,
			p,
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

func TestParity16(t *testing.T) {
	for _, tC := range []struct {
		b    uint16
		want flags
	}{
		{
			0b0000001100001111,
			p,
		},
		{
			0b0000001100001110,
			none,
		},
	} {
		if got := parity16(tC.b); got != tC.want {
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
			"CALL adr",
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						SP: 0x07,
					},
				},
				mem: rom("00 CD 0A 00 00 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						SP: 0x05,
						PC: 0x0A,
					},
				},
				mem: rom("00 CD 0A 00 00 00 01 00"),
			},
			nil,
		},
		{
			"INX B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0xFF,
					},
				},
				mem: rom("03"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						B:  0x00,
					},
				},
				mem: rom("03"),
			},
			nil,
		},
		{
			"INX D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0xFF,
					},
				},
				mem: rom("13"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0x01,
						D:  0x00,
					},
				},
				mem: rom("13"),
			},
			nil,
		},
		{
			"INX H",
			&Computer{
				cpu: cpu{
					registers: registers{
						L: 0xFF,
					},
				},
				mem: rom("23"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						H:  0x01,
						L:  0x00,
					},
				},
				mem: rom("23"),
			},
			nil,
		},
		{
			"INX SP",
			&Computer{
				cpu: cpu{
					registers: registers{
						SP: 0x0F,
					},
				},
				mem: rom("33"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						SP: 0x10,
					},
				},
				mem: rom("33"),
			},
			nil,
		},
		{
			"JMP adr",
			&Computer{
				mem: rom("C3 0A 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x0A,
					},
				},
				mem: rom("C3 0A 00"),
			},
			nil,
		},
		{
			"LDAX D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0x00,
						E: 0x02,
					},
				},
				mem: rom("1a 00 ff"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x00,
						E:  0x02,
					},
					alu: alu{
						A: 0xFF,
					},
				},
				mem: rom("1a 00 ff"),
			},
			nil,
		},
		{
			"LXI B, D16",
			&Computer{
				mem: rom("01 0B 01"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x03,
						B:  0x01,
						C:  0x0B,
					},
				},
				mem: rom("01 0B 01"),
			},
			nil,
		},
		{
			"LXI D, D16",
			&Computer{
				mem: rom("11 0B 01"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x03,
						D:  0x01,
						E:  0x0B,
					},
				},
				mem: rom("11 0B 01"),
			},
			nil,
		},
		{
			"LXI H, D16",
			&Computer{
				mem: rom("21 0B 01"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x03,
						H:  0x01,
						L:  0x0B,
					},
				},
				mem: rom("21 0B 01"),
			},
			nil,
		},
		{
			"LXI SP, D16",
			&Computer{
				mem: rom("31 0B 01"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x03,
						SP: 0x010B,
					},
				},
				mem: rom("31 0B 01"),
			},
			nil,
		},
		{
			"MOV A, A",
			&Computer{
				cpu: cpu{
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("7F"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("7F"),
			},
			nil,
		},
		{
			"MOV A, B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0x01,
					},
				},
				mem: rom("78"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("78"),
			},
			nil,
		},
		{
			"MOV A, C",
			&Computer{
				cpu: cpu{
					registers: registers{
						C: 0x01,
					},
				},
				mem: rom("79"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("79"),
			},
			nil,
		},
		{
			"MOV A, D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0x01,
					},
				},
				mem: rom("7A"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("7A"),
			},
			nil,
		},
		{
			"MOV A, E",
			&Computer{
				cpu: cpu{
					registers: registers{
						E: 0x01,
					},
				},
				mem: rom("7B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("7B"),
			},
			nil,
		},
		{
			"MOV A, H",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x01,
					},
				},
				mem: rom("7C"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						H:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("7C"),
			},
			nil,
		},
		{
			"MOV A, L",
			&Computer{
				cpu: cpu{
					registers: registers{
						L: 0x01,
					},
				},
				mem: rom("7D"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						L:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("7D"),
			},
			nil,
		},
		{
			"MOV B, A",
			&Computer{
				cpu: cpu{
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("47"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("47"),
			},
			nil,
		},
		{
			"MOV B, B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0x01,
					},
				},
				mem: rom("40"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
					},
				},
				mem: rom("40"),
			},
			nil,
		},
		{
			"MOV B, C",
			&Computer{
				cpu: cpu{
					registers: registers{
						C: 0x01,
					},
				},
				mem: rom("41"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						C:  0x01,
					},
				},
				mem: rom("41"),
			},
			nil,
		},
		{
			"MOV B, D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0x01,
					},
				},
				mem: rom("42"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						D:  0x01,
					},
				},
				mem: rom("42"),
			},
			nil,
		},
		{
			"MOV B, E",
			&Computer{
				cpu: cpu{
					registers: registers{
						E: 0x01,
					},
				},
				mem: rom("43"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						E:  0x01,
					},
				},
				mem: rom("43"),
			},
			nil,
		},
		{
			"MOV B, H",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x01,
					},
				},
				mem: rom("44"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						H:  0x01,
					},
				},
				mem: rom("44"),
			},
			nil,
		},
		{
			"MOV B, L",
			&Computer{
				cpu: cpu{
					registers: registers{
						L: 0x01,
					},
				},
				mem: rom("45"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						L:  0x01,
					},
				},
				mem: rom("45"),
			},
			nil,
		},
		{
			"MOV C, A",
			&Computer{
				cpu: cpu{
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("4F"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("4F"),
			},
			nil,
		},
		{
			"MOV C, B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0x01,
					},
				},
				mem: rom("48"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						B:  0x01,
					},
				},
				mem: rom("48"),
			},
			nil,
		},
		{
			"MOV C, C",
			&Computer{
				cpu: cpu{
					registers: registers{
						C: 0x01,
					},
				},
				mem: rom("49"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
					},
				},
				mem: rom("49"),
			},
			nil,
		},
		{
			"MOV C, D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0x01,
					},
				},
				mem: rom("4A"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						D:  0x01,
					},
				},
				mem: rom("4A"),
			},
			nil,
		},
		{
			"MOV C, E",
			&Computer{
				cpu: cpu{
					registers: registers{
						E: 0x01,
					},
				},
				mem: rom("4B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						E:  0x01,
					},
				},
				mem: rom("4B"),
			},
			nil,
		},
		{
			"MOV C, H",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x01,
					},
				},
				mem: rom("4C"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						H:  0x01,
					},
				},
				mem: rom("4C"),
			},
			nil,
		},
		{
			"MOV C, L",
			&Computer{
				cpu: cpu{
					registers: registers{
						L: 0x01,
					},
				},
				mem: rom("4D"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						L:  0x01,
					},
				},
				mem: rom("4D"),
			},
			nil,
		},
		{
			"MOV D, A",
			&Computer{
				cpu: cpu{
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("57"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("57"),
			},
			nil,
		},
		{
			"MOV D, B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0x01,
					},
				},
				mem: rom("50"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						D:  0x01,
					},
				},
				mem: rom("50"),
			},
			nil,
		},
		{
			"MOV D, C",
			&Computer{
				cpu: cpu{
					registers: registers{
						C: 0x01,
					},
				},
				mem: rom("51"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						D:  0x01,
					},
				},
				mem: rom("51"),
			},
			nil,
		},
		{
			"MOV D, D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0x01,
					},
				},
				mem: rom("52"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
					},
				},
				mem: rom("52"),
			},
			nil,
		},
		{
			"MOV D, E",
			&Computer{
				cpu: cpu{
					registers: registers{
						E: 0x01,
					},
				},
				mem: rom("53"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
						E:  0x01,
					},
				},
				mem: rom("53"),
			},
			nil,
		},
		{
			"MOV D, H",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x01,
					},
				},
				mem: rom("54"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
						H:  0x01,
					},
				},
				mem: rom("54"),
			},
			nil,
		},
		{
			"MOV D, L",
			&Computer{
				cpu: cpu{
					registers: registers{
						L: 0x01,
					},
				},
				mem: rom("55"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
						L:  0x01,
					},
				},
				mem: rom("55"),
			},
			nil,
		},
		{
			"MOV E, A",
			&Computer{
				cpu: cpu{
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("5F"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("5F"),
			},
			nil,
		},
		{
			"MOV E, B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0x01,
					},
				},
				mem: rom("58"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						E:  0x01,
					},
				},
				mem: rom("58"),
			},
			nil,
		},
		{
			"MOV E, C",
			&Computer{
				cpu: cpu{
					registers: registers{
						C: 0x01,
					},
				},
				mem: rom("59"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						E:  0x01,
					},
				},
				mem: rom("59"),
			},
			nil,
		},
		{
			"MOV E, D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0x01,
					},
				},
				mem: rom("5A"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
						E:  0x01,
					},
				},
				mem: rom("5A"),
			},
			nil,
		},
		{
			"MOV E, E",
			&Computer{
				cpu: cpu{
					registers: registers{
						E: 0x01,
					},
				},
				mem: rom("5B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0x01,
					},
				},
				mem: rom("5B"),
			},
			nil,
		},
		{
			"MOV E, H",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x01,
					},
				},
				mem: rom("5C"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0x01,
						H:  0x01,
					},
				},
				mem: rom("5C"),
			},
			nil,
		},
		{
			"MOV E, L",
			&Computer{
				cpu: cpu{
					registers: registers{
						L: 0x01,
					},
				},
				mem: rom("5D"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0x01,
						L:  0x01,
					},
				},
				mem: rom("5D"),
			},
			nil,
		},
		{
			"MOV H, A",
			&Computer{
				cpu: cpu{
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("67"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						H:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("67"),
			},
			nil,
		},
		{
			"MOV H, B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0x01,
					},
				},
				mem: rom("60"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						H:  0x01,
					},
				},
				mem: rom("60"),
			},
			nil,
		},
		{
			"MOV H, C",
			&Computer{
				cpu: cpu{
					registers: registers{
						C: 0x01,
					},
				},
				mem: rom("61"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						H:  0x01,
					},
				},
				mem: rom("61"),
			},
			nil,
		},
		{
			"MOV H, D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0x01,
					},
				},
				mem: rom("62"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
						H:  0x01,
					},
				},
				mem: rom("62"),
			},
			nil,
		},
		{
			"MOV H, E",
			&Computer{
				cpu: cpu{
					registers: registers{
						E: 0x01,
					},
				},
				mem: rom("63"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0x01,
						H:  0x01,
					},
				},
				mem: rom("63"),
			},
			nil,
		},
		{
			"MOV H, H",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x01,
					},
				},
				mem: rom("64"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						H:  0x01,
					},
				},
				mem: rom("64"),
			},
			nil,
		},
		{
			"MOV H, L",
			&Computer{
				cpu: cpu{
					registers: registers{
						L: 0x01,
					},
				},
				mem: rom("65"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						H:  0x01,
						L:  0x01,
					},
				},
				mem: rom("65"),
			},
			nil,
		},
		{
			"MOV L, A",
			&Computer{
				cpu: cpu{
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("6F"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						L:  0x01,
					},
					alu: alu{
						A: 0x01,
					},
				},
				mem: rom("6F"),
			},
			nil,
		},
		{
			"MOV L, B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0x01,
					},
				},
				mem: rom("68"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0x01,
						L:  0x01,
					},
				},
				mem: rom("68"),
			},
			nil,
		},
		{
			"MOV L, C",
			&Computer{
				cpu: cpu{
					registers: registers{
						C: 0x01,
					},
				},
				mem: rom("69"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0x01,
						L:  0x01,
					},
				},
				mem: rom("69"),
			},
			nil,
		},
		{
			"MOV L, D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0x01,
					},
				},
				mem: rom("6A"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0x01,
						L:  0x01,
					},
				},
				mem: rom("6A"),
			},
			nil,
		},
		{
			"MOV L, E",
			&Computer{
				cpu: cpu{
					registers: registers{
						E: 0x01,
					},
				},
				mem: rom("6B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0x01,
						L:  0x01,
					},
				},
				mem: rom("6B"),
			},
			nil,
		},
		{
			"MOV L, H",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x01,
					},
				},
				mem: rom("6C"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						L:  0x01,
						H:  0x01,
					},
				},
				mem: rom("6C"),
			},
			nil,
		},
		{
			"MOV L, L",
			&Computer{
				cpu: cpu{
					registers: registers{
						L: 0x01,
					},
				},
				mem: rom("6D"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						L:  0x01,
					},
				},
				mem: rom("6D"),
			},
			nil,
		},
		{
			"MOV M, A",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x00,
						L: 0x03,
					},
					alu: alu{
						A: 0xff,
					},
				},
				mem: rom("77 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						H:  0x00,
						L:  0x03,
					},
					alu: alu{
						A: 0xff,
					},
				},
				mem: rom("77 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, B",
			&Computer{
				cpu: cpu{
					registers: registers{
						B: 0xff,
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("70 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						B:  0xff,
						H:  0x00,
						L:  0x03,
					},
				},
				mem: rom("70 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, C",
			&Computer{
				cpu: cpu{
					registers: registers{
						C: 0xff,
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("71 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						C:  0xff,
						H:  0x00,
						L:  0x03,
					},
				},
				mem: rom("71 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, D",
			&Computer{
				cpu: cpu{
					registers: registers{
						D: 0xff,
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("72 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						D:  0xff,
						H:  0x00,
						L:  0x03,
					},
				},
				mem: rom("72 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, E",
			&Computer{
				cpu: cpu{
					registers: registers{
						E: 0xff,
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("73 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						E:  0xff,
						H:  0x00,
						L:  0x03,
					},
				},
				mem: rom("73 00 00 FF"),
			},
			nil,
		},
		{
			"MOV M, H",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("74 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						H:  0x00,
						L:  0x03,
					},
				},
				mem: rom("74 00 00 00"),
			},
			nil,
		},
		{
			"MOV M, L",
			&Computer{
				cpu: cpu{
					registers: registers{
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("75 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
						H:  0x00,
						L:  0x03,
					},
				},
				mem: rom("75 00 00 03"),
			},
			nil,
		},
		{
			"MVI A, D8",
			&Computer{
				mem: rom("3E 0B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x02,
					},
					alu: alu{
						A: 0x0B,
					},
				},
				mem: rom("3E 0B"),
			},
			nil,
		},
		{
			"MVI B, D8",
			&Computer{
				mem: rom("06 0B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x02,
						B:  0x0B,
					},
				},
				mem: rom("06 0B"),
			},
			nil,
		},
		{
			"MVI C, D8",
			&Computer{
				mem: rom("0E 0B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x02,
						C:  0x0B,
					},
				},
				mem: rom("0E 0B"),
			},
			nil,
		},
		{
			"MVI D, D8",
			&Computer{
				mem: rom("16 0B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x02,
						D:  0x0B,
					},
				},
				mem: rom("16 0B"),
			},
			nil,
		},
		{
			"MVI E, D8",
			&Computer{
				mem: rom("1E 0B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x02,
						E:  0x0B,
					},
				},
				mem: rom("1E 0B"),
			},
			nil,
		},
		{
			"MVI H, D8",
			&Computer{
				mem: rom("26 0B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x02,
						H:  0x0B,
					},
				},
				mem: rom("26 0B"),
			},
			nil,
		},
		{
			"MVI L, D8",
			&Computer{
				mem: rom("2E 0B"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x02,
						L:  0x0B,
					},
				},
				mem: rom("2E 0B"),
			},
			nil,
		},
		{
			"NOP",
			&Computer{
				mem: rom("00"),
			},
			&Computer{
				cpu: cpu{
					registers: registers{
						PC: 0x01,
					},
				},
				mem: rom("00"),
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
