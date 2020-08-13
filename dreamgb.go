package main

import (
	"fmt"
)

func main() {
	window := Window{}
	config := Config{}
	cpu := CPU{}
	config.fetchFlags()
	// load bootrom
	if config.bootrom {
		fmt.Println("use bootrom")
		cpu.bus.cartridge.loadBootROM()
		fmt.Println(cpu.bus.cartridge.ROM)
	} else {
		// skip bootrom
		cpu.initRegisters()
		cpu.initIO()
		fmt.Println("skip bootrom")
	}
	window.init()
	window.loop()
}