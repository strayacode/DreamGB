package main

import (
	"fmt"
	"strconv"
	"os"
)

type Bus struct {
	WRAM [0x2000]byte
	HRAM [0x80]byte
	ppu PPU
	timer Timer
	interrupt Interrupt
	apu APU
	cartridge Cartridge
	input Input
	serial Serial
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
	case addr >= 0xC000 && addr <= 0xDFFF:
		return bus.WRAM[addr - 0xC000]
	case addr >= 0xE000 && addr <= 0xFDFF:
		return bus.WRAM[addr - 0xE000]
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
		fmt.Println("unimplemented memory read! 0x" + strconv.FormatUint(uint64(addr), 16))
		// os.Exit(0)
		return 0
	}
}

func (bus *Bus) read16(addr uint16) uint16 {
	// remember gameboy is little-endian: if given data 0xFFFE then 0xFE is read/stored first then 0xFF
	return uint16(bus.read(addr + 1)) << 8 | uint16(bus.read(addr))
}

func (bus *Bus) write(addr uint16, data byte) {
	cpu.tick(4)
	switch {
	// no writes allowed to rom mem region
	case addr >= 0x8000 && addr <= 0x9FFF:
		bus.ppu.VRAM[addr - 0x8000] = data
	case addr >= 0xA000 && addr <= 0xBFFF:
		bus.cartridge.rambank.bank[bus.cartridge.rambank.bankptr][addr - 0xA000] = data
	case addr >= 0xC000 && addr <= 0xDFFF:
		bus.WRAM[addr - 0xC000] = data
	case addr >= 0xE000 && addr <= 0xFDFF:
		bus.WRAM[addr - 0xE000] = data
	case addr >= 0xFE00 && addr <= 0xFE9F:
		bus.ppu.OAM[addr - 0xFE00] = data
	case addr >= 0xFF00 && addr <= 0xFF7F:
		bus.writeIO(addr, data)
	case addr >= 0xFF80 && addr <= 0xFFFE:
		bus.HRAM[addr - 0xFF80] = data
	case addr == 0xFFFF:
		bus.interrupt.IE = data
	default:
		fmt.Println("unimplemented write! 0x" + strconv.FormatUint(uint64(addr), 16))
		os.Exit(0)
	}
}

func (bus *Bus) readIO(addr uint16) byte {
	switch addr {
	case 0xFF00:
		return bus.input.P1
	case 0xFF01:
		fmt.Println(string(bus.serial.SB))
		return bus.serial.SB
	case 0xFF02:
		return bus.serial.SC
	case 0xFF04:
		return bus.timer.DIV
	case 0xFF05:
		return bus.timer.TIMA
	case 0xFF0F:
		return bus.interrupt.IF
	case 0xFF10:
		return bus.apu.NR10 
	case 0xFF11:
		return bus.apu.NR11 
	case 0xFF12:
		return bus.apu.NR12 
	case 0xFF13:
		return bus.apu.NR13 
	case 0xFF14:
		return bus.apu.NR14 
	case 0xFF16:
		return bus.apu.NR21 
	case 0xFF17:
		return bus.apu.NR22 
	case 0xFF19:
		return bus.apu.NR24 
	case 0xFF1A:
		return bus.apu.NR30 
	case 0xFF1B:
		return bus.apu.NR31 
	case 0xFF1C:
		return bus.apu.NR32 
	case 0xFF1D:
		return bus.apu.NR33
	case 0xFF1E:
		return bus.apu.NR34 
	case 0xFF20:
		return bus.apu.NR41 
	case 0xFF21:
		return bus.apu.NR42 
	case 0xFF22:
		return bus.apu.NR43 
	case 0xFF23:
		return bus.apu.NR44 
	case 0xFF24:
		return bus.apu.NR50 
	case 0xFF25:
		return bus.apu.NR51 
	case 0xFF26:
		return bus.apu.NR52 
	case 0xFF40:
		return bus.ppu.LCDC
	case 0xFF41:
		return bus.ppu.LCDCSTAT
	case 0xFF42:
		return bus.ppu.SCY
	case 0xFF43:
		return bus.ppu.SCX
	case 0xFF44:
		return bus.ppu.LY
	case 0xFF45:
		return bus.ppu.LYC
	case 0xFF46:
		return bus.ppu.DMA
	case 0xFF47:
		return bus.ppu.BGP
	case 0xFF48:
		return bus.ppu.OBP0
	case 0xFF49:
		return bus.ppu.OBP1
	case 0xFF4A:
		return bus.ppu.WY
	case 0xFF4B:
		return bus.ppu.WX
	default:
		fmt.Println(addr, "IO read not implemented yet!")
		// os.Exit(3)
		return 0xFF
	}
	return 0xFF
}

func (bus *Bus) writeIO(addr uint16, data byte) {
	switch addr {
	case 0xFF00:
		bus.input.P1 = data
	case 0xFF01:
		bus.serial.SB = data
		fmt.Println(string(data))
	case 0xFF02:
		bus.serial.SC = data
	case 0xFF04:
		bus.timer.DIV = data
	case 0xFF05:
		bus.timer.TIMA = data
	case 0xFF06:
		bus.timer.TMA = data
	case 0xFF07:
		bus.timer.TAC = data
	case 0xFF10:
		bus.apu.NR10 = data
	case 0xFF11:
		bus.apu.NR11 = data
	case 0xFF12:
		bus.apu.NR12 = data
	case 0xFF14:
		bus.apu.NR14 = data
	case 0xFF16:
		bus.apu.NR21 = data
	case 0xFF0F:
		bus.interrupt.IF = data
	case 0xFF24:
		bus.apu.NR50 = data
	case 0xFF25:
		bus.apu.NR51 = data
	case 0xFF26:
		bus.apu.NR52 = data
	case 0xFF30, 0xFF31, 0xFF32, 0xFF33, 0xFF34, 0xFF35, 0xFF36, 0xFF37, 0xFF38, 0xFF39, 0xFF3A, 0xFF3B, 0xFF3C, 0xFF3D, 0xFF3E, 0xFF3F:
		bus.apu.WAVEPATTERN[addr - 0xFF30] = data
	case 0xFF40:
		bus.ppu.LCDC = data
	case 0xFF41:
		bus.ppu.LCDCSTAT = data
		// todo: block writes to first 2 bits
	case 0xFF42:
		bus.ppu.SCY = data
	case 0xFF43:
		bus.ppu.SCX = data
	case 0xFF47:
		bus.ppu.BGP = data
	case 0xFF50:
		bus.cartridge.unmapBootROM()
	default:
		fmt.Println("unimplemented write! 0x" + strconv.FormatUint(uint64(addr), 16))
	}
}