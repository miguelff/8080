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
	kilobyte               = 1 << 10
	MemSize                = 16 * kilobyte
	RomSize                = 8 * kilobyte
	ErrEOM   ComputerError = "reached end of memory"
)

// registerArray contains 8 registers: 6 single-word registers
// (B-L); and two double-word registers: the stack pointer (SP)
// and program counter (PC)
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

// alu (arithmetic-logic unit) contains 5 flags (zero, carry,
// sign, parity, and auxiliary carry); an accumulator register
// (ACC) a temporary register (TMP) and a temporary accumulator
// register (TACC)
type alu struct {
	ZF bool
	CF bool
	SF bool
	PF bool
	AF bool

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
	var err error

	opcode, err := c.currentByte()
	if err != nil {
		return err
	}

	if instruction := instructionTable[opcode]; instruction != nil {
		err = instruction(c)
	} else {
		err = fmt.Errorf("unimplemented opcode %02X", opcode)
	}
	c.PC++

	return err
}

func (c *Computer) nextByte() (byte, error) {
	c.PC++
	return c.currentByte()
}

func (c *Computer) currentByte() (byte, error) {
	if int(c.PC) > len(c.mem) {
		return 0, ErrEOM
	}
	return c.mem[c.PC], nil
}

type instruction func(*Computer) error

var instructionTable = []instruction{
	0x00: nop,
	0x01: lxiB,
}

// nop (do nothing)
func nop(_ *Computer) error {
	return nil
}

// LXI B, D16. Move to register pair B (registers B, C), the 16 bits
// denoted by the following 2 bytes (in little endian form)
func lxiB(c *Computer) error {
	b, err := c.nextByte()
	if err != nil {
		return err
	}
	c.C = b

	b, err = c.nextByte()
	if err != nil {
		return err
	}
	c.B = b
	return nil
}
