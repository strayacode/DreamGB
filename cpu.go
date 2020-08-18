package main

// import (
// 	"fmt"
// )

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
	cycles int // amount of cycles elapsed per frame
	opcode byte
	halt bool
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

// executes an opcode and increments the cycles (useful for stepmode)
func (cpu *CPU) step() {
	cpu.opcode = cpu.fetchOpcode()
	cpu.executeOpcode()
	// cpu.debugCPU()
}

// advances the cpu by n cycles
func (cpu *CPU) tick(cycles int) {
	cpu.cycles += 4
	cpu.bus.ppu.cycles += 4
	cpu.bus.ppu.update()
}

func (cpu *CPU) fetchOpcode() byte {
	cpu.PC++
	return cpu.bus.read(cpu.PC - 1)
}

func (cpu *CPU) executeOpcode() {
	opcodes[cpu.opcode].exec(cpu)
}

func (cpu *CPU) ZFlag(value bool) {
	if value {
		cpu.F |= (1 << 7)
	} else {
		cpu.F &= 0x7F
	}
}

func (cpu *CPU) NFlag(value bool) {
	if value {
		cpu.F |= (1 << 6)
	} else {
		cpu.F &= 0xBF
	}
}

func (cpu *CPU) HFlag(value bool) {
	if value {
		cpu.F |= (1 << 5)
	} else {
		cpu.F &= 0xDF
	}
}

func (cpu *CPU) CFlag(value bool) {
	if value {
		cpu.F |= (1 << 4)
	} else {
		cpu.F &= 0xEF
	}
}

func (cpu *CPU) getZFlag() byte {
	return (cpu.F & 0x80) >> 7
}

func (cpu *CPU) getNFlag() byte {
	return (cpu.F & 0x40) >> 6
}

func (cpu *CPU) getHFlag() byte {
	return (cpu.F & 0x20) >> 5
}

func (cpu *CPU) getCFlag() byte {
	return (cpu.F & 0x10) >> 4
}



