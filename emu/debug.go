package emu

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/miguelff/8080/encoding"

	"github.com/miguelff/8080/dasm"
)

// DebugFilter is a predicate indicating wether or not
// to issue a debug trace for the give opcode
type DebugFilter func(opcode byte) bool

// DebugAll debugs all opcodes
func DebugAll(_ byte) bool { return true }

// DebugNone doesn't debug any opcode
func DebugNone(_ byte) bool { return false }

// MakeDebugFilter creates a DebugFilter that will select the
// opcodes denoted by the given string.
//
// MakeDebugFilter("all") will debug all symbols
// MakeDebugFilter("C9 CD") will debug CALL and RET instructions
func MakeDebugFilter(def string) DebugFilter {
	if def == "all" {
		return DebugAll
	} else {
		opcodes := encoding.HexToBin(def)

		return func(opcode byte) bool {
			for i := range opcodes {
				if opcode == opcodes[i] {
					return true
				}
			}
			return false
		}
	}
}

func (c *Computer) debug(prev *Computer) {
	context := make([]byte, 4)
	copy(context, prev.Mem[prev.PC:])
	assembly, err := dasm.DisassembleFirst(context)

	if err != nil {
		fmt.Printf("Error dissassembing bytes: %v\n", err)
	} else {
		fmt.Printf("(%s) PC: %04X  Mem(%04X-%04X): %s | %s \n", prev.Flags, prev.PC, prev.PC, prev.PC+4, hex.EncodeToString(context), strings.TrimSpace(assembly))
	}
	fmt.Println(prev.diff(c))
}

func (c *Computer) diff(other *Computer) string {
	var sb strings.Builder
	if c.A != other.A {
		sb.WriteString(fmt.Sprintf("- A: %02X → %02X\n", c.A, other.A))
	}
	if c.B != other.B {
		sb.WriteString(fmt.Sprintf("- B: %02X → %02X\n", c.B, other.B))
	}
	if c.C != other.C {
		sb.WriteString(fmt.Sprintf("- C: %02X → %02X\n", c.C, other.C))
	}
	if c.D != other.D {
		sb.WriteString(fmt.Sprintf("- D: %02X → %02X\n", c.D, other.D))
	}
	if c.E != other.E {
		sb.WriteString(fmt.Sprintf("- E: %02X → %02X\n", c.E, other.E))
	}
	if c.H != other.H {
		sb.WriteString(fmt.Sprintf("- H: %02X → %02X\n", c.H, other.H))
	}
	if c.L != other.L {
		sb.WriteString(fmt.Sprintf("- E: %02X → %02X\n", c.E, other.E))
	}
	if c.SP != other.SP {
		sb.WriteString(fmt.Sprintf("- SP: %04X → %04X\n", c.SP, other.SP))
	}
	if c.PC != other.PC {
		sb.WriteString(fmt.Sprintf("- PC: %04X → %04X\n", c.PC, other.PC))
	}
	if c.Flags != other.Flags {
		sb.WriteString(fmt.Sprintf("- Flags: %s → %s\n", c.Flags, other.Flags))
	}
	return sb.String()
}
