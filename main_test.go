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

// expectedEntry := []uint8{INS_NOP, INS_JP, 0x50, 0x01}
// assert.EqualValues(t, NINTENDO_LOGO, bus.Cartridge.Content[ADDR_NINTENDO_LOGO_START:ADDR_NINTENDO_LOGO_END+1])
// assert.Equal(t, DEST_CODE_JAPAN_AND_OVERSEAS, bus.Cartridge.Read(ADDR_DESTINATION_CODE))
// assert.Equal(t, LICENSEE_CODE_NINTENDO, bus.Cartridge.Read(ADDR_LICENSEE_CODE_OLD))
// assert.Equal(t, uint8(0x01), bus.Cartridge.Read(ADDR_ROM_VERSION_NUM))
// assert.Equal(t, uint(32), GetRomSizeKB(&bus))
// assert.Equal(t, CART_TYPE_ROM_ONLY, bus.Cartridge.Read(ADDR_CART_TYPE))
// assert.Equal(t, uint8(0x00), bus.Cartridge.Read(ADDR_SGB_FLAG))

// type RomHeader struct {
// 	Entry           []uint8 // 4 bytes
// 	NintendoLogo    []uint8 // 48 bytes
// 	Title           []uint8 // 12 bytes
// 	LicenseeCodeNew uint16
// 	SgbFlag         uint8
// 	RomType         uint8
// 	RomSize         uint8
// 	RamSize         uint8
// 	DestinationCode uint8
// 	LicenseeCodeOld uint16
// 	Version         uint8
// 	Checksum        uint8
// 	GlobalChecksum  uint16
// }

// func RomHeaderFromBytes(b []uint8) RomHeader {
// 	header := RomHeader{}

// 	header.Entry = make([]uint8, 4)
// 	copy(header.Entry, b[:4])

// 	header.NintendoLogo = make([]uint8, 48)
// 	header.Title = make([]uint8, 12)

// 	return header
// 	// header.
// }

const (
	// ADDR_DESTINATION_CODE uint16 = 0x014A
	DEST_JP_AND_OVERSEAS uint8 = 0x00
	// DEST_CODE_OVERSEAS_ONLY uint8  = 0x01

	// ADDR_LICENSEE_CODE_OLD uint16 = 0x014B
	LICENSEE_NINTENDO uint8 = 0x01

	// ADDR_CART_TYPE     uint16 = 0x0147
	CART_TYPE_ROM_ONLY uint8 = 0x00

	// ADDR_ROM_VERSION_NUM uint16 = 0x014C

	// addrChecksum    uint16 = 0x014D
	// start           uint16 = 0x0134
	// ADDR_HEADER_END uint16 = 0x014C

	// ADDR_ENTRY_START uint16 = 0x0100
	// ADDR_ENTRY_END   uint16 = 0x0103

	// ADDR_NINTENDO_LOGO_START uint16 = 0x0104
	// ADDR_NINTENDO_LOGO_END   uint16 = 0x0133

	// ADDR_SGB_FLAG uint16 = 0x0146
)

var NINTENDO_LOGO = []uint8{
	0xCE, 0xED, 0x66, 0x66, 0xCC, 0x0D, 0x00, 0x0B, 0x03, 0x73,
	0x00, 0x83, 0x00, 0x0C, 0x00, 0x0D, 0x00, 0x08, 0x11, 0x1F,
	0x88, 0x89, 0x00, 0x0E, 0xDC, 0xCC, 0x6E, 0xE6, 0xDD, 0xDD,
	0xD9, 0x99, 0xBB, 0xBB, 0x67, 0x63, 0x6E, 0x0E, 0xEC, 0xCC,
	0xDD, 0xDC, 0x99, 0x9F, 0xBB, 0xB9, 0x33, 0x3E,
}

// type Instruction uint8

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
