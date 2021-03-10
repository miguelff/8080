// dasm.go provides tooling for disassembling 8080 machine code
package dasm

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"strings"
)

type byteReader struct {
	*bufio.Reader
	cursor int
}

// ReadByte reads a single byte from the reader, keeping track of the number of bytes read
func (r *byteReader) ReadByte() (byte, error) {
	b, err := r.Reader.ReadByte()
	if err == nil {
		r.cursor++
	}
	return b, err
}

func newByteReader(r io.Reader, cursor int) *byteReader {
	return &byteReader{Reader: bufio.NewReader(r), cursor: cursor}
}

type instDasm func(*byteReader, *bufio.Writer) error

var instructions = []instDasm{
	0x00: inst8("NOP"),
	0x01: inst24("LXI B,D16"),
	0x02: inst8("STAX B"),
	0x03: inst8("INX B"),
	0x04: inst8("INR B"),
	0x05: inst8("DCR B"),
	0x06: inst16("MVI B,D8"),
	0x07: inst8("RLC"),
	0x09: inst8("DAD B"),
	0x0a: inst8("LDAX B"),
	0x0b: inst8("DCX B"),
	0x0c: inst8("INR C"),
	0x0d: inst8("DCR C"),
	0x0e: inst16("MVI C,D8"),
	0x0f: inst8("RRC"),
	0x11: inst24("LXI D,D16"),
	0x12: inst8("STAX D"),
	0x13: inst8("INX D"),
	0x14: inst8("INR D"),
	0x15: inst8("DCR D"),
	0x16: inst16("MVI D,D8"),
	0x17: inst8("RAL"),
	0x19: inst8("DAD D"),
	0x1a: inst8("LDAX D"),
	0x1b: inst8("DCX D"),
	0x1c: inst8("INR E"),
	0x1d: inst8("DCR E"),
	0x1e: inst16("MVI E,D8"),
	0x1f: inst8("RAR"),
	0x21: inst24("LXI H,D16"),
	0x22: inst24("SHLD adr"),
	0x23: inst8("INX H"),
	0x24: inst8("INR H"),
	0x25: inst8("DCR H"),
	0x26: inst16("MVI H,D8"),
	0x27: inst8("DAA"),
	0x29: inst8("DAD H"),
	0x2a: inst24("LHLD adr"),
	0x2b: inst8("DCX H"),
	0x2c: inst8("INR L"),
	0x2d: inst8("DCR L"),
	0x2e: inst16("MVI L,D8"),
	0x2f: inst8("CMA"),
	0x31: inst24("LXI SP,D16"),
	0x32: inst24("STA adr"),
	0x33: inst8("INX SP"),
	0x34: inst8("INR M"),
	0x35: inst8("DCR M"),
	0x36: inst16("MVI M,D8"),
	0x37: inst8("STC"),
	0x39: inst8("DAD SP"),
	0x3a: inst24("LDA adr"),
	0x3b: inst8("DCX SP"),
	0x3c: inst8("INR A"),
	0x3d: inst8("DCR A"),
	0x3e: inst16("MVI A,D8"),
	0x3f: inst8("CMC"),
	0x40: inst8("MOV B,B"),
	0x41: inst8("MOV B,C"),
	0x42: inst8("MOV B,D"),
	0x43: inst8("MOV B,E"),
	0x44: inst8("MOV B,H"),
	0x45: inst8("MOV B,L"),
	0x46: inst8("MOV B,M"),
	0x47: inst8("MOV B,A"),
	0x48: inst8("MOV C,B"),
	0x49: inst8("MOV C,C"),
	0x4a: inst8("MOV C,D"),
	0x4b: inst8("MOV C,E"),
	0x4c: inst8("MOV C,H"),
	0x4d: inst8("MOV C,L"),
	0x4e: inst8("MOV C,M"),
	0x4f: inst8("MOV C,A"),
	0x50: inst8("MOV D,B"),
	0x51: inst8("MOV D,C"),
	0x52: inst8("MOV D,D"),
	0x53: inst8("MOV D,E"),
	0x54: inst8("MOV D,H"),
	0x55: inst8("MOV D,L"),
	0x56: inst8("MOV D,M"),
	0x57: inst8("MOV D,A"),
	0x58: inst8("MOV E,B"),
	0x59: inst8("MOV E,C"),
	0x5a: inst8("MOV E,D"),
	0x5b: inst8("MOV E,E"),
	0x5c: inst8("MOV E,H"),
	0x5d: inst8("MOV E,L"),
	0x5e: inst8("MOV E,M"),
	0x5f: inst8("MOV E,A"),
	0x60: inst8("MOV H,B"),
	0x61: inst8("MOV H,C"),
	0x62: inst8("MOV H,D"),
	0x63: inst8("MOV H,E"),
	0x64: inst8("MOV H,H"),
	0x65: inst8("MOV H,L"),
	0x66: inst8("MOV H,M"),
	0x67: inst8("MOV H,A"),
	0x68: inst8("MOV L,B"),
	0x69: inst8("MOV L,C"),
	0x6a: inst8("MOV L,D"),
	0x6b: inst8("MOV L,E"),
	0x6c: inst8("MOV L,H"),
	0x6d: inst8("MOV L,L"),
	0x6e: inst8("MOV L,M"),
	0x6f: inst8("MOV L,A"),
	0x70: inst8("MOV M,B"),
	0x71: inst8("MOV M,C"),
	0x72: inst8("MOV M,D"),
	0x73: inst8("MOV M,E"),
	0x74: inst8("MOV M,H"),
	0x75: inst8("MOV M,L"),
	0x76: inst8("HLT"),
	0x77: inst8("MOV M,A"),
	0x78: inst8("MOV A,B"),
	0x79: inst8("MOV A,C"),
	0x7a: inst8("MOV A,D"),
	0x7b: inst8("MOV A,E"),
	0x7c: inst8("MOV A,H"),
	0x7d: inst8("MOV A,L"),
	0x7e: inst8("MOV A,M"),
	0x7f: inst8("MOV A,A"),
	0x80: inst8("ADD B"),
	0x81: inst8("ADD C"),
	0x82: inst8("ADD D"),
	0x83: inst8("ADD E"),
	0x84: inst8("ADD H"),
	0x85: inst8("ADD L"),
	0x86: inst8("ADD M"),
	0x87: inst8("ADD A"),
	0x88: inst8("ADC B"),
	0x89: inst8("ADC C"),
	0x8a: inst8("ADC D"),
	0x8b: inst8("ADC E"),
	0x8c: inst8("ADC H"),
	0x8d: inst8("ADC L"),
	0x8e: inst8("ADC M"),
	0x8f: inst8("ADC A"),
	0x90: inst8("SUB B"),
	0x91: inst8("SUB C"),
	0x92: inst8("SUB D"),
	0x93: inst8("SUB E"),
	0x94: inst8("SUB H"),
	0x95: inst8("SUB L"),
	0x96: inst8("SUB M"),
	0x97: inst8("SUB A"),
	0x98: inst8("SBB B"),
	0x99: inst8("SBB C"),
	0x9a: inst8("SBB D"),
	0x9b: inst8("SBB E"),
	0x9c: inst8("SBB H"),
	0x9d: inst8("SBB L"),
	0x9e: inst8("SBB M"),
	0x9f: inst8("SBB A"),
	0xa0: inst8("ANA B"),
	0xa1: inst8("ANA C"),
	0xa2: inst8("ANA D"),
	0xa3: inst8("ANA E"),
	0xa4: inst8("ANA H"),
	0xa5: inst8("ANA L"),
	0xa6: inst8("ANA M"),
	0xa7: inst8("ANA A"),
	0xa8: inst8("XRA B"),
	0xa9: inst8("XRA C"),
	0xaa: inst8("XRA D"),
	0xab: inst8("XRA E"),
	0xac: inst8("XRA H"),
	0xad: inst8("XRA L"),
	0xae: inst8("XRA M"),
	0xaf: inst8("XRA A"),
	0xb0: inst8("ORA B"),
	0xb1: inst8("ORA C"),
	0xb2: inst8("ORA D"),
	0xb3: inst8("ORA E"),
	0xb4: inst8("ORA H"),
	0xb5: inst8("ORA L"),
	0xb6: inst8("ORA M"),
	0xb7: inst8("ORA A"),
	0xb8: inst8("CMP B"),
	0xb9: inst8("CMP C"),
	0xba: inst8("CMP D"),
	0xbb: inst8("CMP E"),
	0xbc: inst8("CMP H"),
	0xbd: inst8("CMP L"),
	0xbe: inst8("CMP M"),
	0xbf: inst8("CMP A"),
	0xc0: inst8("RNZ"),
	0xc1: inst8("POP B"),
	0xc2: inst24("JNZ adr"),
	0xc3: inst24("JMP adr"),
	0xc4: inst24("CNZ adr"),
	0xc5: inst8("PUSH B"),
	0xc6: inst16("ADI D8"),
	0xc7: inst8("RST 0"),
	0xc8: inst8("RZ"),
	0xc9: inst8("RET"),
	0xca: inst24("JZ adr"),
	0xcc: inst24("CZ adr"),
	0xcd: inst24("CALL adr"),
	0xce: inst16("ACI D8"),
	0xcf: inst8("RST 1"),
	0xd0: inst8("RNC"),
	0xd1: inst8("POP D"),
	0xd2: inst24("JNC adr"),
	0xd3: inst16("OUT D8"),
	0xd4: inst24("CNC adr"),
	0xd5: inst8("PUSH D"),
	0xd6: inst16("SUI D8"),
	0xd7: inst8("RST 2"),
	0xd8: inst8("RC"),
	0xda: inst24("JC adr"),
	0xdb: inst16("IN D8"),
	0xdc: inst24("CC adr"),
	0xde: inst16("SBI D8"),
	0xdf: inst8("RST 3"),
	0xe0: inst8("RPO"),
	0xe1: inst8("POP H"),
	0xe2: inst24("JPO adr"),
	0xe3: inst8("XTHL"),
	0xe4: inst24("CPO adr"),
	0xe5: inst8("PUSH H"),
	0xe6: inst16("ANI D8"),
	0xe7: inst8("RST 4"),
	0xe8: inst8("RPE"),
	0xe9: inst8("PCHL"),
	0xea: inst24("JPE adr"),
	0xeb: inst8("XCHG"),
	0xec: inst24("CPE adr"),
	0xee: inst16("XRI D8"),
	0xef: inst8("RST 5"),
	0xf0: inst8("RP"),
	0xf1: inst8("POP PSW"),
	0xf2: inst24("JP adr"),
	0xf3: inst8("DI"),
	0xf4: inst24("CP adr"),
	0xf5: inst8("PUSH PSW"),
	0xf6: inst16("ORI D8"),
	0xf7: inst8("RST 6"),
	0xf8: inst8("RM"),
	0xf9: inst8("SPHL"),
	0xfa: inst24("JM adr"),
	0xfb: inst8("EI"),
	0xfc: inst24("CM adr"),
	0xfe: inst16("CPI D8"),
	0xff: inst8("RST 7"),
}

// Disassemble reads machine code from the reader, and writes assembly code to the writer
func Disassemble(r io.Reader, w io.Writer) error {
	return DisassembleFrom(r, w, 0)
}

// DisassembleFrom reads machine code from the reader starting at the given offset, and writes assembly code to the writer
func DisassembleFrom(r io.Reader, w io.Writer, offset int) error {
	br := newByteReader(r, offset)
	bw := bufio.NewWriter(w)
	defer bw.Flush()

	if _, err := br.Discard(offset); err != nil {
		return err
	}

	lastOp := byte(len(instructions) - 1)
	for op, err := br.ReadByte(); err == nil; op, err = br.ReadByte() {
		if op > lastOp {
			return fmt.Errorf("unkown op code 0x%02X", op)
		}

		if instDasm := instructions[op]; instDasm != nil {
			err = instDasm(br, bw)
			if err != nil {
				return err
			}
		}
	}

	return nil
}

// DisassembleFirst dissasemble a single machine instruction
func DisassembleFirst(bin []byte) (string, error) {
	buf := new(bytes.Buffer)
	bw := bufio.NewWriter(buf)

	br := newByteReader(bytes.NewReader(bin), 0)
	op, err := br.ReadByte()

	if err != nil {
		return "", err
	}

	lastOp := byte(len(instructions) - 1)
	if op > lastOp {
		return "", fmt.Errorf("unkown op code 0x%02X", op)
	}

	if instDasm := instructions[op]; instDasm != nil {
		err = instDasm(br, bw)
		if err != nil {
			return "", err
		}
		err := bw.Flush()
		if err != nil {
			return "", err
		}
		return buf.String(), nil
	} else {
		return "", fmt.Errorf("unkown op code 0x%02X", op)
	}
}

func inst8(inst string) instDasm {
	return func(_ *byteReader, w *bufio.Writer) error {
		_, err := w.WriteString(inst + "\n")
		return err
	}
}

func inst16(inst string) instDasm {
	fmtStr := strings.Replace(inst, "D8", "$%02X\n", 1)

	return func(r *byteReader, w *bufio.Writer) error {
		b, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("%s requires 2 arguments, error parsing arg 1 (low byte) at byte %02X: %w", inst, r.cursor, err)
		}

		_, err = w.WriteString(fmt.Sprintf(fmtStr, b))
		return err
	}
}

func inst24(inst string) instDasm {
	inst = strings.Replace(inst, "adr", "D16", 1)
	fmtStr := strings.Replace(inst, "D16", "$%02X%02X\n", 1)

	return func(r *byteReader, w *bufio.Writer) error {
		var err error
		lb, err := r.ReadByte()

		if err != nil {
			return fmt.Errorf("%s requires 2 arguments, error parsing arg 1 (low byte) at byte %2X: %w", inst, r.cursor, err)
		}
		hb, err := r.ReadByte()
		if err != nil {
			return fmt.Errorf("%s requires 2 arguments, error parsing arg 2 (high byte) at byte %2X: %w", inst, r.cursor, err)
		}

		_, err = w.WriteString(fmt.Sprintf(fmtStr, hb, lb))
		return err
	}
}
