package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

const (
	F_ZERO       uint8 = 7
	F_SUBTRACT   uint8 = 6
	F_HALF_CARRY uint8 = 5
	F_CARRY      uint8 = 4
)

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

func (r *Registers) GetPair(r1 *uint8, r2 *uint8) uint16 {
	return uint16(*r1)<<8 | uint16(*r2)
}

func (r *Registers) SetPair(r1 *uint8, r2 *uint8, value uint16) {
	high := uint8(value >> 8)
	low := uint8(value & 0x00FF)
	*r1 = high
	*r2 = low
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

func TestFlags(t *testing.T) {
	regs := Registers{F: 0b1111_0000}

	assert.True(t, regs.GetFlag(F_ZERO))

	regs.SetFlag(F_ZERO, false)
	assert.False(t, regs.GetFlag(F_ZERO))
	assert.Equal(t, uint8(0b0111_0000), regs.F)

	regs.SetFlag(F_SUBTRACT, false)
	assert.False(t, regs.GetFlag(F_SUBTRACT))

	regs.SetFlag(F_SUBTRACT, false)
	assert.False(t, regs.GetFlag(F_SUBTRACT))

	assert.Equal(t, uint8(0b0011_0000), regs.F)
}
