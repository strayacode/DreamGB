package main

import (
	// "fmt"
)

type PPU struct {
	VRAM [0x2000]byte
	OAM [0xA0]byte
	LCDC byte
	LCDCSTAT byte
	SCY byte
	SCX byte
	LYC byte
	LY byte
	BGP byte
	OBP0 byte
	OBP1 byte
	WX byte
	WY byte
	cycles int // counts the cycles per line
	DMA byte
	// add accurate cpu oam and vram accesses
}
// how it works: each line goes for 456 t cycles for the 144 scanlines and vblank takes 4560
func (ppu *PPU) update() {
	// fmt.Println(ppu.LCDCSTAT & 0x03)
	switch ppu.LCDCSTAT & 0x03 {
	case 0:
		// HBlank. always occurs once 456 t cycles have passed on a scanline
		if ppu.cycles >= 456 {
			ppu.cycles = 0
			ppu.LY++
			if ppu.LY >= 144 {
				// enter mode 1
				ppu.LCDCSTAT |= 0x02
				ppu.LCDCSTAT &= 0xFD
			} else {
				// enter 2 (new scanline)
				ppu.LCDCSTAT |= 0x02
				ppu.LCDCSTAT &= 0xFE
			}
		}
	case 1:
		if ppu.cycles >= 456 {
			// check LYC interrupt later but increment LY
			ppu.LY++
			if ppu.LY > 153 {
				ppu.LCDCSTAT |= 0x02
				ppu.LCDCSTAT &= 0xFE
				ppu.LY = 0
			}
		}
	case 2:
		// OAM search. always 80 t cycles
		// cycles will be 0 once OAM search is entered as it is the first mode on a scanline
		if ppu.cycles >= 80 {
			ppu.LCDCSTAT |= 0x03 // switch to mode 3
		}
	case 3:
		// read oam and vram and generate framebuffer for that scanline
		// total cycles can range from 168 to 291 t cycles
		// TODO add proper ppu timing
		if ppu.cycles >= 252 {
			ppu.LCDCSTAT &= 0xFC // switch to mode 0
		}
	}
}