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
