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
		c.Registers.Pc++
	} else if opcode == 0xC3 {
		low := c.Bus.Read(c.Registers.Pc + 1)
		high := c.Bus.Read(c.Registers.Pc + 2)
		c.Registers.Pc = uint16(low) | (uint16(high) << 8)
	} else if opcode == 0xAF {
		c.Registers.SetFlag(7, false)
		c.Registers.SetFlag(6, false)
		c.Registers.SetFlag(5, false)
		c.Registers.SetFlag(4, false)
		c.Registers.Pc++
	} else {
		panic(fmt.Sprintf("not implemented\nopcode: %02X\naddress: %04X", opcode, c.Registers.Pc))
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
	cartridge, err := CartridgeFromFile("tetris.gb")
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
