package main

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
)

type Cpu struct {
	Registers *Registers
	Bus       *Bus
}

func (c *Cpu) Step() {
	opcode := c.Bus.Read(c.Registers.Pc)

	if opcode == 0x00 {
		c.Registers.Pc += 1
	} else if opcode == 0xC3 {
		low := c.Bus.Read(c.Registers.Pc + 1)
		high := c.Bus.Read(c.Registers.Pc + 2)
		c.Registers.Pc = uint16(low) | (uint16(high) << 8)
	} else if opcode == 0xAF { // XOR A
		c.Registers.SetFlag(7, false)
		c.Registers.SetFlag(6, false)
		c.Registers.SetFlag(5, false)
		c.Registers.SetFlag(4, false)
		c.Registers.Pc += 1
	} else if opcode == 0x21 { // LD BC, d16
		low := c.Bus.Read(c.Registers.Pc + 1)
		high := c.Bus.Read(c.Registers.Pc + 2)
		c.Registers.B = low
		c.Registers.C = high
		c.Registers.Pc += 3
	} else if opcode == 0x0E { // LD C, d8
		c.Registers.C = c.Bus.Read(c.Registers.Pc + 1)
		c.Registers.Pc += 2
	} else if opcode == 0x06 { // LD B, d8
		c.Registers.B = c.Bus.Read(c.Registers.Pc + 1)
		c.Registers.Pc += 2
	} else if opcode == 0x32 { // LD (HL+), A
		// FIXME:
		// addr := uint16(c.Registers.H) | uint16(c.Registers.L)<<8
		// c.Bus.Write(addr, c.Registers.A)

		hl := c.Registers.GetPair(&c.Registers.H, &c.Registers.L)
		c.Registers.SetPair(&c.Registers.H, &c.Registers.L, hl-1)

		// c.Registers.L -= 1
		c.Registers.Pc += 1

	} else if opcode == 0x05 { // DEC B
		c.Registers.B -= 1
		c.Registers.SetFlag(F_SUBTRACT, true)

		if c.Registers.B == 0 {
			c.Registers.SetFlag(F_ZERO, true)
		}
		c.Registers.Pc += 1
	} else if opcode == 0x0D { // DEC D
		c.Registers.D -= 1
		c.Registers.SetFlag(F_SUBTRACT, true)

		if c.Registers.D == 0 {
			c.Registers.SetFlag(F_ZERO, true)
		}
		c.Registers.Pc += 1
	} else if opcode == 0x20 { // JNZ r8
		if !c.Registers.GetFlag(F_ZERO) {
			r8 := int8(c.Bus.Read(c.Registers.Pc + 1))
			c.Registers.Pc = uint16(int32(c.Registers.Pc+2) + int32(r8))
		} else {
			c.Registers.Pc += 2
		}
	} else if opcode == 0x3E { // LD A, d8
		d8 := c.Bus.Read(c.Registers.Pc + 1)
		c.Registers.A = d8
		c.Registers.Pc += 2
	} else {
		panic(fmt.Sprintf("(not implemented) opcode: %02X, address: %04X", opcode, c.Registers.Pc))
	}

}

func TestCpuStepBlankCartridge(t *testing.T) {
	cpu := Cpu{
		Registers: &Registers{},
		Bus: &Bus{
			Cartridge: BlankCartridge(),
		},
	}

	assert.Equal(t, uint16(0), cpu.Registers.Pc)

	cpu.Step()
	assert.Equal(t, uint16(1), cpu.Registers.Pc)

	cpu.Step()
	assert.Equal(t, uint16(2), cpu.Registers.Pc)
}

func TestCpuStepTetris(t *testing.T) {
	cartridge, err := CartridgeFromFile("roms/tetris.gb")
	assert.NoError(t, err)

	cpu := Cpu{
		Registers: &Registers{Pc: 0x0100},
		Bus: &Bus{
			Cartridge: cartridge,
		},
	}

	for {
		t.Logf("address: %04X, opcode: %02X\n", cpu.Registers.Pc, cpu.Bus.Read(cpu.Registers.Pc))
		cpu.Step()
	}
}
