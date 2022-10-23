package main

import (
	"os"
	"strings"
)

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
