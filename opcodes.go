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

func (cpu *CPU) LDIHL_A() {
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), cpu.A)
	cpu.L++
	if cpu.L == 0x00 {
		cpu.H++
	}
}

func (cpu *CPU) LDA_mem(addr uint16) {
	cpu.A = cpu.bus.read(addr)
}

func (cpu *CPU) LDA_HLI() {
	cpu.A = cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.L++
	if cpu.L == 0 {
		cpu.H++
	}
}

func (cpu *CPU) LDr16_r8(r1 byte, r2 byte, r3 byte) {
	cpu.bus.write(uint16(r1) << 8 | uint16(r2), r3)
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

func (cpu *CPU) JR_i8() {
	offset := int8(cpu.bus.read(cpu.PC) + 1)
	cpu.PC += uint16(offset)
	cpu.tick(4)
}

func (cpu *CPU) JP_u16() {
	cpu.PC = cpu.bus.read16(cpu.PC)
	cpu.tick(4)
}

func (cpu *CPU) LDr8_u8(r8 *byte) {
	*r8 = cpu.bus.read(cpu.PC)
	cpu.PC++
}

func (cpu *CPU) LDr8_r8(r1 *byte, r2 byte) {
	*r1 = r2
}

func (cpu *CPU) LDmem_A(addr uint16) {
	cpu.bus.write(addr, cpu.A)
}

func (cpu *CPU) CALL_u16() {
	cpu.tick(4)
	// TODO: read u16 first
	hi := byte((cpu.PC + 2) >> 8)
	lo := byte(((cpu.PC + 2) & 0xFF))
	cpu.SP--
	cpu.bus.write(cpu.SP, hi)
	cpu.SP--
	cpu.bus.write(cpu.SP, lo)
	cpu.PC = cpu.bus.read16(cpu.PC)
}

func (cpu *CPU) PUSH_r16(r1 byte, r2 byte) {
	cpu.tick(4)
	cpu.SP--
	cpu.bus.write(cpu.SP, r1)
	cpu.SP--
	cpu.bus.write(cpu.SP, r2)
}

func (cpu *CPU) POP_r16(r1 *byte, r2 *byte) {
	*r2 = cpu.bus.read(cpu.SP)
	cpu.SP++
	*r1 = cpu.bus.read(cpu.SP)
	cpu.SP++
}

func (cpu *CPU) INC_r8(r8 *byte) {
	cpu.HFlag(((*r8 & 0xF) + (1 & 0xF)) & 0x10 == 0x10)
	*r8++
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
}

func (cpu *CPU) INC_r16(r1 *byte, r2 *byte) {
	*r2++
	if *r2 == 0 {
		*r1++
	}
	cpu.tick(4)
}

func (cpu *CPU) DEC_r8(r8 *byte) {
	cpu.HFlag(((*r8 & 0xF) - (1 & 0xF)) & 0x10 == 0x10)
	*r8--
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(true)
	
}

func (cpu *CPU) RL_r8(r8 *byte) {
	tempcarry := 0
	if cpu.getCFlag() {
		tempcarry = 1
	}
	newcarry := *r8 >> 7
	*r8 = (*r8 << 1 | byte(tempcarry))
	cpu.CFlag(newcarry == 1)
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RLA() {
	tempcarry := 0
	if cpu.getCFlag() {
		tempcarry = 1
	}
	newcarry := cpu.A >> 7
	cpu.A = (cpu.A << 1 | byte(tempcarry))
	cpu.CFlag(newcarry == 1)
	cpu.ZFlag(false)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RLCA() {
	tempcarry := cpu.A >> 7
	cpu.CFlag(tempcarry == 1)
	cpu.A = (cpu.A << 1 | tempcarry)
	cpu.ZFlag(false)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RET() {
	cpu.PC = cpu.bus.read16(cpu.SP)
	cpu.SP += 2
	cpu.tick(4)
}

func (cpu *CPU) RET_NC() {
	cpu.tick(4)
	if !cpu.getCFlag() {
		cpu.PC = cpu.bus.read16(cpu.SP)
		cpu.SP += 2
		cpu.tick(4)
	}
	
}

func (cpu *CPU) CPA_r8(value byte) {
	result := cpu.A - value
	cpu.ZFlag(result == 0)
	cpu.NFlag(true)
	cpu.HFlag(((cpu.A & 0xF) - (value & 0xF)) & 0x10 == 0x10)
	cpu.CFlag(value > cpu.A)
}

func (cpu *CPU) DI() {
	// todo interrupt delay or whatever
	cpu.bus.interrupt.IME = 0
}

func (cpu *CPU) ADDSP_i8() {
	i8 := int8(cpu.bus.read(cpu.PC))
	cpu.HFlag((((cpu.SP) & 0xF) + (uint16(i8) & 0xF)) >= 0x10)
	cpu.CFlag((uint16(i8) & 0xFF) + (cpu.SP & 0xFF) > 0xFF)
	cpu.tick(8)
	cpu.SP += uint16(i8)
	cpu.ZFlag(false)
	cpu.NFlag(false)
	cpu.PC++
}

func (cpu *CPU) ADDA_r8(r8 byte) {
	result := uint16(cpu.A) + uint16(r8)
	cpu.HFlag(((cpu.A & 0xF) + ((r8) & 0xF)) & 0x10 == 0x10)
	cpu.A += r8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.CFlag(result > 0xFF)
}

func (cpu *CPU) ANDA_r8(r8 byte) {
	cpu.A &= r8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(true)
	cpu.CFlag(false)
}

func (cpu *CPU) ADDr16_r16(r1 *byte, r2 *byte, r3 byte, r4 byte) {
	r16_1 := uint16(*r1) << 8 | uint16(*r2)
	r16_2 := uint16(r3) << 8 | uint16(r4)
	cpu.CFlag(uint32(r16_1) + uint32(r16_2) > 0xFFFF)
	cpu.HFlag(((r16_1 & 0x0FFF) + (r16_2 & 0x0FFF)) & 0x1000 == 0x1000)
	result := r16_1 + r16_2
	cpu.NFlag(false)
	*r1 = byte((result >> 8) & 0xFF)
	*r2 = byte(result)
}

func (cpu *CPU) DECmem_HL() {
	result := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.HFlag(((result & 0xF) - (1 & 0xF)) & 0x10 == 0x10)
	result--
	cpu.ZFlag(result == 0)
	cpu.NFlag(true)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), result)
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
	Opcode{"NOP", func(cpu *CPU) { cpu.NOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC BC", func(cpu *CPU) { cpu.INC_r16(&cpu.B, &cpu.C) }}, Opcode{"INC B", func(cpu *CPU) { cpu.INC_r8(&cpu.B) }}, Opcode{"DEC B", func(cpu *CPU) { cpu.DEC_r8(&cpu.B) }}, Opcode{"LD B, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.B) }}, Opcode{"RLCA", func(cpu *CPU) { cpu.RLCA() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"ADD HL, BC", func(cpu *CPU) { cpu.ADDr16_r16(&cpu.H, &cpu.L, cpu.B, cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC C", func(cpu *CPU) { cpu.INC_r8(&cpu.C) }}, Opcode{"DEC C", func(cpu *CPU) { cpu.DEC_r8(&cpu.C) }}, Opcode{"LD C, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD DE, u16", func(cpu *CPU) { cpu.LDr16_u16(&cpu.D, &cpu.E) }}, Opcode{"LD (DE), A", func(cpu *CPU) { cpu.LDr16_r8(cpu.D, cpu.E, cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC D", func(cpu *CPU) { cpu.INC_r8(&cpu.D) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"RLA", func(cpu *CPU) { cpu.RLA() }}, Opcode{"JR i8", func(cpu *CPU) { cpu.JR_i8() }}, Opcode{"ADD HL, DE", func(cpu *CPU) { cpu.ADDr16_r16(&cpu.H, &cpu.L, cpu.D, cpu.E) }}, Opcode{"LD A, (DE)", func(cpu *CPU) { cpu.LDr8_r8(&cpu.A, cpu.bus.read(uint16(cpu.D) << 8 | uint16(cpu.E))) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC E", func(cpu *CPU) { cpu.INC_r8(&cpu.E) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"JR NZ, i8", func(cpu *CPU) { cpu.JRNZ_i8() }}, Opcode{"LD HL, u16", func(cpu *CPU) { cpu.LDr16_u16(&cpu.H, &cpu.L) }}, Opcode{"LD (HL+), A", func(cpu *CPU) { cpu.LDIHL_A() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD H, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.H) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD A, (HL+)", func(cpu *CPU) { cpu.LDA_HLI() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD SP, u16", func(cpu *CPU) { cpu.LDSP_u16() }}, Opcode{"LD (HL-), A", func(cpu *CPU) { cpu.LDDHL_A() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"DEC (HL)", func(cpu *CPU) { cpu.DECmem_HL() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC A", func(cpu *CPU) { cpu.INC_r8(&cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD A, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"LD B, B", func(cpu *CPU) { cpu.LDr8_r8(&cpu.B, cpu.B) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD B, A", func(cpu *CPU) { cpu.LDr8_r8(&cpu.B, cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD C, A", func(cpu *CPU) { cpu.LDr8_r8(&cpu.C, cpu.A) }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD (HL), A", func(cpu *CPU) { cpu.LDr16_r8(cpu.H, cpu.L, cpu.A) }}, Opcode{"LD A, B", func(cpu *CPU) { cpu.LDr8_r8(&cpu.A, cpu.B) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"ADD A, B", func(cpu *CPU) { cpu.ADDA_r8(cpu.B) }}, Opcode{"ADD A, C", func(cpu *CPU) { cpu.ADDA_r8(cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"ADD A, A", func(cpu *CPU) { cpu.ADDA_r8(cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"AND A, A", func(cpu *CPU) { cpu.ANDA_r8(cpu.A) }}, Opcode{"XOR A, B", func(cpu *CPU) { cpu.XORA_r8(cpu.B) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"XOR A, A", func(cpu *CPU) { cpu.XORA_r8(cpu.A) }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CP A, H", func(cpu *CPU) { cpu.CPA_r8(cpu.H) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"POP BC", func(cpu *CPU) { cpu.POP_r16(&cpu.B, &cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"JP u16", func(cpu *CPU) { cpu.JP_u16() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"PUSH BC", func(cpu *CPU) { cpu.PUSH_r16(cpu.B, cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"RET", func(cpu *CPU) { cpu.RET() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CB", func(cpu *CPU) { cpu.CB() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CALL u16", func(cpu *CPU) { cpu.CALL_u16() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"RET NC", func(cpu *CPU) { cpu.RET_NC() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"LD (FF00+u8), A", func(cpu *CPU) { cpu.LDmem_A(0xFF00 + uint16(cpu.fetchOpcode())) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD (FF00 + C), A", func(cpu *CPU) { cpu.LDmem_A(0xFF00 + uint16(cpu.C)) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"ADD SP, i8", func(cpu *CPU) { cpu.ADDSP_i8() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD (u16), A", func(cpu *CPU) { cpu.LDmem_A(cpu.bus.read16(cpu.PC)) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"LD A, (FF00+u8)", func(cpu *CPU) { cpu.LDA_mem(0xFF00 + uint16(cpu.fetchOpcode())) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"DI", func(cpu *CPU) { cpu.DI() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CP A, u8", func(cpu *CPU) { cpu.CPA_r8(cpu.fetchOpcode()) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
}

var cbopcodes = [256]Opcode {
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"RL C", func(cpu *CPU) { cpu.RL_r8(&cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_CBOP() }}, 
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