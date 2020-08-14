package main

type PPU struct {
	VRAM [0x2000]byte
	OAM [0xA0]byte
	LCDC byte
	SCY byte
	SCX byte
	LYC byte
	BGP byte
	OBP0 byte
	OBP1 byte
	WX byte
	WY byte
	cycles byte
}