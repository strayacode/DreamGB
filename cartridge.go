package main

import (
	// "io/ioutil"
)

type Cartridge struct {
	ROM [0x4000]byte
	rombank ROMBank
	rambank RAMBank
}

type ROMBank struct {
	bankptr uint16 // 0-512
	bank [512][0x4000]byte
}

type RAMBank struct {
	bankptr byte // 0-16
	bank [16][0x2000]byte
}

