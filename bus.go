package main

import (
	"fmt"
)

type Bus struct {
	WRAM [0x2000]byte
	HRAM [0x80]byte
	ppu PPU
	timer Timer
	interrupt Interrupt
	apu APU
	cartridge Cartridge
}

func (bus *Bus) read(addr uint16) byte {
	cpu.tick(4)
	switch {
	case addr >= 0x0000 && addr <= 0x3FFF:
		return bus.cartridge.rombank.bank[0][addr]
	case addr >= 0x4000 && addr <= 0x7FFF:
		return bus.cartridge.rombank.bank[bus.cartridge.rombank.bankptr][addr - 0x4000]
	case addr >= 0x8000 && addr <= 0x9FFF:
		return bus.ppu.VRAM[addr - 0x8000]
	case addr >= 0xA000 && addr <= 0xBFFF:
		return bus.cartridge.rambank.bank[bus.cartridge.rambank.bankptr][addr - 0xA000]
	case addr >= 0xC000 && addr <= 0xCFFF:
		return bus.WRAM[addr - 0xC000]
	case addr >= 0xE000 && addr <= 0xFDFF:
		return bus.read(addr - 0x2000)
	case addr >= 0xFE00 && addr <= 0xFE9F:
		return bus.ppu.OAM[addr - 0xFE00]
	// TODO: check behaviour for unusable area later
	case addr >= 0xFF00 && addr <= 0xFF7F:
		return bus.readIO(addr)
	case addr >= 0xFF80 && addr <= 0xFFFE:
		return bus.HRAM[addr - 0xFF80]
	case addr == 0xFFFF:
		return bus.interrupt.IE
	default:
		fmt.Println("unimplemented memory read!")
		return 0
	}
	
}

func (bus *Bus) readIO(addr uint16) byte {
	switch {
	default:
		return 0
		fmt.Println(addr)
	}
	return 0
}