package emu

import (
	"encoding/hex"
	"fmt"
	"strings"

	"github.com/miguelff/8080/dasm"
)

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

func (c *Computer) diff(other *Computer) diff {
	var d diff
	if c.A != other.A {
		d = append(d, makeChange("A", fmt.Sprintf("%02X", c.A), fmt.Sprintf("%02X", other.A)))
	}
	if c.B != other.B {
		d = append(d, makeChange("B", fmt.Sprintf("%02X", c.B), fmt.Sprintf("%02X", other.B)))
	}
	if c.C != other.C {
		d = append(d, makeChange("C", fmt.Sprintf("%02X", c.C), fmt.Sprintf("%02X", other.C)))
	}
	if c.D != other.D {
		d = append(d, makeChange("D", fmt.Sprintf("%02X", c.D), fmt.Sprintf("%02X", other.D)))
	}
	if c.E != other.E {
		d = append(d, makeChange("E", fmt.Sprintf("%02X", c.E), fmt.Sprintf("%02X", other.E)))
	}
	if c.H != other.H {
		d = append(d, makeChange("H", fmt.Sprintf("%02X", c.H), fmt.Sprintf("%02X", other.H)))
	}
	if c.L != other.L {
		d = append(d, makeChange("L", fmt.Sprintf("%02X", c.L), fmt.Sprintf("%02X", other.L)))
	}
	if c.SP != other.SP {
		d = append(d, makeChange("SP", fmt.Sprintf("%04X", c.SP), fmt.Sprintf("%04X", other.SP)))
	}
	if c.PC != other.PC {
		d = append(d, makeChange("PC", fmt.Sprintf("%04X", c.PC), fmt.Sprintf("%04X", other.PC)))
	}
	if c.Flags != other.Flags {
		d = append(d, makeChange("Flags", c.Flags.String(), other.Flags.String()))
	}
	return d
}

type diff []change

func (d diff) String() string {
	if len(d) == 0 {
		return "No changes"
	}

	sb := strings.Builder{}
	for i := range d {
		sb.WriteString(fmt.Sprintf("\t- %s\n", d[i].String()))
	}
	return sb.String()
}

type change struct {
	name string
	old  string
	new  string
}

func makeChange(name, old, new string) change {
	return change{name, old, new}
}

func (c *change) String() string {
	return fmt.Sprintf("%s changed: %s → %s", c.name, c.old, c.new)
}
