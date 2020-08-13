package main

type CPU struct {
	bus Bus
	A byte
	B byte
	C byte
	D byte
	E byte
	F byte
	H byte
	L byte
	PC uint16
	SP uint16
}

func (cpu *CPU) initRegisters() {
	cpu.A = 0x01
	cpu.F = 0xB0
	cpu.B = 0x00
	cpu.C = 0x13
	cpu.D = 0x00
	cpu.E = 0xD8
	cpu.H = 0x01
	cpu.L = 0x4D
	cpu.SP = 0xFFFE
	cpu.PC = 0x100
}

func (cpu *CPU) initIO() {
	cpu.bus.timer.TIMA, cpu.bus.timer.TMA, cpu.bus.timer.TAC = 0x00, 0x00, 0x00
	cpu.bus.apu.NR10 = 0x80
	cpu.bus.apu.NR11 = 0xBF
	cpu.bus.apu.NR12 = 0xF3
	cpu.bus.apu.NR14 = 0xBF
	cpu.bus.apu.NR21 = 0x3F
	cpu.bus.apu.NR22 = 0x00
	cpu.bus.apu.NR24 = 0xBF
	cpu.bus.apu.NR30 = 0x7F
	cpu.bus.apu.NR31 = 0xFF
	cpu.bus.apu.NR32 = 0x9F
	cpu.bus.apu.NR33 = 0xBF
	cpu.bus.apu.NR41 = 0xFF
	cpu.bus.apu.NR42, cpu.bus.apu.NR43 = 0x00, 0x00
	cpu.bus.apu.NR44 = 0xBF
	cpu.bus.apu.NR50 = 0x77
	cpu.bus.apu.NR51 = 0xF3
	cpu.bus.apu.NR52 = 0xF1
	cpu.bus.ppu.LCDC = 0x91
	cpu.bus.ppu.SCY, cpu.bus.ppu.SCX, cpu.bus.ppu.LYC = 0x00, 0x00, 0x00
	cpu.bus.ppu.BGP = 0xFC
	cpu.bus.ppu.OBP0, cpu.bus.ppu.OBP1 = 0xFF, 0xFF
	cpu.bus.ppu.WY, cpu.bus.ppu.WX, cpu.bus.interrupt.IE = 0x00, 0x00, 0x00
}