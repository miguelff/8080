package emu

// Instruction is the algorithm that emulates a certain opcode within the computer
type Instruction func(*Computer) error

// it is an opcode to instruction table
var it = []Instruction{
	0x00: nop,
	0x05: dcrb,
	0x01: lxib,
	0x02: staxb,
	0x03: inxb,
	0x04: inrb,
	0x06: mvib,
	0x09: dadb,
	0x0A: ldaxb,
	0x0C: inrc,
	0x0D: dcrc,
	0x0E: mvic,
	0x11: lxid,
	0x12: staxd,
	0x13: inxd,
	0x14: inrd,
	0x15: dcrd,
	0x16: mvid,
	0x19: dadd,
	0x1A: ldaxd,
	0x1C: inre,
	0x1D: dcre,
	0x1E: mvie,
	0x20: dcrh,
	0x21: lxih,
	0x23: inxh,
	0x24: inrh,
	0x26: mvih,
	0x29: dadh,
	0x2C: inrl,
	0x2D: dcrl,
	0x2E: mvil,
	0x31: lxisp,
	0x32: sta,
	0x33: inxsp,
	0x36: mvim,
	0x39: dadsp,
	0x3C: inra,
	0x3D: dcra,
	0x3E: mvia,
	0x40: movbb,
	0x41: movbc,
	0x42: movbd,
	0x43: movbe,
	0x44: movbh,
	0x45: movbl,
	0x46: movfrommb,
	0x47: movba,
	0x48: movcb,
	0x49: movcc,
	0x4A: movcd,
	0x4B: movce,
	0x4C: movch,
	0x4D: movcl,
	0x4E: movfrommc,
	0x4F: movca,
	0x50: movdb,
	0x51: movdc,
	0x52: movdd,
	0x53: movde,
	0x54: movdh,
	0x55: movdl,
	0x56: movfrommd,
	0x57: movda,
	0x58: moveb,
	0x59: movec,
	0x5A: moved,
	0x5B: movee,
	0x5C: moveh,
	0x5D: movel,
	0x5E: movfromme,
	0x5F: movea,
	0x60: movhb,
	0x61: movhc,
	0x62: movhd,
	0x63: movhe,
	0x64: movhh,
	0x65: movhl,
	0x66: movfrommh,
	0x67: movha,
	0x68: movlb,
	0x69: movlc,
	0x6A: movld,
	0x6B: movle,
	0x6C: movlh,
	0x6D: movll,
	0x6E: movfromml,
	0x6F: movla,
	0x70: movtomb,
	0x71: movtomc,
	0x72: movtomd,
	0x73: movtome,
	0x74: movtomh,
	0x75: movtoml,
	0x77: movtoma,
	0x78: movab,
	0x79: movac,
	0x7A: movad,
	0x7B: movae,
	0x7C: movah,
	0x7D: moval,
	0x7E: movfromma,
	0x7F: movaa,
	0x80: addb,
	0x81: addc,
	0x82: addd,
	0x83: adde,
	0x84: addh,
	0x85: addl,
	0x86: addm,
	0x87: adda,
	0x88: adcb,
	0x89: adcc,
	0x8A: adcd,
	0x8B: adce,
	0x8C: adch,
	0x8D: adcl,
	0x8E: adcm,
	0x8F: adca,
	0x90: subb,
	0x91: subc,
	0x92: subd,
	0x93: sube,
	0x94: subh,
	0x95: subl,
	0x97: suba,
	0x98: sbbb,
	0x99: sbbc,
	0x9A: sbbd,
	0x9B: sbbe,
	0x9C: sbbh,
	0x9D: sbbl,
	0x9F: sbba,
	0xA0: anab,
	0xA1: anac,
	0xA2: anad,
	0xA3: anae,
	0xA4: anah,
	0xA5: anal,
	0xA7: anaa,
	0xA8: xraa,
	0xA9: xrab,
	0xAA: xrac,
	0xAB: xrad,
	0xAC: xrae,
	0xAD: xrah,
	0xAF: xral,
	0xB0: orab,
	0xB1: orac,
	0xB2: orad,
	0xB3: orae,
	0xB4: orah,
	0xB5: oral,
	0xB7: oraa,
	0xB8: cmpb,
	0xB9: cmpc,
	0xBA: cmpd,
	0xBB: cmpe,
	0xBC: cmph,
	0xBD: cmpl,
	0xBF: cmpa,
	0xC2: jnz,
	0xC3: jmp,
	0xC7: rst0,
	0xC9: ret,
	0xCD: call,
	0xCF: rst1,
	0xD5: pushd,
	0xD7: rst2,
	0xDF: rst3,
	0xE6: ani,
	0xE7: rst4,
	0xEF: rst5,
	0xFE: cpi,
	0xF7: rst6,
	0xFF: rst7,
}

// 0x8F ADC A | A <- A + A + CY (Z, S, P, CY, AC)
func adca(c *Computer) error {
	return add(c, c.A, c.Flags.carry())
}

// 0x88 ADC B |	A <- A + B + CY (Z, S, P, CY, AC)
func adcb(c *Computer) error {
	return add(c, c.B, c.Flags.carry())
}

// 0x89 ADC C |	A <- A + C + CY (Z, S, P, CY, AC)
func adcc(c *Computer) error {
	return add(c, c.C, c.Flags.carry())
}

// 0x8A ADC D |	A <- A + D + CY (Z, S, P, CY, AC)
func adcd(c *Computer) error {
	return add(c, c.D, c.Flags.carry())
}

// 0x8B ADC E | A <- A + E + CY (Z, S, P, CY, AC)
func adce(c *Computer) error {
	return add(c, c.E, c.Flags.carry())
}

// 0x8C ADC E | A <- A + E + CY (Z, S, P, CY, AC)
func adch(c *Computer) error {
	return add(c, c.H, c.Flags.carry())
}

// 0x8D ADC L | A <- A + L + CY (Z, S, P, CY, AC)
func adcl(c *Computer) error {
	return add(c, c.L, c.Flags.carry())
}

// 0x8E ADC M | A <- A + (HL) + CY (Z, S, P, CY, AC)
func adcm(c *Computer) error {
	v, err := c.read8Indirect()
	if err != nil {
		return err
	}
	return add(c, v, c.Flags.carry())
}

// 0x87 ADD A | A <- A + A (Z, S, P, CY, AC)
func adda(c *Computer) error {
	return add(c, c.A, false)
}

// 0x80 ADD B |	A <- A + B (Z, S, P, CY, AC)
func addb(c *Computer) error {
	return add(c, c.B, false)
}

// 0x81 ADD C | A <- A + C (Z, S, P, CY, AC)
func addc(c *Computer) error {
	return add(c, c.C, false)
}

// 0x82 ADD D | A <- A + D (Z, S, P, CY, AC)
func addd(c *Computer) error {
	return add(c, c.D, false)
}

// 0x83 ADD E | A <- A + E (Z, S, P, CY, AC)
func adde(c *Computer) error {
	return add(c, c.E, false)
}

// 0x84 ADD E | A <- A + E (Z, S, P, CY, AC)
func addh(c *Computer) error {
	return add(c, c.H, false)
}

// 0x85 ADD L | A <- A + L (Z, S, P, CY, AC)
func addl(c *Computer) error {
	return add(c, c.L, false)
}

// 0x86 ADD M | A <- A + (HL) (Z, S, P, CY, AC)
func addm(c *Computer) error {
	v, err := c.read8Indirect()
	if err != nil {
		return err
	}
	return add(c, v, false)
}

// 0xA7 ANA A | A <- A & A (Z, S, P, CY)
func anaa(c *Computer) error {
	return ana(c, c.A)
}

// 0xA0 ANA B | A <- A & B (Z, S, P, CY)
func anab(c *Computer) error {
	return ana(c, c.B)
}

// 0xA1 ANA C | A <- A & C (Z, S, P, CY)
func anac(c *Computer) error {
	return ana(c, c.C)
}

// 0xA2 ANA D | A <- A & D (Z, S, P, CY)
func anad(c *Computer) error {
	return ana(c, c.D)
}

// 0xA3 ANA E | A <- A & E (Z, S, P, CY)
func anae(c *Computer) error {
	return ana(c, c.E)
}

// 0xA4 ANA H | A <- A & H (Z, S, P, CY)
func anah(c *Computer) error {
	return ana(c, c.H)
}

// 0xA5 ANA L | A <- A & L (Z, S, P, CY)
func anal(c *Computer) error {
	return ana(c, c.L)
}

// 0xE6 ANI D8 | A <- A & data (Z, S, P, CY, AC)
// CY and AC are cleared
func ani(c *Computer) error {
	v, err := c.read8(c.PC + 1)
	if err != nil {
		return err
	}
	and := c.A & v
	flags := zero8(and) | sign8(and) | parity8(and)

	c.Flags = flags
	c.A = and
	c.PC += 2
	return nil
}

// 0xCD: CALL adr | (SP-1)<-PC.hi;(SP-2)<-PC.lo;SP<-SP-2;PC=adr
func call(c *Computer) error {
	addr, err := c.read16(c.PC + 1)
	if err != nil {
		return err
	}
	return calladdr(c, addr)
}

// 0xBF	CMP A | A - A (Z, S, P, CY, AC)
func cmpa(c *Computer) error {
	return sub(c, c.A, false)
}

// 0xB8	CMP B | A - B (Z, S, P, CY, AC)
func cmpb(c *Computer) error {
	return sub(c, c.B, false)
}

// 0xB9 CMP C | A - C (Z, S, P, CY, AC)
func cmpc(c *Computer) error {
	return sub(c, c.C, false)
}

// 0xBA CMP D | A - D (Z, S, P, CY, AC)
func cmpd(c *Computer) error {
	return sub(c, c.D, false)
}

//0xBB CMP E | A - E (Z, S, P, CY, AC)
func cmpe(c *Computer) error {
	return sub(c, c.E, false)
}

// 0xBC	CMP H | A - H (Z, S, P, CY, AC)
func cmph(c *Computer) error {
	return sub(c, c.H, false)
}

// 0xBD	CMP L | A - L (Z, S, P, CY, AC)
func cmpl(c *Computer) error {
	return sub(c, c.L, false)
}

//0xFE CPI D8 | A - data (Z, S, P, CY, AC)
func cpi(c *Computer) error {
	v, err := c.read8(c.PC + 1)
	if err != nil {
		return err
	}
	err = sub(c, v, false)
	if err != nil {
		return err
	}
	c.PC++
	return nil
}

// 0x09	DAD B | HL = HL + BC (CY)
func dadb(c *Computer) error {
	return dad(c, c.BC())
}

// 0x19	DAD D | HL = HL + DE (CY)
func dadd(c *Computer) error {
	return dad(c, c.DE())
}

// 0x29	DAD H | HL = HL + HL (CY)
func dadh(c *Computer) error {
	return dad(c, c.HL())
}

// 0x39	DAD B | HL = HL + SP (CY)
func dadsp(c *Computer) error {
	return dad(c, c.SP)
}

// 0x3D	DCR A | A <- A -1 (Z, S, P, AC)
func dcra(c *Computer) error {
	return dcr(c, &c.A)
}

// 0x05	DCR B | B <- B -1 (Z, S, P, AC)
func dcrb(c *Computer) error {
	return dcr(c, &c.B)
}

// 0x0D	DCR C | C <- C -1 (Z, S, P, AC)
func dcrc(c *Computer) error {
	return dcr(c, &c.C)
}

// 0x15	DCR D | D <- D -1 (Z, S, P, AC)
func dcrd(c *Computer) error {
	return dcr(c, &c.D)
}

// 0x1D	DCR E | E <- E -1 (Z, S, P, AC)
func dcre(c *Computer) error {
	return dcr(c, &c.E)
}

// 0x20	DCR H | H <- H -1 (Z, S, P, AC)
func dcrh(c *Computer) error {
	return dcr(c, &c.H)
}

// 0x2D	DCR L | L <- L -1 (Z, S, P, AC)
func dcrl(c *Computer) error {
	return dcr(c, &c.L)
}

// 0x3C	INR A | A <- A+1 (Z, S, P, AC)
func inra(c *Computer) error {
	return inr(c, &c.A)
}

// 0x04	INR B | B <- B+1 (Z, S, P, AC)
func inrb(c *Computer) error {
	return inr(c, &c.B)
}

// 0x0C	INR C | C <- C+1 (Z, S, P, AC)
func inrc(c *Computer) error {
	return inr(c, &c.C)
}

// 0x14	INR D | D <- D+1 (Z, S, P, AC)
func inrd(c *Computer) error {
	return inr(c, &c.D)
}

// 0x1C	INR E | E <- E+1 (Z, S, P, AC)
func inre(c *Computer) error {
	return inr(c, &c.E)
}

// 0x24	INR H | H <- H+1 (Z, S, P, AC)
func inrh(c *Computer) error {
	return inr(c, &c.H)
}

// 0x2C	INR L | L <- L+1 (Z, S, P, AC)
func inrl(c *Computer) error {
	return inr(c, &c.L)
}

// 0x03: INX BC | BC <- BC + 1
func inxb(c *Computer) error {
	return inx(c, &c.B, &c.C)
}

// 0x13: INX DE | DE <- DE + 1
func inxd(c *Computer) error {
	return inx(c, &c.D, &c.E)
}

// 0x23: INX H | H <- H + 1
func inxh(c *Computer) error {
	return inx(c, &c.H, &c.L)
}

// 0x33: INX SP | SP <- SP + 1
func inxsp(c *Computer) error {
	return inx16(c, &c.SP)
}

// 0xC3: JMP adr | PC <- adr.
// Jump to the address denoted by the next two bytes.
func jmp(c *Computer) error {
	addr, err := c.read16(c.PC + 1)
	if err != nil {
		return err
	}
	c.PC = addr
	return nil
}

// 0xC2: JNZ adr | if NZ, PC <- addr
// Jump to the address denoted by the next two bytes if the zero flag is set
func jnz(c *Computer) error {
	if !c.Flags.zero() {
		return jmp(c)
	} else {
		c.PC += 3
		return nil
	}
}

// 0x0A: LDAX B | A <- (BC)
func ldaxb(c *Computer) error {
	return ldax(c, c.B, c.C)
}

// 0x1A: LDAX D | A <- (DE)
func ldaxd(c *Computer) error {
	return ldax(c, c.D, c.E)
}

// 0x01: LXI B | D16. B <- byte 3, C <- byte 2
func lxib(c *Computer) error {
	return lxi(c, &c.B, &c.C)
}

// 0x11: LXI D | D16. D <- byte 3, E <- byte 2
func lxid(c *Computer) error {
	return lxi(c, &c.D, &c.E)
}

// 0x21: LXI H, D161 | H <- byte 3, L <- byte 2
func lxih(c *Computer) error {
	return lxi(c, &c.H, &c.L)
}

// 0x31: LXI SP, D16 | SP.hi <- byte 3, SP.lo <- byte 2
// Resets the stack pointer to a given value
func lxisp(c *Computer) error {
	return lxi16(c, &c.SP)
}

// 0x7F: MOV A, A | A <- A
func movaa(c *Computer) error {
	return nop(c)
}

// 0x78: MOV A, B | A <- B
func movab(c *Computer) error {
	return mov(c, &c.A, &c.B)
}

// 0x79: MOV A, C | A <- C
func movac(c *Computer) error {
	return mov(c, &c.A, &c.C)
}

// 0x7A: MOV A, D | A <- D
func movad(c *Computer) error {
	return mov(c, &c.A, &c.D)
}

// 0x7B: MOV A, E | A <- E
func movae(c *Computer) error {
	return mov(c, &c.A, &c.E)
}

// 0x7C: MOV A, H | A <- H
func movah(c *Computer) error {
	return mov(c, &c.A, &c.H)
}

// 0x7D: MOV A, L | A <- L
func moval(c *Computer) error {
	return mov(c, &c.A, &c.L)
}

// 0x47: MOV B, A | B <- A
func movba(c *Computer) error {
	return mov(c, &c.B, &c.A)
}

// 0x40: MOV B, B | B <- B
func movbb(c *Computer) error {
	return nop(c)
}

// 0x41: MOV B, C | B <- C
func movbc(c *Computer) error {
	return mov(c, &c.B, &c.C)
}

// 0x42: MOV B, D | B <- D
func movbd(c *Computer) error {
	return mov(c, &c.B, &c.D)
}

// 0x43: MOV B, E | B <- E
func movbe(c *Computer) error {
	return mov(c, &c.B, &c.E)
}

// 0x44: MOV B, H | B <- H
func movbh(c *Computer) error {
	return mov(c, &c.B, &c.H)
}

// 0x45: MOV B, L | B <- L
func movbl(c *Computer) error {
	return mov(c, &c.B, &c.L)
}

// 0x4F: MOV C, A | C <- A
func movca(c *Computer) error {
	return mov(c, &c.C, &c.A)
}

// 0x48: MOV C, B | C <- B
func movcb(c *Computer) error {
	return mov(c, &c.C, &c.B)
}

// 0x49: MOV C, C | C <- C
func movcc(c *Computer) error {
	return nop(c)
}

// 0x4A: MOV C, D | C <- D
func movcd(c *Computer) error {
	return mov(c, &c.C, &c.D)
}

// 0x4B: MOV C, E | C <- E
func movce(c *Computer) error {
	return mov(c, &c.C, &c.E)
}

// 0x4C: MOV C, H | C <- H
func movch(c *Computer) error {
	return mov(c, &c.C, &c.H)
}

// 0x4D: MOV C, L | C <- L
func movcl(c *Computer) error {
	return mov(c, &c.C, &c.L)
}

// 0x57: MOV D, A | D <- A
func movda(c *Computer) error {
	return mov(c, &c.D, &c.A)
}

// 0x50: MOV D, B | D <- B
func movdb(c *Computer) error {
	return mov(c, &c.D, &c.B)
}

// 0x51: MOV D, C | D <- C
func movdc(c *Computer) error {
	return mov(c, &c.D, &c.C)
}

// 0x52: MOV D, D | D <- D
func movdd(c *Computer) error {
	return nop(c)
}

// 0x53: MOV D, E | D <- E
func movde(c *Computer) error {
	return mov(c, &c.D, &c.E)
}

// 0x54: MOV D, H | D <- H
func movdh(c *Computer) error {
	return mov(c, &c.D, &c.H)
}

// 0x55: MOV D, L | D <- L
func movdl(c *Computer) error {
	return mov(c, &c.D, &c.L)
}

// 0x5F: MOV E, A | E <- A
func movea(c *Computer) error {
	return mov(c, &c.E, &c.A)
}

// 0x58: MOV E, B | E <- B
func moveb(c *Computer) error {
	return mov(c, &c.E, &c.B)
}

// 0x59: MOV E, C | E <- C
func movec(c *Computer) error {
	return mov(c, &c.E, &c.C)
}

// 0x5A: MOV E, D | E <- D
func moved(c *Computer) error {
	return mov(c, &c.E, &c.D)
}

// 0x5B: MOV E, E | E <- E
func movee(c *Computer) error {
	return nop(c)
}

// 0x5C: MOV E, H | E <- H
func moveh(c *Computer) error {
	return mov(c, &c.E, &c.H)
}

// 0x46	MOV B,M | B <- (HL)
func movfrommb(c *Computer) error {
	return movfromm(c, &c.B)
}

// 0x4E	MOV C,M | C <- (HL)
func movfrommc(c *Computer) error {
	return movfromm(c, &c.C)
}

// 0x56	MOV D,M | D <- (HL)
func movfrommd(c *Computer) error {
	return movfromm(c, &c.D)
}

// 0x5E	MOV E,M | E <- (HL)
func movfromme(c *Computer) error {
	return movfromm(c, &c.E)
}

// 0x66	MOV H,M | H <- (HL)
func movfrommh(c *Computer) error {
	return movfromm(c, &c.H)
}

// 0x6E	MOV L,M | L <- (HL)
func movfromml(c *Computer) error {
	return movfromm(c, &c.L)
}

// 0x7E	MOV A,M | A <- (HL)
func movfromma(c *Computer) error {
	return movfromm(c, &c.A)
}

// 0x5D: MOV E, L | E <- L
func movel(c *Computer) error {
	return mov(c, &c.E, &c.L)
}

// 0x67: MOV H, A | H <- A
func movha(c *Computer) error {
	return mov(c, &c.H, &c.A)
}

// 0x60: MOV H, B | H <- B
func movhb(c *Computer) error {
	return mov(c, &c.H, &c.B)
}

// 0x61: MOV H, C | H <- C
func movhc(c *Computer) error {
	return mov(c, &c.H, &c.C)
}

// 0x62: MOV H, D | H <- D
func movhd(c *Computer) error {
	return mov(c, &c.H, &c.D)
}

// 0x63: MOV H, E | H <- E
func movhe(c *Computer) error {
	return mov(c, &c.H, &c.E)
}

// 0x64: MOV H, H | H <- H
func movhh(c *Computer) error {
	return nop(c)
}

// 0x65: MOV H, L | H <- L
func movhl(c *Computer) error {
	return mov(c, &c.H, &c.L)
}

// 0x6F: MOV L, A | L <- A
func movla(c *Computer) error {
	return mov(c, &c.L, &c.A)
}

// 0x68: MOV L, B | L <- B
func movlb(c *Computer) error {
	return mov(c, &c.L, &c.B)
}

// 0x69: MOV L, C | L <- C
func movlc(c *Computer) error {
	return mov(c, &c.L, &c.C)
}

// 0x6A: MOV L, D | L <- D
func movld(c *Computer) error {
	return mov(c, &c.L, &c.D)
}

// 0x6B: MOV L, E | L <- E
func movle(c *Computer) error {
	return mov(c, &c.L, &c.E)
}

// 0x6C: MOV L, H | L <- H
func movlh(c *Computer) error {
	return mov(c, &c.L, &c.H)
}

// 0x6D: MOV L, L | L <- L
func movll(c *Computer) error {
	return nop(c)
}

// 0x77: MOV M,A | (HL) <- A
func movtoma(c *Computer) error {
	return movtom(c, c.A)
}

// 0x77: MOV M,B | (HL) <- B
func movtomb(c *Computer) error {
	return movtom(c, c.B)
}

// 0x77: MOV M,C | (HL) <- C
func movtomc(c *Computer) error {
	return movtom(c, c.C)
}

// 0x77: MOV M,D | (HL) <- D
func movtomd(c *Computer) error {
	return movtom(c, c.D)
}

// 0x77: MOV M,E | (HL) <- E
func movtome(c *Computer) error {
	return movtom(c, c.E)
}

// 0x77: MOV M,H | (HL) <- H
func movtomh(c *Computer) error {
	return movtom(c, c.H)
}

// 0x77: MOV M,L | (HL) <- L
func movtoml(c *Computer) error {
	return movtom(c, c.L)
}

// 0x3E: MVI A, D8 | A <- byte 2
func mvia(c *Computer) error {
	return mvi(c, &c.A)
}

// 0x06: MVI B, D8 | B <- byte 2
func mvib(c *Computer) error {
	return mvi(c, &c.B)
}

// 0x0E: MVI C, D8 | C <- byte 2
func mvic(c *Computer) error {
	return mvi(c, &c.C)
}

// 0x16: MVI D, D8 | D <- byte 2
func mvid(c *Computer) error {
	return mvi(c, &c.D)
}

// 0x1E: MVI E, D8 | E <- byte 2
func mvie(c *Computer) error {
	return mvi(c, &c.E)
}

// 0x26: MVI H, D8 | H <- byte 2
func mvih(c *Computer) error {
	return mvi(c, &c.H)
}

// 0x2E: MVI L, D8 | L <- byte 2
func mvil(c *Computer) error {
	return mvi(c, &c.L)
}

// 0x36: MVI M, D8 | (H L) <- byte 2
func mvim(c *Computer) error {
	v, err := c.read8(c.PC + 1)
	if err != nil {
		return err
	}
	err = movtom(c, v)
	if err != nil {
		return err
	}
	c.PC++
	return nil
}

// 0x00: NOP
func nop(c *Computer) error {
	c.PC++
	return nil
}

// 0xb7	ORA A (Z, S, P, CY, AC) | A <- A | A
func oraa(c *Computer) error {
	return ora(c, c.A)
}

// 0xb0	ORA B (Z, S, P, CY, AC) | A <- A | B
func orab(c *Computer) error {
	return ora(c, c.B)
}

// 0xb1	ORA C (Z, S, P, CY, AC) | A <- A | C
func orac(c *Computer) error {
	return ora(c, c.C)
}

// 0xb2	ORA D (Z, S, P, CY, AC) | A <- A | D
func orad(c *Computer) error {
	return ora(c, c.D)
}

// 0xb3	ORA E (Z, S, P, CY, AC) | A <- A | E
func orae(c *Computer) error {
	return ora(c, c.E)
}

// 0xb4	ORA H (Z, S, P, CY, AC) | A <- A | H
func orah(c *Computer) error {
	return ora(c, c.H)
}

// 0xb5	ORA L (Z, S, P, CY, AC) | A <- A | L
func oral(c *Computer) error {
	return ora(c, c.L)
}

// 0xD5	PUSH D | (sp-2)<-E; (sp-1)<-D; sp <- sp - 2
func pushd(c *Computer) error {
	err := push(c, c.DE())
	if err != nil {
		return err
	}
	c.PC += 2
	return nil
}

// 0xC9 RET | PC.lo <- (sp); PC.hi<-(sp+1); SP <- SP+2
func ret(c *Computer) error {
	pc, err := pop(c)
	if err != nil {
		return err
	}
	c.PC = pc + 1
	return nil
}

//0xC7 RST 0 | CALL $0
func rst0(c *Computer) error {
	return calladdr(c, 0x0)
}

//0xCF RST 1 | CALL $8
func rst1(c *Computer) error {
	return calladdr(c, 0x08)
}

//0xD7 RST 2 | CALL $10
func rst2(c *Computer) error {
	return calladdr(c, 0x10)
}

//0xDF RST 3 | CALL $18
func rst3(c *Computer) error {
	return calladdr(c, 0x18)
}

//0xE7 RST 4 | CALL $20
func rst4(c *Computer) error {
	return calladdr(c, 0x20)
}

//0xEF RST 5 | CALL $28
func rst5(c *Computer) error {
	return calladdr(c, 0x28)
}

//0xF7 RST 6 | CALL $30
func rst6(c *Computer) error {
	return calladdr(c, 0x30)
}

//0xFF RST 7 | CALL $38
func rst7(c *Computer) error {
	return calladdr(c, 0x38)
}

// 0x9F SBB A | A <- A - A - CY (Z, S, P, CY, AC)
func sbba(c *Computer) error {
	return sub(c, c.A, c.Flags.carry())
}

// 0x98 SBB B | A <- A - B - CY (Z, S, P, CY, AC)
func sbbb(c *Computer) error {
	return sub(c, c.B, c.Flags.carry())
}

// 0x99 SBB C | A <- A - C - CY (Z, S, P, CY, AC)
func sbbc(c *Computer) error {
	return sub(c, c.C, c.Flags.carry())
}

// 0x9A SBB D | A <- A - D - CY (Z, S, P, CY, AC)
func sbbd(c *Computer) error {
	return sub(c, c.D, c.Flags.carry())
}

// 0x9B SBB E | A <- A - E - CY (Z, S, P, CY, AC)
func sbbe(c *Computer) error {
	return sub(c, c.E, c.Flags.carry())
}

// 0x9C SBB E | A <- A - E - CY (Z, S, P, CY, AC)
func sbbh(c *Computer) error {
	return sub(c, c.H, c.Flags.carry())
}

// 0x9D SBB L | A <- A - L - CY (Z, S, P, CY, AC)
func sbbl(c *Computer) error {
	return sub(c, c.L, c.Flags.carry())
}

// 0x32 STA addr
func sta(c *Computer) error {
	addr, err := c.read16(c.PC + 1)
	if err != nil {
		return err
	}
	err = c.write8(addr, c.A)
	if err != nil {
		return err
	}
	c.PC += 3
	return nil
}

// 0x02 STAX B | (BC) <- A
func staxb(c *Computer) error {
	return stax(c, c.B, c.C)
}

// 0x12 STAX D | (DE) <- A
func staxd(c *Computer) error {
	return stax(c, c.D, c.E)
}

// 0x97 SUB A | A <- A - A (Z, S, P, CY, AC)
func suba(c *Computer) error {
	return sub(c, c.A, false)
}

// 0x90 SUB B | A <- A - B (Z, S, P, CY, AC)
func subb(c *Computer) error {
	return sub(c, c.B, false)
}

// 0x91 SUB C | A <- A - C (Z, S, P, CY, AC)
func subc(c *Computer) error {
	return sub(c, c.C, false)
}

// 0x92 SUB D | A <- A - D (Z, S, P, CY, AC)
func subd(c *Computer) error {
	return sub(c, c.D, false)
}

// 0x93 SUB E | A <- A - E (Z, S, P, CY, AC)
func sube(c *Computer) error {
	return sub(c, c.E, false)
}

// 0x94 SUB E | A <- A - E (Z, S, P, CY, AC)
func subh(c *Computer) error {
	return sub(c, c.H, false)
}

// 0x95 SUB L | A <- A - L (Z, S, P, CY, AC)
func subl(c *Computer) error {
	return sub(c, c.L, false)
}

// 0xA8 XRA A | A <- A XOR A (Z, S, P, CY)
func xraa(c *Computer) error {
	return xra(c, c.A)
}

// 0xA9 XRA B | A <- A XOR B (Z, S, P, CY)
func xrab(c *Computer) error {
	return xra(c, c.B)
}

// 0xAA XRA C | A <- A XOR C (Z, S, P, CY)
func xrac(c *Computer) error {
	return xra(c, c.C)
}

// 0xAB XRA D | A <- A XOR D (Z, S, P, CY)
func xrad(c *Computer) error {
	return xra(c, c.D)
}

// 0xAC XRA E | A <- A XOR E (Z, S, P, CY)
func xrae(c *Computer) error {
	return xra(c, c.E)
}

// 0xAD XRA H | A <- A XOR H (Z, S, P, CY)
func xrah(c *Computer) error {
	return xra(c, c.H)
}

// 0xAF XRA L | A <- A XOR L (Z, S, P, CY)
func xral(c *Computer) error {
	return xra(c, c.L)
}

/****************************************/
/*			Support functions 			*/
/****************************************/

func add(c *Computer, v byte, carry bool) error {
	if carry {
		v++
	}
	sum := c.A + v

	flags := zero8(sum) | sign8(sum) | parity8(sum)
	if sum < c.A {
		flags |= cf
	}
	if c.A&0x0F+v&0x0F >= 0x10 {
		flags |= hc
	}

	c.A = sum
	c.Flags = flags
	c.PC++
	return nil
}

func ana(c *Computer, v byte) error {
	c.A &= v
	c.Flags = zero8(c.A) | sign8(c.A) | parity8(c.A) | (c.Flags & hc)
	c.PC++
	return nil
}

func calladdr(c *Computer, addr uint16) error {
	err := push(c, c.PC)
	if err != nil {
		return err
	}

	c.PC = addr
	return nil
}

func dad(c *Computer, d16 uint16) error {
	s := c.HL()
	sum := s + d16

	c.H = byte(sum >> 8)
	c.L = byte(sum & 0x00FF)

	// affect only the carry flag
	if sum < s {
		c.Flags |= cf
	} else {
		c.Flags &= 0xFF ^ cf
	}

	c.PC++
	return nil
}

func dcr(c *Computer, reg *byte) error {
	sum := *reg - 1

	flags := zero8(sum) | sign8(sum) | parity8(sum) | (c.Flags & cf)

	if *reg&0x0F >= 0x01 {
		flags |= hc
	}

	*reg = sum
	c.Flags = flags
	c.PC++
	return nil
}

func inr(c *Computer, reg *byte) error {
	sum := *reg + 1

	flags := zero8(sum) | sign8(sum) | parity8(sum) | (c.Flags & cf)

	// there was auxiliary carry if there was carry between bit 3 and bit 4 of the resulting value.
	if sum&0x10 == 0x10 {
		flags |= hc
	}

	*reg = sum
	c.Flags = flags
	c.PC++
	return nil
}

func inx(c *Computer, msreg, lsreg *byte) error {
	incr := (uint16(*msreg)<<8 + uint16(*lsreg)) + 1
	*msreg = byte(incr >> 8)
	*lsreg = byte(incr & 0xFF)
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

func ldax(c *Computer, msb, lsb byte) error {
	addr := uint16(msb)<<8 + uint16(lsb)
	v, err := c.read8(addr)
	if err != nil {
		return err
	}

	c.A = v
	c.PC++
	return nil
}

func lxi(c *Computer, msreg, lsreg *byte) error {
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
	c.PC += 3
	return nil
}

func mov(c *Computer, dstreg, srcreg *byte) error {
	*dstreg = *srcreg
	c.PC++
	return nil
}

func movfromm(c *Computer, reg *byte) error {
	v, err := c.read8Indirect()
	if err != nil {
		return err
	}
	*reg = v
	c.PC++
	return nil
}

func movtom(c *Computer, v byte) error {
	err := c.write8Indirect(v)
	if err != nil {
		return err
	}
	c.PC++
	return nil
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

func ora(c *Computer, v byte) error {
	c.A |= v
	c.Flags = zero8(c.A) | sign8(c.A) | parity8(c.A)
	c.PC++
	return nil
}

func pop(c *Computer) (uint16, error) {
	v, err := c.read16(c.SP)
	if err != nil {
		return 0, err
	}
	c.SP += 2
	return v, nil
}

func push(c *Computer, d16 uint16) error {
	lsb := byte(d16 & 0x00FF)
	msb := byte(d16 >> 8)

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

func stax(c *Computer, msb, lsb byte) error {
	addr := uint16(msb)<<8 + uint16(lsb)
	err := c.write8(addr, c.A)
	c.PC++
	return err
}

func sub(c *Computer, v byte, borrow bool) error {
	if borrow {
		v++
	}
	sub := c.A + (^v + 1)

	flags := zero8(sub) | sign8(sub) | parity8(sub)
	// there was borrow (cf = 1) if v is higher than the minuend
	if sub > c.A {
		flags |= cf
	}
	// there was auxiliary carry if there was carry between bit 3 and bit 4 of the resulting value
	// it seems counterintuitive, that the behavior is the same of additions, but check this stackexchange
	// answer: https://retrocomputing.stackexchange.com/a/12560
	if c.A&0x0F+v&0x0F >= 0x10 {
		flags |= hc
	}

	c.A = sub
	c.Flags = flags
	c.PC++
	return nil
}

func xra(c *Computer, v byte) error {
	c.A ^= v
	c.Flags = zero8(c.A) | sign8(c.A) | parity8(c.A)
	c.PC++
	return nil
}
