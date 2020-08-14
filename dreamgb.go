package main

var (
	cpu = CPU{}
	window = Window{}
	config = Config{}
)

func main() {
	config.fetchFlags()
	// load bootrom
	if config.bootrom {
		cpu.bus.cartridge.loadBootROM()
		// TODO load cartridge
	} else {
		// skip bootrom
		cpu.initRegisters()
		cpu.initIO()
		// TODO load cartridge
	}
	window.init()
	window.loop()
}