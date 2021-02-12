// emu implements the 8080 emulator
package emu

import (
	"fmt"
)

// ComputerError denotes an error condition in the computer
type ComputerError string

func (e ComputerError) Error() string {
	return string(e)
}

const (
	kilobyte = 1 << 10
	MemSize  = 16 * kilobyte
	RomSize  = 8 * kilobyte
)

// registerArray contains 8 registers: 6 8-bit registers
// (B-L); and two 16-bit registers: the stack pointer (SP)
// and program counter (PC)
//
// 8 bit registers come in pairs (B-C, D-E, H-L) and some opcodes
// operate on the pair itself, for instance LXI B, D16 loads two bytes
// in registers B (most significant byte) and C (least significant byte)
type registerArray struct {
	B  byte
	C  byte
	D  byte
	E  byte
	H  byte
	L  byte
	SP uint16
	PC uint16
}

// alu (arithmetic-logic unit) contains 5 flags (zero, sign, parity,
// carry, and auxiliary carry), an accumulator register (ACC),
// a temporary register (TMP) and a temporary accumulator
// register (TACC).
type alu struct {
	Z  bool
	S  bool
	P  bool
	CY bool
	AC bool

	ACC  byte
	TMP  byte
	TACC byte
}

// cpu is the central processing unit comprised of the
// registers and alu
type cpu struct {
	registerArray
	alu
}

// Memory represents the computer memory
type memory []byte

// Computer connects the memory and the cpu
type Computer struct {
	cpu
	mem memory
}

// Load loads the ROM into the computer main memory
func (c *Computer) Load(rom []byte) {
	c.mem = make(memory, MemSize)
	copy(c.mem[:RomSize], rom)
}

// Step executes one instruction of the code pointed by the Program Counter (PC) of the CPU
func (c *Computer) Step() error {
	opcode, err := c.readD8(c.PC)
	if err != nil {
		return err
	}

	if int(opcode) > len(instructionTable) || instructionTable[opcode] == nil {
		return fmt.Errorf("unimplemented opcode %02X", opcode)
	}

	err = instructionTable[opcode](c)
	if err != nil {
		return err
	}
	return nil
}

func (c *Computer) readD16(addr uint16) (uint16, error) {
	l, err := c.readD8(addr)
	if err != nil {
		return 0, err
	}

	h, err := c.readD8(addr + 1)
	if err != nil {
		return 0, err
	}

	return uint16(h)<<8 + uint16(l), nil
}

func (c *Computer) readD8(addr uint16) (byte, error) {
	if int(addr) > len(c.mem) {
		return 0, ComputerError(fmt.Sprintf("segfault accessing %04X", addr))
	}
	return c.mem[addr], nil
}

func (c *Computer) writeD8(addr uint16, d8 byte) error {
	if int(addr) > len(c.mem) {
		return ComputerError(fmt.Sprintf("segfault accessing %04X", addr))
	}
	c.mem[addr] = d8
	return nil
}

type instruction func(*Computer) error

var instructionTable = []instruction{
	0x00: nop,
	0x06: mvib,
	0x31: lxisp,
	0xC3: jmp,
	0xCD: call,
}

// 0x00: NOP. Move to the next instruction
func nop(c *Computer) error {
	c.PC++
	return nil
}

// 0x06: MVI B, D8. B <- byte 2
// Loads word into B register
func mvib(c *Computer) error {
	return loadD8Register(c, &c.B)
}

// 0x31: LXI SP, D16 | SP.hi <- byte 3, SP.lo <- byte 2
// Reset the stack pointer to a given value
func lxisp(c *Computer) error {
	return loadD16Register(c, &c.SP)
}

// 0xC3: JMP adr | PC <= adr.
// Jump to the address denoted by the next two bytes.
func jmp(c *Computer) error {
	return loadD16Register(c, &c.PC)
}

// 0xCD: CALL adr | (SP-1)<-PC.hi;(SP-2)<-PC.lo;SP<-SP-2;PC=adr
// CALL pushes the program counter (PC) into the stack (SP), and
// updates the program counter to point to adr.
func call(c *Computer) error {
	err := pushD16(c, c.PC)
	if err != nil {
		return err
	}

	jmp(c)
	return nil
}

func pushD16(c *Computer, d16 uint16) error {
	hi := byte(d16 & 0x00FF)
	lo := byte(d16 >> 8)

	err := c.writeD8(c.SP-1, hi)
	if err != nil {
		return err
	}

	err = c.writeD8(c.SP-2, lo)
	if err != nil {
		return err
	}
	c.SP -= 2
	return nil
}

func loadD8Register(c *Computer, register *byte) error {
	c.PC++
	w, err := c.readD8(c.PC)
	c.PC++
	if err != nil {
		return err
	}

	*register = w
	return nil
}

func loadD16Register(c *Computer, register *uint16) error {
	c.PC++
	dw, err := c.readD16(c.PC)
	c.PC += 2
	if err != nil {
		return err
	}

	*register = dw
	return nil
}
