package emu

import (
	"reflect"
	"testing"

	"github.com/miguelff/8080/encoding"
)

func rom(rom string) memory {
	return encoding.HexToBin(rom)
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
					registerArray: registerArray{
						PC: 0x01,
						SP: 0x07,
					},
				},
				mem: rom("00 CD 0A 00 00 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						B: 0xFF,
					},
				},
				mem: rom("03"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						D: 0xFF,
					},
				},
				mem: rom("13"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						L: 0xFF,
					},
				},
				mem: rom("23"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						SP: 0x0F,
					},
				},
				mem: rom("33"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
						D: 0x00,
						E: 0x02,
					},
				},
				mem: rom("1a 00 ff"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
						PC: 0x03,
						SP: 0x010B,
					},
				},
				mem: rom("31 0B 01"),
			},
			nil,
		},
		{
			"MOV M, A",
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
						B: 0xff,
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("70 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						C: 0xff,
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("71 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						D: 0xff,
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("72 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						E: 0xff,
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("73 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("74 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
						H: 0x00,
						L: 0x03,
					},
				},
				mem: rom("75 00 00 00"),
			},
			&Computer{
				cpu: cpu{
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
					registerArray: registerArray{
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
