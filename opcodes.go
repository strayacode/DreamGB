package main

import (
	"fmt"
	"strconv"
	"os"
)


type Opcode struct {
	description string
	exec func(cpu *CPU)
}

func (cpu *CPU) NOP() {

}

func (cpu *CPU) LDSP_u16() {
	cpu.SP = cpu.bus.read16(cpu.PC)
	cpu.PC += 2
}

func (cpu *CPU) LDr16_u16(r1 *byte, r2 *byte) {
	*r2 = cpu.bus.read(cpu.PC)
	*r1 = cpu.bus.read(cpu.PC + 1)
	cpu.PC += 2
}

func (cpu *CPU) LDDHL_A() {
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), cpu.A)
	cpu.L--
	if cpu.L == 0xFF {
		cpu.H--
	}
}

func (cpu *CPU) LDHL_A() {
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), cpu.A)
}

func (cpu *CPU) XORA_r8(r8 byte) {
	cpu.A ^= r8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
}

func (cpu *CPU) CB() {
	// get the cb opcode
	cpu.opcode = cpu.bus.read(cpu.PC)
	cbopcodes[cpu.opcode].exec(cpu)
	cpu.PC++ // since none of the cb opcodes change PC we increment after the instruction is done
}

func (cpu *CPU) BITu3_r8(bit int, r8 byte) {
	cpu.ZFlag((r8 & (1 << bit)) == 0)
	cpu.NFlag(false)
	cpu.HFlag(true)
}

func (cpu *CPU) JRNZ_i8() {
	offset := int8(cpu.bus.read(cpu.PC) + 1)
	// branch NZ
	if !cpu.getZFlag() {
		cpu.PC += uint16(offset)
		cpu.tick(4)
		
	} else {
		cpu.PC++
	}
}

func (cpu *CPU) LDr8_u8(r8 *byte) {
	*r8 = cpu.bus.read(cpu.PC)
	cpu.PC++
}

func (cpu *CPU) LDr8_r8(r1 *byte, r2 byte) {
	*r1 = r2
}

func (cpu *CPU) LDA_mem(addr uint16) {
	cpu.bus.write(addr, cpu.A)
}

func (cpu *CPU) CALL_u16() {
}

func (cpu *CPU) INC_r8(r8 *byte) {
	cpu.HFlag(((cpu.C & 0xF) + (1 & 0xF)) & 0x10 == 0x10)
	*r8++
	cpu.ZFlag(cpu.C == 0)
	cpu.NFlag(false)
}

func (cpu *CPU) UND_OP() {
	fmt.Println("undefined opcode! 0x" + strconv.FormatUint(uint64(cpu.opcode), 16))
	os.Exit(0)
}

func (cpu *CPU) UND_CBOP() {
	fmt.Println("undefined cb opcode! 0x" + strconv.FormatUint(uint64(cpu.opcode), 16))
	os.Exit(0)
}

var opcodes = [256]Opcode {
	Opcode{"NOP", func(cpu *CPU) { cpu.NOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC C", func(cpu *CPU) { cpu.INC_r8(&cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD C, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD DE, u16", func(cpu *CPU) { cpu.LDr16_u16(&cpu.D, &cpu.E) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD A, (DE)", func(cpu *CPU) { cpu.LDr8_r8(&cpu.A, cpu.bus.read(uint16(cpu.D) << 8 | uint16(cpu.E))) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"JR NZ, i8", func(cpu *CPU) { cpu.JRNZ_i8() }}, Opcode{"LD HL, u16", func(cpu *CPU) { cpu.LDr16_u16(&cpu.H, &cpu.L) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD SP, u16", func(cpu *CPU) { cpu.LDSP_u16() }}, Opcode{"LD (HL-), A", func(cpu *CPU) { cpu.LDDHL_A() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD A, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD B, A", func(cpu *CPU) { cpu.LDr8_r8(&cpu.B, cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD (HL), A", func(cpu *CPU) { cpu.LDHL_A() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"XOR A, A", func(cpu *CPU) { cpu.XORA_r8(cpu.A) }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CB", func(cpu *CPU) { cpu.CB() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"LD (FF00+u8), A", func(cpu *CPU) { cpu.LDA_mem(0xFF00 + uint16(cpu.fetchOpcode())) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD (FF00 + C), A", func(cpu *CPU) { cpu.LDA_mem(0xFF00 + uint16(cpu.C)) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
}

var cbopcodes = [256]Opcode {
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"BIT 7, H", func(cpu *CPU) { cpu.BITu3_r8(7, cpu.H) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }},
}