// emu implements the 8080 emulator
package emu

import (
	"fmt"
	"math"
	"strings"
)

const (
	kilobyte = 1 << 10
	// MemSize is the whole amount of Memory in the computer
	MemSize = 16 * kilobyte
	// RomSize is the size of the ROM area
	RomSize = 8 * kilobyte
)

// ComputerError denotes an error condition in the computer
type ComputerError string

// Error implements the error interface
func (e ComputerError) Error() string {
	return string(e)
}

// Flags hold the settings of the five condition bits, i.e., Carry, Zero, Sign, Parity, and Auxiliary Carry.
// The format of this byte is, according to the 8080 assembly programming manual:
//
//	| | | |A| | | | |
//	|S|Z|0|C|0|P|1|C|
//
// S  State of Sign bit: the number if negative if interpreted as signed (most significant bit of the operation result is 1)
// Z  State of Zero bit: the result of an operation was zero (all bits of the operation result are zero)
// 0  always 0 (ignored)
// AC State of auxiliary carry bit: like carry but applied to the least significant 4 bits of the number.
// 0  always 0 (ignored)
// P  State of Parity bit: the amount of ones in the result of an operation is even
// 1  always 1 (ignored)
// C  State of carry bit: an addition exceeded the maximum number, or there was borrow in a substraction.
type Flags byte

const (
	none Flags = iota
	cf   Flags = 1 << (iota - 1)
	_
	pf
	_
	acf
	_
	zf
	sf
)

// String returns a string containing a code for each of the flags set
func (f Flags) String() string {
	var flags []string
	if f.zero() {
		flags = append(flags, "Z")
	}
	if f.sign() {
		flags = append(flags, "S")
	}
	if f.parity() {
		flags = append(flags, "P")
	}
	if f.carry() {
		flags = append(flags, "CY")
	}
	if f.auxiliaryCarry() {
		flags = append(flags, "AC")
	}
	return strings.Join(flags, " ")
}

// carry returns whether the carry flag is set
func (f Flags) carry() bool {
	return (f & cf) != 0
}

// auxiliaryCarry returns whether the auxiliary carry flag is set
func (f Flags) auxiliaryCarry() bool {
	return (f & acf) != 0
}

// parity returns whether the parity flag is set
func (f Flags) parity() bool {
	return (f & pf) != 0
}

// sign returns whether the sign flag is set
func (f Flags) sign() bool {
	return (f & sf) != 0
}

// zero returns whether the zero flag is set
func (f Flags) zero() bool {
	return (f & zf) != 0
}

// parity8 calculates the parity of the given byte, and returns a Flags value with the parity8 flag set appropriately
func parity8(b byte) Flags {
	i := b ^ (b >> 1)
	i = i ^ (i >> 2)
	i = i ^ (i >> 4)
	if i&1 == 0 {
		return pf
	}
	return 0
}

// sign8 calculates the sign of the given byte, and returns a Flags value with the sign flag set appropriately
func sign8(b byte) Flags {
	if b&0x80 == 0x80 {
		return sf
	}
	return none
}

// zero8 calculates the zero flag of the given byte, and returns a Flags value with the sign flag set appropriately
func zero8(b byte) Flags {
	if b == 0x0 {
		return zf
	}
	return none
}

// CPU is the central processing unit comprised of the registers and arithmetic-logic unit (ALU).
//
// For simplicity, we inline the structs of the alu and the registers bank in this struct.
//
// The 8080 processor has 8 registers in its registry bank:
//  * six 8-bit registers  (B-L). -there's another general purpose registry (A), but the hardware for it belongs to the
// ALU, more on that below-
//  * two 16-bit registers: the stack pointer (SP) and program counter (PC)
//
// the 8 bit registers come in pairs (B-C, D-E, H-L) and some opcodes operate on the pair itself, for instance LXI B, D16.
//
// The alu (arithmetic-logic unit) contains 5 Flags (zero, sign, parity, carry, and auxiliary carry), and special
// registers that belong to the ALU and not the register array: The accumulator registry (A) is used to store the result
// of several arithmetic operations. While logically most opcodes treat the A registry as a general purpose one, this
// resides in the ALU.
type CPU struct {
	A byte
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte

	SP uint16
	PC uint16

	Flags Flags
}

// BC returns the register-pair B-C
func (c *CPU) BC() uint16 {
	return uint16(c.B)<<8 + uint16(c.C)
}

// DE returns the register-pair D-E
func (c *CPU) DE() uint16 {
	return uint16(c.D)<<8 + uint16(c.E)
}

// HL returns the register-pair H-L
func (c *CPU) HL() uint16 {
	return uint16(c.H)<<8 + uint16(c.L)
}

// Computer connects the Memory and the cpu
type Computer struct {
	CPU
	Mem []byte
}

// newComputer creates a new computer with the cpu and memory states given
func newComputer(c CPU, m []byte) *Computer {
	return &Computer{
		CPU: c,
		Mem: m,
	}
}

// Load loads the ROM into a newly created computer main Memory
func Load(rom []byte) *Computer {
	c := newComputer(CPU{}, make([]byte, MemSize))
	copy(c.Mem[:RomSize], rom)
	return c
}

// snapshot creates a copy of the current state of the computer
func (c *Computer) snapshot() *Computer {
	return newComputer(c.CPU, c.Mem)
}

func (c *Computer) String() string {
	template := `
╔═══════════════════════════════════════════════════════╗
║                          CPU                          ║
╠═══════════════════════════════════════════════════════╣
║ A    ┆ $REGA              ║  B-C  ┆  $REGB            ║
║ D-E  ┆ $REGD              ║  H-L  ┆  $REGH            ║
║ SP   ┆ $REGS              ║  PC   ┆  $REGP            ║
╟───────────────────────────────────────────────────────╢
║                           ║ Flags ┆  $FLAG_VALUES     ║
╠═══════════════════════════════════════════════════════╣
║                        Memory                         ║  
╠═══════════════════════════════════════════════════════╣
║ 0000: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ║
╚═══════════════════════════════════════════════════════╝
`

	s := strings.Replace(template, "$REGA", fmt.Sprintf("%02X   ", c.A), 1)
	s = strings.Replace(s, "$REGB", fmt.Sprintf("%04X ", c.BC()), 1)
	s = strings.Replace(s, "$REGD", fmt.Sprintf("%04X ", c.DE()), 1)
	s = strings.Replace(s, "$REGH", fmt.Sprintf("%04X ", c.HL()), 1)
	s = strings.Replace(s, "$REGS", fmt.Sprintf("%04X ", c.SP), 1)
	s = strings.Replace(s, "$REGP", fmt.Sprintf("%04X ", c.PC), 1)
	s = strings.Replace(s, " $FLAG_VALUES  ", fmt.Sprintf("%-15s", c.Flags.String()), 1)

	memory := strings.Builder{}
	buf := make([]byte, 0x10)
	rows := int(math.Ceil(float64(len(c.Mem)) / 0x10))
	for r := 0; r < rows; r++ {
		memory.WriteString("║ ")
		copy(buf, c.Mem[r*0x10:])
		memory.WriteString(fmt.Sprintf("%03X0: ", r))
		for pos := range buf {
			memory.WriteString(fmt.Sprintf("%02x ", buf[pos]))
		}
		memory.WriteString("║\n")
	}

	s = strings.Replace(s, "║ 0000: 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 00 ║\n", memory.String(), 1)
	return s
}

// Step executes one instruction of the code pointed by the Program Counter (PC) of the CPU
func (c *Computer) Step(df DebugFilter) error {
	op, err := c.read8(c.PC)
	if err != nil {
		return err
	}
	if int(op) >= len(it) || it[op] == nil {
		return fmt.Errorf("unimplemented op %02X", op)
	}

	if df != nil && df(op) {
		prev := c.snapshot()
		defer c.debug(prev)
	}

	err = it[op](c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Computer) read16(addr uint16) (uint16, error) {
	l, err := c.read8(addr)
	if err != nil {
		return 0, err
	}

	h, err := c.read8(addr + 1)
	if err != nil {
		return 0, err
	}

	return uint16(h)<<8 + uint16(l), nil
}

func (c *Computer) read8(addr uint16) (byte, error) {
	if int(addr) > len(c.Mem) {
		return 0, ComputerError(fmt.Sprintf("segfault accessing %04X", addr))
	}
	return c.Mem[addr], nil
}

func (c *Computer) write8(addr uint16, d8 byte) error {
	if int(addr) > len(c.Mem) {
		return ComputerError(fmt.Sprintf("segfault accessing %04X", addr))
	}
	c.Mem[addr] = d8
	return nil
}

func (c *Computer) read8Indirect() (byte, error) {
	return c.read8(c.HL())
}

func (c *Computer) write8Indirect(v byte) error {
	return c.write8(c.HL(), v)
}
