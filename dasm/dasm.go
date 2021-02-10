// dasm.go provides tooling for disassembling 8080 machine code
package dasm

import (
	"bufio"
	"fmt"
	"io"
	"strings"
)

type byteReader struct {
	*bufio.Reader
	cursor int
}

// ReadByte reads a single byte from the reader, keeping
// track of the number of bytes read
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
	0x00: singleWordInst("NOP"),
	0x01: tripleWordInst("LXI B,D16"),
	0x02: singleWordInst("STAX B"),
	0x03: singleWordInst("INX B"),
	0x04: singleWordInst("INR B"),
	0x05: singleWordInst("DCR B"),
	0x06: doubleWordInst("MVI B,D8"),
	0x07: singleWordInst("RLC"),
	0x09: singleWordInst("DAD B"),
	0x0a: singleWordInst("LDAX B"),
	0x0b: singleWordInst("DCX B"),
	0x0c: singleWordInst("INR C"),
	0x0d: singleWordInst("DCR C"),
	0x0e: doubleWordInst("MVI C,D8"),
	0x0f: singleWordInst("RRC"),
	0x11: tripleWordInst("LXI D,D16"),
	0x12: singleWordInst("STAX D"),
	0x13: singleWordInst("INX D"),
	0x14: singleWordInst("INR D"),
	0x15: singleWordInst("DCR D"),
	0x16: doubleWordInst("MVI D,D8"),
	0x17: singleWordInst("RAL"),
	0x19: singleWordInst("DAD D"),
	0x1a: singleWordInst("LDAX D"),
	0x1b: singleWordInst("DCX D"),
	0x1c: singleWordInst("INR E"),
	0x1d: singleWordInst("DCR E"),
	0x1e: doubleWordInst("MVI E,D8"),
	0x1f: singleWordInst("RAR"),
	0x21: tripleWordInst("LXI H,D16"),
	0x22: tripleWordInst("SHLD adr"),
	0x23: singleWordInst("INX H"),
	0x24: singleWordInst("INR H"),
	0x25: singleWordInst("DCR H"),
	0x26: doubleWordInst("MVI H,D8"),
	0x27: singleWordInst("DAA"),
	0x29: singleWordInst("DAD H"),
	0x2a: tripleWordInst("LHLD adr"),
	0x2b: singleWordInst("DCX H"),
	0x2c: singleWordInst("INR L"),
	0x2d: singleWordInst("DCR L"),
	0x2e: doubleWordInst("MVI L,D8"),
	0x2f: singleWordInst("CMA"),
	0x31: tripleWordInst("LXI SP,D16"),
	0x32: tripleWordInst("STA adr"),
	0x33: singleWordInst("INX SP"),
	0x34: singleWordInst("INR M"),
	0x35: singleWordInst("DCR M"),
	0x36: doubleWordInst("MVI M,D8"),
	0x37: singleWordInst("STC"),
	0x39: singleWordInst("DAD SP"),
	0x3a: tripleWordInst("LDA adr"),
	0x3b: singleWordInst("DCX SP"),
	0x3c: singleWordInst("INR A"),
	0x3d: singleWordInst("DCR A"),
	0x3e: doubleWordInst("MVI A,D8"),
	0x3f: singleWordInst("CMC"),
	0x40: singleWordInst("MOV B,B"),
	0x41: singleWordInst("MOV B,C"),
	0x42: singleWordInst("MOV B,D"),
	0x43: singleWordInst("MOV B,E"),
	0x44: singleWordInst("MOV B,H"),
	0x45: singleWordInst("MOV B,L"),
	0x46: singleWordInst("MOV B,M"),
	0x47: singleWordInst("MOV B,A"),
	0x48: singleWordInst("MOV C,B"),
	0x49: singleWordInst("MOV C,C"),
	0x4a: singleWordInst("MOV C,D"),
	0x4b: singleWordInst("MOV C,E"),
	0x4c: singleWordInst("MOV C,H"),
	0x4d: singleWordInst("MOV C,L"),
	0x4e: singleWordInst("MOV C,M"),
	0x4f: singleWordInst("MOV C,A"),
	0x50: singleWordInst("MOV D,B"),
	0x51: singleWordInst("MOV D,C"),
	0x52: singleWordInst("MOV D,D"),
	0x53: singleWordInst("MOV D,E"),
	0x54: singleWordInst("MOV D,H"),
	0x55: singleWordInst("MOV D,L"),
	0x56: singleWordInst("MOV D,M"),
	0x57: singleWordInst("MOV D,A"),
	0x58: singleWordInst("MOV E,B"),
	0x59: singleWordInst("MOV E,C"),
	0x5a: singleWordInst("MOV E,D"),
	0x5b: singleWordInst("MOV E,E"),
	0x5c: singleWordInst("MOV E,H"),
	0x5d: singleWordInst("MOV E,L"),
	0x5e: singleWordInst("MOV E,M"),
	0x5f: singleWordInst("MOV E,A"),
	0x60: singleWordInst("MOV H,B"),
	0x61: singleWordInst("MOV H,C"),
	0x62: singleWordInst("MOV H,D"),
	0x63: singleWordInst("MOV H,E"),
	0x64: singleWordInst("MOV H,H"),
	0x65: singleWordInst("MOV H,L"),
	0x66: singleWordInst("MOV H,M"),
	0x67: singleWordInst("MOV H,A"),
	0x68: singleWordInst("MOV L,B"),
	0x69: singleWordInst("MOV L,C"),
	0x6a: singleWordInst("MOV L,D"),
	0x6b: singleWordInst("MOV L,E"),
	0x6c: singleWordInst("MOV L,H"),
	0x6d: singleWordInst("MOV L,L"),
	0x6e: singleWordInst("MOV L,M"),
	0x6f: singleWordInst("MOV L,A"),
	0x70: singleWordInst("MOV M,B"),
	0x71: singleWordInst("MOV M,C"),
	0x72: singleWordInst("MOV M,D"),
	0x73: singleWordInst("MOV M,E"),
	0x74: singleWordInst("MOV M,H"),
	0x75: singleWordInst("MOV M,L"),
	0x76: singleWordInst("HLT"),
	0x77: singleWordInst("MOV M,A"),
	0x78: singleWordInst("MOV A,B"),
	0x79: singleWordInst("MOV A,C"),
	0x7a: singleWordInst("MOV A,D"),
	0x7b: singleWordInst("MOV A,E"),
	0x7c: singleWordInst("MOV A,H"),
	0x7d: singleWordInst("MOV A,L"),
	0x7e: singleWordInst("MOV A,M"),
	0x7f: singleWordInst("MOV A,A"),
	0x80: singleWordInst("ADD B"),
	0x81: singleWordInst("ADD C"),
	0x82: singleWordInst("ADD D"),
	0x83: singleWordInst("ADD E"),
	0x84: singleWordInst("ADD H"),
	0x85: singleWordInst("ADD L"),
	0x86: singleWordInst("ADD M"),
	0x87: singleWordInst("ADD A"),
	0x88: singleWordInst("ADC B"),
	0x89: singleWordInst("ADC C"),
	0x8a: singleWordInst("ADC D"),
	0x8b: singleWordInst("ADC E"),
	0x8c: singleWordInst("ADC H"),
	0x8d: singleWordInst("ADC L"),
	0x8e: singleWordInst("ADC M"),
	0x8f: singleWordInst("ADC A"),
	0x90: singleWordInst("SUB B"),
	0x91: singleWordInst("SUB C"),
	0x92: singleWordInst("SUB D"),
	0x93: singleWordInst("SUB E"),
	0x94: singleWordInst("SUB H"),
	0x95: singleWordInst("SUB L"),
	0x96: singleWordInst("SUB M"),
	0x97: singleWordInst("SUB A"),
	0x98: singleWordInst("SBB B"),
	0x99: singleWordInst("SBB C"),
	0x9a: singleWordInst("SBB D"),
	0x9b: singleWordInst("SBB E"),
	0x9c: singleWordInst("SBB H"),
	0x9d: singleWordInst("SBB L"),
	0x9e: singleWordInst("SBB M"),
	0x9f: singleWordInst("SBB A"),
	0xa0: singleWordInst("ANA B"),
	0xa1: singleWordInst("ANA C"),
	0xa2: singleWordInst("ANA D"),
	0xa3: singleWordInst("ANA E"),
	0xa4: singleWordInst("ANA H"),
	0xa5: singleWordInst("ANA L"),
	0xa6: singleWordInst("ANA M"),
	0xa7: singleWordInst("ANA A"),
	0xa8: singleWordInst("XRA B"),
	0xa9: singleWordInst("XRA C"),
	0xaa: singleWordInst("XRA D"),
	0xab: singleWordInst("XRA E"),
	0xac: singleWordInst("XRA H"),
	0xad: singleWordInst("XRA L"),
	0xae: singleWordInst("XRA M"),
	0xaf: singleWordInst("XRA A"),
	0xb0: singleWordInst("ORA B"),
	0xb1: singleWordInst("ORA C"),
	0xb2: singleWordInst("ORA D"),
	0xb3: singleWordInst("ORA E"),
	0xb4: singleWordInst("ORA H"),
	0xb5: singleWordInst("ORA L"),
	0xb6: singleWordInst("ORA M"),
	0xb7: singleWordInst("ORA A"),
	0xb8: singleWordInst("CMP B"),
	0xb9: singleWordInst("CMP C"),
	0xba: singleWordInst("CMP D"),
	0xbb: singleWordInst("CMP E"),
	0xbc: singleWordInst("CMP H"),
	0xbd: singleWordInst("CMP L"),
	0xbe: singleWordInst("CMP M"),
	0xbf: singleWordInst("CMP A"),
	0xc0: singleWordInst("RNZ"),
	0xc1: singleWordInst("POP B"),
	0xc2: tripleWordInst("JNZ adr"),
	0xc3: tripleWordInst("JMP adr"),
	0xc4: tripleWordInst("CNZ adr"),
	0xc5: singleWordInst("PUSH B"),
	0xc6: doubleWordInst("ADI D8"),
	0xc7: singleWordInst("RST 0"),
	0xc8: singleWordInst("RZ"),
	0xc9: singleWordInst("RET"),
	0xca: tripleWordInst("JZ adr"),
	0xcc: tripleWordInst("CZ adr"),
	0xcd: tripleWordInst("CALL adr"),
	0xce: doubleWordInst("ACI D8"),
	0xcf: singleWordInst("RST 1"),
	0xd0: singleWordInst("RNC"),
	0xd1: singleWordInst("POP D"),
	0xd2: tripleWordInst("JNC adr"),
	0xd3: doubleWordInst("OUT D8"),
	0xd4: tripleWordInst("CNC adr"),
	0xd5: singleWordInst("PUSH D"),
	0xd6: doubleWordInst("SUI D8"),
	0xd7: singleWordInst("RST 2"),
	0xd8: singleWordInst("RC"),
	0xda: tripleWordInst("JC adr"),
	0xdb: doubleWordInst("IN D8"),
	0xdc: tripleWordInst("CC adr"),
	0xde: doubleWordInst("SBI D8"),
	0xdf: singleWordInst("RST 3"),
	0xe0: singleWordInst("RPO"),
	0xe1: singleWordInst("POP H"),
	0xe2: tripleWordInst("JPO adr"),
	0xe3: singleWordInst("XTHL"),
	0xe4: tripleWordInst("CPO adr"),
	0xe5: singleWordInst("PUSH H"),
	0xe6: doubleWordInst("ANI D8"),
	0xe7: singleWordInst("RST 4"),
	0xe8: singleWordInst("RPE"),
	0xe9: singleWordInst("PCHL"),
	0xea: tripleWordInst("JPE adr"),
	0xeb: singleWordInst("XCHG"),
	0xec: tripleWordInst("CPE adr"),
	0xee: doubleWordInst("XRI D8"),
	0xef: singleWordInst("RST 5"),
	0xf0: singleWordInst("RP"),
	0xf1: singleWordInst("POP PSW"),
	0xf2: tripleWordInst("JP adr"),
	0xf3: singleWordInst("DI"),
	0xf4: tripleWordInst("CP adr"),
	0xf5: singleWordInst("PUSH PSW"),
	0xf6: doubleWordInst("ORI D8"),
	0xf7: singleWordInst("RST 6"),
	0xf8: singleWordInst("RM"),
	0xf9: singleWordInst("SPHL"),
	0xfa: tripleWordInst("JM adr"),
	0xfb: singleWordInst("EI"),
	0xfc: tripleWordInst("CM adr"),
	0xfe: doubleWordInst("CPI D8"),
	0xff: singleWordInst("RST 7"),
}

// Disassemble reads machine code from the reader, and writes
// assembly code to the writer
func Disassemble(r io.Reader, w io.Writer) error {
	return DisassembleFrom(r, w, 0)
}

// DisassembleFrom reads machine code from the reader starting
// at the given offset, and writes assembly code to the writer
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
			return fmt.Errorf("unkown op code %x", op)
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

func singleWordInst(inst string) instDasm {
	return func(_ *byteReader, w *bufio.Writer) error {
		_, err := w.WriteString(inst + "\n")
		return err
	}
}

func doubleWordInst(inst string) instDasm {
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

func tripleWordInst(inst string) instDasm {
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
