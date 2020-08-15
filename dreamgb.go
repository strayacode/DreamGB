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
		// process: load cartridge, load bootrom into first 0xFF bytes and then once write occurs to 0xFF50 unmap the bootrom
		cpu.bus.cartridge.loadCartridge()
		cpu.bus.cartridge.loadBootROM()
		// TODO load cartridge
	} else {
		// skip bootrom
		cpu.initRegisters()
		cpu.initIO()
		cpu.bus.cartridge.loadCartridge()
		// TODO load cartridge
	}
	window.init()
	window.loop()
}