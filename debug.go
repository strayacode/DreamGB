package main

import (
	"fmt"
	"strconv"
)


func (cpu *CPU) debugCPU() {
	fmt.Println("A: ", strconv.FormatUint(uint64(cpu.A), 16), "F ", strconv.FormatUint(uint64(cpu.F), 16))
	fmt.Println("B: ", strconv.FormatUint(uint64(cpu.B), 16), "C: ", strconv.FormatUint(uint64(cpu.C), 16))
	fmt.Println("D: ", strconv.FormatUint(uint64(cpu.D), 16), "E: ", strconv.FormatUint(uint64(cpu.E), 16))
	fmt.Println("H: ", strconv.FormatUint(uint64(cpu.H), 16), "L: ", strconv.FormatUint(uint64(cpu.L), 16))
	fmt.Println("PC: ", strconv.FormatUint(uint64(cpu.PC), 16), "SP: ", strconv.FormatUint(uint64(cpu.SP), 16))
	fmt.Println("Opcode: 0x" + strconv.FormatUint(uint64(cpu.opcode), 16))
	fmt.Println("Cycles: " + strconv.Itoa(cpu.cycles))
}

func (cpu *CPU) debugPPU() {
	fmt.Println("LCDC: 0x" + strconv.FormatUint(uint64(cpu.bus.ppu.LCDC), 16))
	fmt.Println("LCDCSTAT: 0x" + strconv.FormatUint(uint64(cpu.bus.ppu.LCDCSTAT), 16))
	fmt.Println("LY: 0x" + strconv.FormatUint(uint64(cpu.bus.ppu.LY), 16))
	fmt.Println("PPU cycles: " + strconv.Itoa(cpu.bus.ppu.cycles))
	fmt.Println(cpu.bus.ppu.BGFIFO)
}