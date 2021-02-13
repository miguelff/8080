// emu implements the 8080 emulator
package emu

import (
	"fmt"
)

const (
	kilobyte = 1 << 10
	// MemSize is the whole amount of memory in the computer
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

// registers contains 8 registers: 6 8-bit registers  (B-L); and two 16-bit registers: the stack pointer (SP) and
// program counter (PC)
//
// 8 bit registers come in pairs (B-C, D-E, H-L) and some opcodes operate on the pair itself, for instance LXI B, D16
// loads two bytes in registers B (most significant byte) and C (least significant byte)
type registers struct {
	B byte
	C byte
	D byte
	E byte
	H byte
	L byte

	SP uint16
	PC uint16
}

// flags encapsulate
type flags byte

const (
	none = iota
	// Zero. Toggled when arithmetic operation was 0.
	z flags = 1 << iota
	// Sign. Toggled when arithmetic operation results in a negative number (i.e. Its most significant bit active)
	s
	// Parity
	p
	// Carry
	cy
	// Auxiliary Carry
	ac
)

// alu (arithmetic-logic unit) contains 5 flags (zero, sign, parity, carry, and auxiliary carry), and special registers
// that belong to the ALU and not the register array: Register (A), a temporary register (TMP) and a temporary
// accumulator register (TACC).
type alu struct {
	flags flags

	A    byte
	TMP  byte
	TACC byte
}

func (a *alu) Z() bool {
	return (a.flags & z) != 0
}

func (a *alu) S() bool {
	return (a.flags & s) != 0
}

func (a *alu) P() bool {
	return (a.flags & p) != 0
}

func (a *alu) CY() bool {
	return (a.flags & cy) != 0
}

func (a *alu) AC() bool {
	return (a.flags & ac) != 0
}

// parity8 calculates the parity of the given byte, and returns a flags value with the parity flag set appropriately
func parity8(result byte) flags {
	i := result ^ (result >> 1)
	i = i ^ (i >> 2)
	i = i ^ (i >> 4)
	if i&1 == 0 {
		return p
	}
	return 0
}

// parity16 calculates the parity of the given 16 bits, and returns a flags value with the parity flag set appropriately
func parity16(result uint16) flags {
	i := result ^ (result >> 1)
	i = i ^ (i >> 2)
	i = i ^ (i >> 4)
	i = i ^ (i >> 8)
	if i&1 == 0 {
		return p
	}
	return 0
}

// cpu is the central processing unit comprised of the  registers and alu
type cpu struct {
	registers
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
	opcode, err := c.read8(c.PC)
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
	if int(addr) > len(c.mem) {
		return 0, ComputerError(fmt.Sprintf("segfault accessing %04X", addr))
	}
	return c.mem[addr], nil
}

func (c *Computer) write8(addr uint16, d8 byte) error {
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
	0x0E: mvic,
	0x11: lxid,
	0x13: inxd,
	0x16: mvid,
	0x1E: mvie,
	0x1A: ldaxd,
	0x21: lxih,
	0x23: inxh,
	0x26: mvih,
	0x2E: mvil,
	0x31: lxisp,
	0x33: inxsp,
	0x3E: mvia,
	0x40: movbb,
	0x41: movbc,
	0x42: movbd,
	0x43: movbe,
	0x44: movbh,
	0x45: movbl,
	0x47: movba,
	0x48: movcb,
	0x49: movcc,
	0x4A: movcd,
	0x4B: movce,
	0x4C: movch,
	0x4D: movcl,
	0x4F: movca,
	0x50: movdb,
	0x51: movdc,
	0x52: movdd,
	0x53: movde,
	0x54: movdh,
	0x55: movdl,
	0x57: movda,
	0x58: moveb,
	0x59: movec,
	0x5A: moved,
	0x5B: movee,
	0x5C: moveh,
	0x5D: movel,
	0x5F: movea,
	0x60: movhb,
	0x61: movhc,
	0x62: movhd,
	0x63: movhe,
	0x64: movhh,
	0x65: movhl,
	0x67: movha,
	0x68: movlb,
	0x69: movlc,
	0x6A: movld,
	0x6B: movle,
	0x6C: movlh,
	0x6D: movll,
	0x6F: movla,
	0x70: movmb,
	0x71: movmc,
	0x72: movmd,
	0x73: movme,
	0x74: movmh,
	0x75: movml,
	0x77: movma,
	0x78: movab,
	0x79: movac,
	0x7A: movad,
	0x7B: movae,
	0x7C: movah,
	0x7D: moval,
	0x7F: movaa,
	0xC3: jmp,
	0xCD: call,
}

// 0xCD: CALL adr | (SP-1)<-PC.hi;(SP-2)<-PC.lo;SP<-SP-2;PC=adr
// CALL pushes the program counter (PC) into the stack (SP), and updates the program counter to point to adr.
func call(c *Computer) error {
	err := push16(c, c.PC)
	if err != nil {
		return err
	}

	err = jmp(c)
	if err != nil {
		return err
	}

	return nil
}

func inx(c *Computer, lsreg, msreg *byte) error {
	incr := (uint16(*msreg)<<8 + uint16(*lsreg)) + 1
	*msreg = byte((incr >> 8) & 0x00ff)
	*lsreg = byte(incr & 0x00ff)
	c.PC++
	return nil
}

func inx16(c *Computer, reg *uint16) error {
	*reg++
	if reg != &c.PC {
		c.PC++
	}
	return nil
}

// 0x03: INX B | B <- B + 1
// Increments B. No condition flags are affected
func inxb(c *Computer) error {
	return inx(c, &c.B, &c.C)
}

// 0x13: INX D | D <- D + 1
// Increments D. No condition flags are affected
func inxd(c *Computer) error {
	return inx(c, &c.D, &c.E)
}

// 0x23: INX H | H <- H + 1
// Increments H. No condition flags are affected
func inxh(c *Computer) error {
	return inx(c, &c.L, &c.H)
}

// 0x33: INX SP | SP <- SP + 1
// Increments SP. No condition flags are affected
func inxsp(c *Computer) error {
	return inx16(c, &c.SP)
}

// 0xC3: JMP adr | PC <- adr.
// Jump to the address denoted by the next two bytes.
func jmp(c *Computer) error {
	return lxi16(c, &c.PC)
}

// 0x1A: LDAX D | A <- (DE)
// Loads into the Accumulator record the value pointed by the address denoted by the DE register group.
func ldaxd(c *Computer) error {
	addr := uint16(c.D)<<8 + uint16(c.E)
	b, err := c.read8(addr)
	if err != nil {
		return err
	}

	c.A = b
	c.PC++
	return nil
}

func lxi(c *Computer, lsreg, msreg *byte) error {
	lsb, err := c.read8(c.PC + 1)
	if err != nil {
		return err
	}

	msb, err := c.read8(c.PC + 2)
	if err != nil {
		return err
	}

	*lsreg, *msreg = lsb, msb
	c.PC += 3
	return nil
}

func lxi16(c *Computer, reg *uint16) error {
	dw, err := c.read16(c.PC + 1)
	if err != nil {
		return err
	}

	*reg = dw
	if reg != &c.PC {
		c.PC += 3
	}
	return nil
}

// 0x01: LXI B | D16. B <- byte 3, C <- byte 2
// Loads double word in registers B and C.
func lxib(c *Computer) error {
	return lxi(c, &c.C, &c.B)
}

// 0x11: LXI D | D16. D <- byte 3, E <- byte 2
// Loads double word in registers D and E.
func lxid(c *Computer) error {
	return lxi(c, &c.E, &c.D)
}

// 0x21: LXI H, D161 | H <- byte 3, L <- byte 2
// Loads double word in the register pair HL
func lxih(c *Computer) error {
	return lxi(c, &c.L, &c.H)
}

// 0x31: LXI SP, D16 | SP.hi <- byte 3, SP.lo <- byte 2
// Resets the stack pointer to a given value
func lxisp(c *Computer) error {
	return lxi16(c, &c.SP)
}

func mov(c *Computer, dstreg, srcreg *byte) error {
	*dstreg = *srcreg
	c.PC++
	return nil
}

// 0x7F: MOV A, A | A <- A
// Copies contents from record A to A
func movaa(c *Computer) error {
	return nop(c)
}

// 0x78: MOV A, B | A <- B
// Copies contents from record B to A
func movab(c *Computer) error {
	return mov(c, &c.A, &c.B)
}

// 0x79: MOV A, C | A <- C
// Copies contents from record C to A
func movac(c *Computer) error {
	return mov(c, &c.A, &c.C)
}

// 0x7A: MOV A, D | A <- D
// Copies contents from record D to A
func movad(c *Computer) error {
	return mov(c, &c.A, &c.D)
}

// 0x7B: MOV A, E | A <- E
// Copies contents from record E to A
func movae(c *Computer) error {
	return mov(c, &c.A, &c.E)
}

// 0x7C: MOV A, H | A <- H
// Copies contents from record H to A
func movah(c *Computer) error {
	return mov(c, &c.A, &c.H)
}

// 0x7D: MOV A, L | A <- L
// Copies contents from record L to A
func moval(c *Computer) error {
	return mov(c, &c.A, &c.L)
}

// 0x47: MOV B, A | B <- A
// Copies contents from record A to B
func movba(c *Computer) error {
	return mov(c, &c.B, &c.A)
}

// 0x40: MOV B, B | B <- B
// Copies contents from record B to B
func movbb(c *Computer) error {
	return nop(c)
}

// 0x41: MOV B, C | B <- C
// Copies contents from record C to B
func movbc(c *Computer) error {
	return mov(c, &c.B, &c.C)
}

// 0x42: MOV B, D | B <- D
// Copies contents from record D to B
func movbd(c *Computer) error {
	return mov(c, &c.B, &c.D)
}

// 0x43: MOV B, E | B <- E
// Copies contents from record E to B
func movbe(c *Computer) error {
	return mov(c, &c.B, &c.E)
}

// 0x44: MOV B, H | B <- H
// Copies contents from record H to B
func movbh(c *Computer) error {
	return mov(c, &c.B, &c.H)
}

// 0x45: MOV B, L | B <- L
// Copies contents from record L to B
func movbl(c *Computer) error {
	return mov(c, &c.B, &c.L)
}

// 0x4F: MOV C, A | C <- A
// Copies contents from record A to C
func movca(c *Computer) error {
	return mov(c, &c.C, &c.A)
}

// 0x48: MOV C, B | C <- B
// Copies contents from record B to C
func movcb(c *Computer) error {
	return mov(c, &c.C, &c.B)
}

// 0x49: MOV C, C | C <- C
// Copies contents from record C to C
func movcc(c *Computer) error {
	return nop(c)
}

// 0x4A: MOV C, D | C <- D
// Copies contents from record D to C
func movcd(c *Computer) error {
	return mov(c, &c.C, &c.D)
}

// 0x4B: MOV C, E | C <- E
// Copies contents from record E to C
func movce(c *Computer) error {
	return mov(c, &c.C, &c.E)
}

// 0x4C: MOV C, H | C <- H
// Copies contents from record H to C
func movch(c *Computer) error {
	return mov(c, &c.C, &c.H)
}

// 0x4D: MOV C, L | C <- L
// Copies contents from record L to C
func movcl(c *Computer) error {
	return mov(c, &c.C, &c.L)
}

// 0x57: MOV D, A | D <- A
// Copies contents from record A to D
func movda(c *Computer) error {
	return mov(c, &c.D, &c.A)
}

// 0x50: MOV D, B | D <- B
// Copies contents from record B to D
func movdb(c *Computer) error {
	return mov(c, &c.D, &c.B)
}

// 0x51: MOV D, C | D <- C
// Copies contents from record C to D
func movdc(c *Computer) error {
	return mov(c, &c.D, &c.C)
}

// 0x52: MOV D, D | D <- D
// Copies contents from record D to D
func movdd(c *Computer) error {
	return nop(c)
}

// 0x53: MOV D, E | D <- E
// Copies contents from record E to D
func movde(c *Computer) error {
	return mov(c, &c.D, &c.E)
}

// 0x54: MOV D, H | D <- H
// Copies contents from record H to D
func movdh(c *Computer) error {
	return mov(c, &c.D, &c.H)
}

// 0x55: MOV D, L | D <- L
// Copies contents from record L to D
func movdl(c *Computer) error {
	return mov(c, &c.D, &c.L)
}

// 0x5F: MOV E, A | E <- A
// Copies contents from record A to E
func movea(c *Computer) error {
	return mov(c, &c.E, &c.A)
}

// 0x58: MOV E, B | E <- B
// Copies contents from record B to E
func moveb(c *Computer) error {
	return mov(c, &c.E, &c.B)
}

// 0x59: MOV E, C | E <- C
// Copies contents from record C to E
func movec(c *Computer) error {
	return mov(c, &c.E, &c.C)
}

// 0x5A: MOV E, D | E <- D
// Copies contents from record D to E
func moved(c *Computer) error {
	return mov(c, &c.E, &c.D)
}

// 0x5B: MOV E, E | E <- E
// Copies contents from record E to E
func movee(c *Computer) error {
	return nop(c)
}

// 0x5C: MOV E, H | E <- H
// Copies contents from record H to E
func moveh(c *Computer) error {
	return mov(c, &c.E, &c.H)
}

// 0x5D: MOV E, L | E <- L
// Copies contents from record L to E
func movel(c *Computer) error {
	return mov(c, &c.E, &c.L)
}

// 0x67: MOV H, A | H <- A
// Copies contents from record A to H
func movha(c *Computer) error {
	return mov(c, &c.H, &c.A)
}

// 0x60: MOV H, B | H <- B
// Copies contents from record B to H
func movhb(c *Computer) error {
	return mov(c, &c.H, &c.B)
}

// 0x61: MOV H, C | H <- C
// Copies contents from record C to H
func movhc(c *Computer) error {
	return mov(c, &c.H, &c.C)
}

// 0x62: MOV H, D | H <- D
// Copies contents from record D to H
func movhd(c *Computer) error {
	return mov(c, &c.H, &c.D)
}

// 0x63: MOV H, E | H <- E
// Copies contents from record E to H
func movhe(c *Computer) error {
	return mov(c, &c.H, &c.E)
}

// 0x64: MOV H, H | H <- H
// Copies contents from record H to H
func movhh(c *Computer) error {
	return nop(c)
}

// 0x65: MOV H, L | H <- L
// Copies contents from record L to H
func movhl(c *Computer) error {
	return mov(c, &c.H, &c.L)
}

// 0x6F: MOV L, A | L <- A
// Copies contents from record A to L
func movla(c *Computer) error {
	return mov(c, &c.L, &c.A)
}

// 0x68: MOV L, B | L <- B
// Copies contents from record B to L
func movlb(c *Computer) error {
	return mov(c, &c.L, &c.B)
}

// 0x69: MOV L, C | L <- C
// Copies contents from record C to L
func movlc(c *Computer) error {
	return mov(c, &c.L, &c.C)
}

// 0x6A: MOV L, D | L <- D
// Copies contents from record D to L
func movld(c *Computer) error {
	return mov(c, &c.L, &c.D)
}

// 0x6B: MOV L, E | L <- E
// Copies contents from record E to L
func movle(c *Computer) error {
	return mov(c, &c.L, &c.E)
}

// 0x6C: MOV L, H | L <- H
// Copies contents from record H to L
func movlh(c *Computer) error {
	return mov(c, &c.L, &c.H)
}

// 0x6D: MOV L, L | L <- L
// Copies contents from record L to L
func movll(c *Computer) error {
	return nop(c)
}

func movm(c *Computer, r byte) error {
	addr := uint16(c.H)<<8 + uint16(c.L)
	err := c.write8(addr, r)
	if err != nil {
		return err
	}
	c.PC++
	return nil
}

// 0x77: MOV M,A. | (HL) <- A
// Writes A to the address pointed by the register pair HL.
func movma(c *Computer) error {
	return movm(c, c.A)
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

func mvi(c *Computer, reg *byte) error {
	w, err := c.read8(c.PC + 1)
	if err != nil {
		return err
	}

	*reg = w
	c.PC += 2
	return nil
}

// 0x3E: MVI A, D8 | A <- byte 2
// Loads word into A register
func mvia(c *Computer) error {
	return mvi(c, &c.A)
}

// 0x06: MVI B, D8 | B <- byte 2
// Loads word into B register
func mvib(c *Computer) error {
	return mvi(c, &c.B)
}

// 0x0E: MVI C, D8 | C <- byte 2
// Loads word into C register
func mvic(c *Computer) error {
	return mvi(c, &c.C)
}

// 0x16: MVI D, D8 | D <- byte 2
// Loads word into D register
func mvid(c *Computer) error {
	return mvi(c, &c.D)
}

// 0x1E: MVI E, D8 | E <- byte 2
// Loads word into E register
func mvie(c *Computer) error {
	return mvi(c, &c.E)
}

// 0x26: MVI H, D8 | H <- byte 2
// Loads word into H register
func mvih(c *Computer) error {
	return mvi(c, &c.H)
}

// 0x2E: MVI L, D8 | L <- byte 2
// Loads word into L register
func mvil(c *Computer) error {
	return mvi(c, &c.L)
}

// 0x00: NOP
// Move to the next instruction
func nop(c *Computer) error {
	c.PC++
	return nil
}

func push16(c *Computer, d16 uint16) error {
	msb := byte(d16 & 0x00FF)
	lsb := byte((d16 & 0xFF00) >> 8)

	err := c.write8(c.SP-1, msb)
	if err != nil {
		return err
	}

	err = c.write8(c.SP-2, lsb)
	if err != nil {
		return err
	}
	c.SP -= 2
	return nil
}
