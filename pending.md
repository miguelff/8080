Opcode	Instruction	size	flags	function
0x02	STAX B	1		(BC) <- A
0x05	DCR B	1	Z, S, P, AC	B <- B-1
0x07	RLC	1	CY	A = A << 1; bit 0 = prev bit 7; CY = prev bit 7
0x08	-			
0x09	DAD B	1	CY	HL = HL + BC
0x0a	LDAX B	1		A <- (BC)
0x0b	DCX B	1		BC = BC-1
0x0d	DCR C	1	Z, S, P, AC	C <-C-1
0x0f	RRC	1	CY	A = A >> 1; bit 7 = prev bit 0; CY = prev bit 0
0x10	-			
0x12	STAX D	1		(DE) <- A
0x15	DCR D	1	Z, S, P, AC	D <- D-1
0x17	RAL	1	CY	A = A << 1; bit 0 = prev CY; CY = prev bit 7
0x18	-			
0x19	DAD D	1	CY	HL = HL + DE
0x1b	DCX D	1		DE = DE-1
0x1d	DCR E	1	Z, S, P, AC	E <- E-1
0x1f	RAR	1	CY	A = A >> 1; bit 7 = prev bit 7; CY = prev bit 0
0x20	-			
0x22	SHLD adr	3		(adr) <-L; (adr+1)<-H
0x25	DCR H	1	Z, S, P, AC	H <- H-1
0x27	DAA	1		special
0x28	-			
0x29	DAD H	1	CY	HL = HL + HI
0x2a	LHLD adr	3		L <- (adr); H<-(adr+1)
0x2b	DCX H	1		HL = HL-1
0x2d	DCR L	1	Z, S, P, AC	L <- L-1
0x2f	CMA	1		A <- !A
0x30	-			
0x32	STA adr	3		(adr) <- A
0x34	INR M	1	Z, S, P, AC	(HL) <- (HL)+1
0x35	DCR M	1	Z, S, P, AC	(HL) <- (HL)-1
0x37	STC	1	CY	CY = 1
0x38	-			
0x39	DAD SP	1	CY	HL = HL + SP
0x3a	LDA adr	3		A <- (adr)
0x3b	DCX SP	1		SP = SP-1
0x3d	DCR A	1	Z, S, P, AC	A <- A-1
0x3f	CMC	1	CY	CY=!CY
0x46	MOV B,M	1		B <- (HL)
0x4e	MOV C,M	1		C <- (HL)
0x56	MOV D,M	1		D <- (HL)
0x5e	MOV E,M	1		E <- (HL)
0x66	MOV H,M	1		H <- (HL)
0x6e	MOV L,M	1		L <- (HL)
0x76	HLT	1		special
0x7e	MOV A,M	1		A <- (HL)
0x86	ADD M	1	Z, S, P, CY, AC	A <- A + (HL)
0x8e	ADC M	1	Z, S, P, CY, AC	A <- A + (HL) + CY
0x90	SUB B	1	Z, S, P, CY, AC	A <- A - B
0x91	SUB C	1	Z, S, P, CY, AC	A <- A - C
0x92	SUB D	1	Z, S, P, CY, AC	A <- A + D
0x93	SUB E	1	Z, S, P, CY, AC	A <- A - E
0x94	SUB H	1	Z, S, P, CY, AC	A <- A + H
0x95	SUB L	1	Z, S, P, CY, AC	A <- A - L
0x96	SUB M	1	Z, S, P, CY, AC	A <- A + (HL)
0x97	SUB A	1	Z, S, P, CY, AC	A <- A - A
0x98	SBB B	1	Z, S, P, CY, AC	A <- A - B - CY
0x99	SBB C	1	Z, S, P, CY, AC	A <- A - C - CY
0x9a	SBB D	1	Z, S, P, CY, AC	A <- A - D - CY
0x9b	SBB E	1	Z, S, P, CY, AC	A <- A - E - CY
0x9c	SBB H	1	Z, S, P, CY, AC	A <- A - H - CY
0x9d	SBB L	1	Z, S, P, CY, AC	A <- A - L - CY
0x9e	SBB M	1	Z, S, P, CY, AC	A <- A - (HL) - CY
0x9f	SBB A	1	Z, S, P, CY, AC	A <- A - A - CY
0xa0	ANA B	1	Z, S, P, CY, AC	A <- A & B
0xa1	ANA C	1	Z, S, P, CY, AC	A <- A & C
0xa2	ANA D	1	Z, S, P, CY, AC	A <- A & D
0xa3	ANA E	1	Z, S, P, CY, AC	A <- A & E
0xa4	ANA H	1	Z, S, P, CY, AC	A <- A & H
0xa5	ANA L	1	Z, S, P, CY, AC	A <- A & L
0xa6	ANA M	1	Z, S, P, CY, AC	A <- A & (HL)
0xa7	ANA A	1	Z, S, P, CY, AC	A <- A & A
0xa8	XRA B	1	Z, S, P, CY, AC	A <- A ^ B
0xa9	XRA C	1	Z, S, P, CY, AC	A <- A ^ C
0xaa	XRA D	1	Z, S, P, CY, AC	A <- A ^ D
0xab	XRA E	1	Z, S, P, CY, AC	A <- A ^ E
0xac	XRA H	1	Z, S, P, CY, AC	A <- A ^ H
0xad	XRA L	1	Z, S, P, CY, AC	A <- A ^ L
0xae	XRA M	1	Z, S, P, CY, AC	A <- A ^ (HL)
0xaf	XRA A	1	Z, S, P, CY, AC	A <- A ^ A
0xb0	ORA B	1	Z, S, P, CY, AC	A <- A | B
0xb1	ORA C	1	Z, S, P, CY, AC	A <- A | C
0xb2	ORA D	1	Z, S, P, CY, AC	A <- A | D
0xb3	ORA E	1	Z, S, P, CY, AC	A <- A | E
0xb4	ORA H	1	Z, S, P, CY, AC	A <- A | H
0xb5	ORA L	1	Z, S, P, CY, AC	A <- A | L
0xb6	ORA M	1	Z, S, P, CY, AC	A <- A | (HL)
0xb7	ORA A	1	Z, S, P, CY, AC	A <- A | A
0xb8	CMP B	1	Z, S, P, CY, AC	A - B
0xb9	CMP C	1	Z, S, P, CY, AC	A - C
0xba	CMP D	1	Z, S, P, CY, AC	A - D
0xbb	CMP E	1	Z, S, P, CY, AC	A - E
0xbc	CMP H	1	Z, S, P, CY, AC	A - H
0xbd	CMP L	1	Z, S, P, CY, AC	A - L
0xbe	CMP M	1	Z, S, P, CY, AC	A - (HL)
0xbf	CMP A	1	Z, S, P, CY, AC	A - A
0xc0	RNZ	1		if NZ, RET
0xc1	POP B	1		C <- (sp); B <- (sp+1); sp <- sp+2
0xc2	JNZ adr	3		if NZ, PC <- adr
0xc4	CNZ adr	3		if NZ, CALL adr
0xc5	PUSH B	1		(sp-2)<-C; (sp-1)<-B; sp <- sp - 2
0xc6	ADI D8	2	Z, S, P, CY, AC	A <- A + byte
0xc7	RST 0	1		CALL $0
0xc8	RZ	1		if Z, RET
0xc9	RET	1		PC.lo <- (sp); PC.hi<-(sp+1); SP <- SP+2
0xca	JZ adr	3		if Z, PC <- adr
0xcb	-			
0xcc	CZ adr	3		if Z, CALL adr
0xce	ACI D8	2	Z, S, P, CY, AC	A <- A + data + CY
0xcf	RST 1	1		CALL $8
0xd0	RNC	1		if NCY, RET
0xd1	POP D	1		E <- (sp); D <- (sp+1); sp <- sp+2
0xd2	JNC adr	3		if NCY, PC<-adr
0xd3	OUT D8	2		special
0xd4	CNC adr	3		if NCY, CALL adr
0xd5	PUSH D	1		(sp-2)<-E; (sp-1)<-D; sp <- sp - 2
0xd6	SUI D8	2	Z, S, P, CY, AC	A <- A - data
0xd7	RST 2	1		CALL $10
0xd8	RC	1		if CY, RET
0xd9	-			
0xda	JC adr	3		if CY, PC<-adr
0xdb	IN D8	2		special
0xdc	CC adr	3		if CY, CALL adr
0xdd	-			
0xde	SBI D8	2	Z, S, P, CY, AC	A <- A - data - CY
0xdf	RST 3	1		CALL $18
0xe0	RPO	1		if PO, RET
0xe1	POP H	1		L <- (sp); H <- (sp+1); sp <- sp+2
0xe2	JPO adr	3		if PO, PC <- adr
0xe3	XTHL	1		L <-> (SP); H <-> (SP+1)
0xe4	CPO adr	3		if PO, CALL adr
0xe5	PUSH H	1		(sp-2)<-L; (sp-1)<-H; sp <- sp - 2
0xe6	ANI D8	2	Z, S, P, CY, AC	A <- A & data
0xe7	RST 4	1		CALL $20
0xe8	RPE	1		if PE, RET
0xe9	PCHL	1		PC.hi <- H; PC.lo <- L
0xea	JPE adr	3		if PE, PC <- adr
0xeb	XCHG	1		H <-> D; L <-> E
0xec	CPE adr	3		if PE, CALL adr
0xed	-			
0xee	XRI D8	2	Z, S, P, CY, AC	A <- A ^ data
0xef	RST 5	1		CALL $28
0xf0	RP	1		if P, RET
0xf1	POP PSW	1		flags <- (sp); A <- (sp+1); sp <- sp+2
0xf2	JP adr	3		if P=1 PC <- adr
0xf3	DI	1		special
0xf4	CP adr	3		if P, PC <- adr
0xf5	PUSH PSW	1		(sp-2)<-flags; (sp-1)<-A; sp <- sp - 2
0xf6	ORI D8	2	Z, S, P, CY, AC	A <- A | data
0xf7	RST 6	1		CALL $30
0xf8	RM	1		if M, RET
0xf9	SPHL	1		SP=HL
0xfa	JM adr	3		if M, PC <- adr
0xfb	EI	1		special
0xfc	CM adr	3		if M, CALL adr
0xfd	-			
0xfe	CPI D8	2	Z, S, P, CY, AC	A - data
0xff	RST 7	1		CALL $38