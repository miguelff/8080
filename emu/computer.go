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
// carry, and auxiliary carry), an accumulator register (A),
// a temporary register (TMP) and a temporary accumulator
// register (TACC).
type alu struct {
	Z  bool
	S  bool
	P  bool
	CY bool
	AC bool

	A    byte
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
	0x01: lxib,
	0x03: inxb,
	0x06: mvib,
	0x11: lxid,
	0x13: inxd,
	0x1A: ldaxd,
	0x21: lxih,
	0x23: inxh,
	0x31: lxisp,
	0x33: inxsp,
	0x70: movmb,
	0x71: movmc,
	0x72: movmd,
	0x73: movme,
	0x74: movmh,
	0x75: movml,
	0x77: movma,
	0xC3: jmp,
	0xCD: call,
}

// 0x00: NOP
// Move to the next instruction
func nop(c *Computer) error {
	c.PC++
	return nil
}

// 0x01: LXI B | D16. B <- byte 3, C <- byte 2
// Loads double word in registers B and C.
func lxib(c *Computer) error {
	return lxi(c, &c.C, &c.B)
}

// 0x03: INX B | B <- B + 1
// Increments B. No condition flags are affected
func inxb(c *Computer) error {
	return inx(c, &c.B, &c.C)
}

// 0x06: MVI B | D8. B <- byte 2
// Loads word into B register
func mvib(c *Computer) error {
	return mvi(c, &c.B)
}

// 0x11: LXI D | D16. D <- byte 3, E <- byte 2
// Loads double word in registers D and E.
func lxid(c *Computer) error {
	return lxi(c, &c.E, &c.D)
}

// 0x13: INX D | D <- D + 1
// Increments D. No condition flags are affected
func inxd(c *Computer) error {
	return inx(c, &c.D, &c.E)
}

// 0x1A: LDAX D | A <- (DE)
// Loads into the Accumulator record the value pointed
// by the address denoted by the DE register group.
func ldaxd(c *Computer) error {
	addr := uint16(c.D)<<8 + uint16(c.E)
	b, err := c.readD8(addr)
	if err != nil {
		return err
	}
	c.A = b
	c.PC++
	return nil
}

// 0x21: LXI H, D161 | H <- byte 3, L <- byte 2
// Loads double word in the register pair HL
func lxih(c *Computer) error {
	return lxi(c, &c.L, &c.H)
}

// 0x23: INX H | H <- H + 1
// Increments H. No condition flags are affected
func inxh(c *Computer) error {
	return inx(c, &c.L, &c.H)
}

// 0x31: LXI SP, D16 | SP.hi <- byte 3, SP.lo <- byte 2
// Resets the stack pointer to a given value
func lxisp(c *Computer) error {
	return lxiD16(c, &c.SP)
}

// 0x33: INX SP | SP <- SP + 1
// Increments SP. No condition flags are affected
func inxsp(c *Computer) error {
	return inxD16(c, &c.SP)
}

// 0x77: MOV M,A. | (HL) <- B
// Writes B to the address pointed by the register pair HL.
func movmb(c *Computer) error {
	return movm(c, c.B)
}

// 0x77: MOV M,A. | (HL) <- C
// Writes C to the address pointed by the register pair HL.
func movmc(c *Computer) error {
	return movm(c, c.C)
}

// 0x77: MOV M,A. | (HL) <- D
// Writes D to the address pointed by the register pair HL.
func movmd(c *Computer) error {
	return movm(c, c.D)
}

// 0x77: MOV M,A. | (HL) <- E
// Writes E to the address pointed by the register pair HL.
func movme(c *Computer) error {
	return movm(c, c.E)
}

// 0x77: MOV M,A. | (HL) <- H
// Writes H to the address pointed by the register pair HL.
func movmh(c *Computer) error {
	return movm(c, c.H)
}

// 0x77: MOV M,A. | (HL) <- L
// Writes L to the address pointed by the register pair HL.
func movml(c *Computer) error {
	return movm(c, c.L)
}

// 0x77: MOV M,A. | (HL) <- A
// Writes A to the address pointed by the register pair HL.
func movma(c *Computer) error {
	return movm(c, c.A)
}

// 0xC3: JMP adr | PC <- adr.
// Jump to the address denoted by the next two bytes.
func jmp(c *Computer) error {
	return lxiD16(c, &c.PC)
}

// 0xCD: CALL adr | (SP-1)<-PC.hi;(SP-2)<-PC.lo;SP<-SP-2;PC=adr
// CALL pushes the program counter (PC) into the stack (SP), and
// updates the program counter to point to adr.
func call(c *Computer) error {
	err := pushD16(c, c.PC)
	if err != nil {
		return err
	}

	err = jmp(c)
	if err != nil {
		return err
	}

	return nil
}

func mvi(c *Computer, register *byte) error {
	w, err := c.readD8(c.PC + 1)
	if err != nil {
		return err
	}

	c.PC += 2
	*register = w
	return nil
}

func inx(c *Computer, lsbRegister *byte, msbRegister *byte) error {
	incr := (uint16(*msbRegister)<<8 + uint16(*lsbRegister)) + 1
	c.PC++
	*msbRegister = byte((incr >> 8) & 0x00ff)
	*lsbRegister = byte(incr & 0x00ff)
	return nil
}

func inxD16(c *Computer, register *uint16) error {
	*register++
	c.PC++
	return nil
}

func lxi(c *Computer, lsbRegister *byte, msbRegister *byte) error {
	lsb, err := c.readD8(c.PC + 1)
	if err != nil {
		return err
	}

	msb, err := c.readD8(c.PC + 2)
	if err != nil {
		return err
	}

	c.PC += 3
	*lsbRegister = lsb
	*msbRegister = msb
	return nil
}

func lxiD16(c *Computer, register *uint16) error {
	dw, err := c.readD16(c.PC + 1)

	if err != nil {
		return err
	}

	c.PC += 3
	*register = dw
	return nil
}

func movm(c *Computer, r byte) error {
	addr := uint16(c.H)<<8 + uint16(c.L)
	err := c.writeD8(addr, r)
	if err != nil {
		return err
	}
	c.PC++
	return nil
}

func pushD16(c *Computer, d16 uint16) error {
	msb := byte(d16 & 0x00FF)
	lsb := byte((d16 & 0xFF00) >> 8)

	err := c.writeD8(c.SP-1, msb)
	if err != nil {
		return err
	}

	err = c.writeD8(c.SP-2, lsb)
	if err != nil {
		return err
	}
	c.SP -= 2
	return nil
}
