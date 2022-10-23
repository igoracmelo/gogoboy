package main

type Registers struct {
	A  uint8
	B  uint8
	C  uint8
	D  uint8
	E  uint8
	F  uint8
	G  uint8
	H  uint8
	L  uint8
	Pc uint16
	Sp uint16
}

func (r *Registers) GetBC() uint16 {
	return (uint16(r.B) << 8) + uint16(r.C)
}

func (r *Registers) GetFlag(shift uint8) bool {
	return ((r.F >> shift) & 0x01) != 0
}

func (r *Registers) SetFlag(shift uint8, value bool) {
	var bit uint8 = 0
	if value {
		bit = 1
	}
	bit = bit << shift

	var mask uint8 = ^(0x01 << shift)
	r.F = (r.F & mask) | bit
}
