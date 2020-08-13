package main

type Bus struct {
	ppu PPU
	timer Timer
	interrupt Interrupt
	apu APU
	cartridge Cartridge
}

func (bus *Bus) read() {

}