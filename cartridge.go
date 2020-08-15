package main

import (
	"io/ioutil"
	"os"
	"fmt"
)

type Cartridge struct {
	rombank ROMBank
	rambank RAMBank
	header Header
	romCache [0x100]byte
}

type Header struct {
	title string
	cartridgeType byte
	ROMSize byte
	RAMSize byte
}

type ROMBank struct {
	bankptr uint16 // 0-512
	bank [512][0x4000]byte
}

type RAMBank struct {
	bankptr byte // 0-16
	bank [16][0x2000]byte
}

func (cartridge *Cartridge) loadBootROM() {
	_, err := os.Stat("bios.rom")
	if os.IsNotExist(err) {
		fmt.Println("no bios file detected!")
		os.Exit(0)
	}
	file, err := ioutil.ReadFile("bios.rom")
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < len(file); i++ {
		cartridge.rombank.bank[0][i] = file[i]
	}
}


// loads cartridge into memory as well as the header
func (cartridge *Cartridge) loadCartridge() {
	_, err := os.Stat(config.rompath)
	if os.IsNotExist(err) {
		fmt.Println(".gb file doesn't exist!")
		os.Exit(0)
	}
	file, err := ioutil.ReadFile(config.rompath)
	if err != nil {
		fmt.Println(err)
	}

	for i := 0; i < 256; i++ {
		cartridge.romCache[i] = file[i]
	}

	for i := 0; i < 512; i++ {
		for j := 0; j < 0x4000; j++ {
			if (i * 0x4000) + j < len(file) {
				cartridge.rombank.bank[i][j] = file[(i * 0x4000) + j]
			}
		}
	}
}

func (cartridge *Cartridge) unmapBootROM() {
	for i := 0; i < 256; i++ {
		cartridge.rombank.bank[0][i] = cartridge.romCache[i]
	}
}