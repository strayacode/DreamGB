package main

import (
	"math"
	"fmt"
)

var (
	colours = [4]RGB{RGB{202, 220, 159}, RGB{15, 56, 15}, RGB{48, 98, 48}, RGB{139, 172, 15}}
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
	LX byte
	BGP byte
	OBP0 byte
	OBP1 byte
	WX byte
	WY byte
	cycles int // counts the cycles per line
	DMA byte
	cpuOAMAccess bool
	cpuVRAMAccess bool
	frameBuffer [144 * 166 * 3]byte
	BGFIFO []Pixel
	SPRFIFO []Pixel
	// add accurate cpu oam and vram accesses
}

type Pixel struct {
	colour byte // 0-3
	palette byte // 0-7 on cgb and only applies to sprites on dmg
	spritePriority byte // only on cgb 
	backgroundPriority byte 
}

type RGB struct {
	R byte
	G byte
	B byte
}

// how it works: each line goes for 456 t cycles for the 144 scanlines and vblank takes 4560
func (ppu *PPU) update() {
	// fmt.Println(ppu.LCDCSTAT & 0x03)
	switch ppu.LCDCSTAT & 0x03 {
	case 0:
		// HBlank. always occurs once 456 t cycles have passed on a scanline
		if ppu.cycles >= 456 {
			ppu.cycles = 0
			ppu.LX = 0
			ppu.LY++
			if ppu.LY >= 144 {
				// fmt.Println("vblank")
				// enter mode 1
				ppu.LCDCSTAT |= 0x01
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
			ppu.cycles = 0
			ppu.LY++
			if ppu.LY > 153 {
				ppu.LCDCSTAT |= 0x02
				ppu.LCDCSTAT &= 0xFE
				ppu.LX = 0
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
		// when fifo is manipulated
		// read oam and vram and generate framebuffer for that scanline
		// total cycles can range from 168 to 291 t cycles
		// TODO add proper ppu timing
		if ppu.cycles >= 252 {
			ppu.drawScanline()
			ppu.LCDCSTAT &= 0xFC // switch to mode 0
		}
	}
}

func (ppu *PPU) drawScanline() {
	for i := 0; i < 20; i++ {
		ppu.fetchRow()
		if ppu.LX >= 160 {
			ppu.BGFIFO = nil
			break
		}
	}
}

func (ppu *PPU) fetchRow() {
	// start of scanline
	if len(ppu.BGFIFO) == 0 {
		// Read Tile #. 2 cycles
		offsetX := math.Floor(float64(ppu.LX + ppu.SCX) / 8)
		offsetY := math.Floor(float64(ppu.LY + ppu.SCY) / 8)
		tileAddr := ppu.getBGMapAddr() + uint16(offsetY * 32) + uint16(offsetX) // returns the address of the tile number
		tileOffset := uint16((ppu.SCY + ppu.LY) % 8) * 2 // get the 2 bytes corresponding to the correct row in the tile
		// next read data 0
		tile0 := ppu.VRAM[ppu.getBGStartAddr() - 0x8000 + uint16(ppu.VRAM[tileAddr - 0x8000]) * 16 + tileOffset]
		// next read data 1
		tile1 := ppu.VRAM[ppu.getBGStartAddr() - 0x8000 + uint16(ppu.VRAM[tileAddr - 0x8000]) * 16 + tileOffset + 1]

		// combine tile0 and tile1 to push 8 pixels to fifo but only if len(fifo) <= 8
		// just immediately place the 8 pixels in fifo if cycles <= 86
		for i := 0; i < 8; i++ {
			selection := ((((1 << (7 - i)) & tile0)) >> (7 - i)) << 1 | ((1 << (7 - i)) & tile1) >> (7 - i)
			colour := (ppu.BGP & (0x03 << (selection * 2))) >> (selection * 2)
			// for now ignore colour
			ppu.BGFIFO = append(ppu.BGFIFO, Pixel{colour, 0, 0, 0}) // ignore last 3 params for now
			// first 8 pixels appended
		}
	} else {
		ppu.pushPixel()
		ppu.pushPixel()
		// Read Tile #. 2 cycles
		offsetX := math.Floor(float64(ppu.LX + ppu.SCX) / 8)
		offsetY := math.Floor(float64(ppu.LY + ppu.SCY) / 8)
		tileAddr := ppu.getBGMapAddr() + uint16(offsetY * 32) + uint16(offsetX) // returns the address of the tile number
		tileOffset := uint16((ppu.SCY + ppu.LY) % 8) * 2 // get the 2 bytes corresponding to the correct row in the tile
		ppu.pushPixel()
		ppu.pushPixel()
		// next read data 0
		tile0 := ppu.VRAM[ppu.getBGStartAddr() - 0x8000 + uint16(ppu.VRAM[tileAddr - 0x8000]) * 16 + tileOffset]
		ppu.pushPixel()
		ppu.pushPixel()
		// next read data 1
		tile1 := ppu.VRAM[ppu.getBGStartAddr() - 0x8000 + uint16(ppu.VRAM[tileAddr - 0x8000]) * 16 + tileOffset + 1]
		ppu.pushPixel()
		ppu.pushPixel()
		// now fifo should be only 8 pixels large
		for i := 0; i < 8; i++ {
			selection := ((((1 << (7 - i)) & tile0)) >> (7 - i)) << 1 | ((1 << (7 - i)) & tile1) >> (7 - i)
			colour := (ppu.BGP & (0x03 << (selection * 2))) >> (selection * 2)
			// for now ignore colour
			ppu.BGFIFO = append(ppu.BGFIFO, Pixel{colour, 0, 0, 0}) // ignore last 3 params for now
			// first 8 pixels appended
		}
	}
	
}

func (ppu *PPU) pushPixel() {
	fmt.Println(ppu.LY, ppu.LX, colours[ppu.BGFIFO[0].colour])
	// push pixel to framebuffer
	ppu.frameBuffer[(ppu.LY * 160 + ppu.LX) * 3] = colours[ppu.BGFIFO[0].colour].R
	ppu.frameBuffer[(ppu.LY * 160 + ppu.LX) * 3 + 1] = colours[ppu.BGFIFO[0].colour].G
	ppu.frameBuffer[(ppu.LY * 160 + ppu.LX) * 3 + 2] = colours[ppu.BGFIFO[0].colour].B
	fmt.Println(ppu.frameBuffer)
	// increment ppu.LX
	ppu.LX++
	
	// remove first pixel from fifo
	ppu.BGFIFO = ppu.BGFIFO[1:]
}

func (ppu *PPU) getBGStartAddr() uint16 {
	if ppu.LCDC & 0x10 >> 4 == 1 {
		return 0x8000
	}
	return 0x9000
}

func (ppu *PPU) getBGMapAddr() uint16 {
	// optimize later
	if ppu.LCDC & 0x08 >> 3 == 1 {
		return 0x9C00
	}
	return 0x9800
}

