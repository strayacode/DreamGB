package main

import (
	"fmt"
)

var (
	cpu = CPU{}
	window = Window{}
	config = Config{}
)

func main() {
	config.fetchFlags()
	// load bootrom
	if config.bootrom {
		fmt.Println("use bootrom")
		cpu.bus.cartridge.loadBootROM()
		fmt.Println(cpu.bus.cartridge.rombank.bank[0])
	} else {
		// skip bootrom
		cpu.initRegisters()
		cpu.initIO()
		fmt.Println("skip bootrom")
	}
	window.init()
	window.loop()
}