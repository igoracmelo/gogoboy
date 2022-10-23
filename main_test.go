package main

import (
	"os"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
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

type Bus struct {
	Cartridge Cartridge
}

func (b *Bus) Read(addr uint16) uint8 {
	if addr < 0x8000 {
		return b.Cartridge.Read(addr)
	}

	panic("not implemented")
}

func (b *Bus) Write(addr uint16, value uint8) {
	if addr < 0x8000 {
		b.Cartridge.Write(addr, value)
		return
	}

	panic("not implemented")
}

type Cpu struct {
	Registers Registers
	Bus       Bus
}

type Cartridge struct {
	Content []uint8
}

func CartridgeFromFile(filename string) (*Cartridge, error) {
	content, err := os.ReadFile(filename)

	if err != nil {
		return nil, err
	}

	return &Cartridge{
		Content: content,
	}, nil
}

func (c *Cartridge) Read(addr uint16) uint8 {
	return c.Content[addr]
}

func (c *Cartridge) Write(addr uint16, value uint8) {
	panic("not implemented")
	// c.Content[addr] = value
}

func (c *Cartridge) Title() string {
	titleBytes := c.Content[0x0134 : 0x0143+1]
	title := strings.Trim(string(titleBytes), "\x00")
	return title
}

func (c *Cartridge) ValidateHeaderChecksum() bool {
	checksum := uint8(0)
	for addr := 0x0134; addr <= 0x014C; addr++ {
		checksum -= c.Content[addr] + 1
	}

	return checksum == c.Content[0x014D]
}

func (c *Cartridge) Entry() []uint8 {
	return c.Content[0x0100 : 0x0103+1]
}

func (c *Cartridge) OldLicenseeCode() uint8 {
	return c.Content[0x014B]
}

func (c *Cartridge) SizeKB() uint32 {
	return 32 * (1 << c.Content[0x0148])
}

func (c *Cartridge) NintendoLogo() []uint8 {
	return c.Content[0x0104 : 0x0133+1]
}

func (c *Cartridge) DestinationCode() uint8 {
	return c.Content[0x014A]
}

func (c *Cartridge) Version() uint8 {
	return c.Content[0x014C]
}

func (c *Cartridge) Type() uint8 {
	return c.Content[0x0147]
}

func (c *Cartridge) SgbFlag() uint8 {
	return c.Content[0x0146]
}

const (
	DEST_JP_AND_OVERSEAS uint8 = 0x00
	LICENSEE_NINTENDO    uint8 = 0x01
	CART_TYPE_ROM_ONLY   uint8 = 0x00
)

var NINTENDO_LOGO = []uint8{
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73,
	0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F,
	0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD,
	0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}

const (
	NOP = 0x00
	JP  = 0xC3
)

func TestCartridgeInfoTetris(t *testing.T) {
	c, err := CartridgeFromFile("tetris.gb")
	assert.NoError(t, err)

	expectedEntry := []uint8{NOP, JP, 0x50, 0x01}
	assert.EqualValues(t, expectedEntry, c.Entry())
	assert.EqualValues(t, NINTENDO_LOGO, c.NintendoLogo())
	assert.Equal(t, DEST_JP_AND_OVERSEAS, c.DestinationCode())
	assert.Equal(t, LICENSEE_NINTENDO, c.OldLicenseeCode())
	assert.Equal(t, uint8(0x01), c.Version())
	assert.Equal(t, uint32(32), c.SizeKB())
	assert.Equal(t, CART_TYPE_ROM_ONLY, c.Type())
	assert.Equal(t, uint8(0x00), c.SgbFlag())
	assert.True(t, c.ValidateHeaderChecksum())
	assert.Equal(t, "TETRIS", c.Title())
}
