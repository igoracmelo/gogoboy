package main

type Bus struct {
	Cartridge *Cartridge
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
