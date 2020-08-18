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
//  8 bit arithmetic and logic instructions
func (cpu *CPU) ADCA_r8(r8 byte) {
	cpu.CFlag(uint16(cpu.A) + uint16(cpu.getCFlag()) + uint16(r8) > 0xFF)
	cpu.HFlag(((cpu.A & 0xF) + (cpu.getCFlag() & 0xF) + (r8 & 0xF)) & 0x10 == 0x10)
	cpu.A += (cpu.getCFlag() + r8)
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
}

func (cpu *CPU) ADCA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.CFlag(uint16(cpu.A) + uint16(cpu.getCFlag()) + uint16(hl) > 0xFF)
	cpu.HFlag((cpu.A & 0xF) + (cpu.getCFlag() & 0xF) + (hl & 0xF) & 0x10 == 0x10)
	cpu.A += (cpu.getCFlag() + hl)
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
}

func (cpu *CPU) ADCA_u8() {
	u8 := cpu.bus.read(cpu.PC)
	cpu.CFlag(uint16(cpu.A) + uint16(cpu.getCFlag()) + uint16(u8) > 0xFF)
	cpu.HFlag(((cpu.A & 0xF) + (cpu.getCFlag() & 0xF) + (u8 & 0xF)) & 0x10 == 0x10)
	cpu.A += (cpu.getCFlag() + u8)
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.PC++
}

func (cpu *CPU) ADDA_r8(r8 byte) {
	cpu.CFlag(uint16(cpu.A) + uint16(r8) > 0xFF)
	cpu.HFlag(((cpu.A & 0xF) + ((r8) & 0xF)) & 0x10 == 0x10)
	cpu.A += r8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
}

func (cpu *CPU) ADDA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.CFlag(uint16(cpu.A) + uint16(hl) > 0xFF)
	cpu.HFlag((cpu.A & 0xF) + (hl & 0xF) & 0x10 == 0x10)
	cpu.A += hl
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
}

func (cpu *CPU) ADDA_u8() {
	u8 := cpu.bus.read(cpu.PC)
	cpu.CFlag(uint16(cpu.A) + uint16(u8) > 0xFF)
	cpu.HFlag(((cpu.A & 0xF) + (u8 & 0xF)) & 0x10 == 0x10)
	cpu.A += u8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.PC++
}

func (cpu *CPU) ANDA_r8(r8 byte) {
	cpu.A &= r8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(true)
	cpu.CFlag(false)
}

func (cpu *CPU) ANDA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.A &= hl
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(true)
	cpu.CFlag(false)
}

func (cpu *CPU) ANDA_u8() {
	u8 := cpu.bus.read(cpu.PC)
	cpu.A &= u8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(true)
	cpu.CFlag(false)
	cpu.PC++
}

func (cpu *CPU) CPA_r8(value byte) {
	cpu.ZFlag((cpu.A - value) == 0)
	cpu.NFlag(true)
	cpu.HFlag(((cpu.A & 0xF) - (value & 0xF)) & 0x10 == 0x10)
	cpu.CFlag(value > cpu.A)
}

func (cpu *CPU) CPA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.ZFlag((cpu.A - hl) == 0)
	cpu.NFlag(true)
	cpu.HFlag(((cpu.A & 0xF) - (hl & 0xF)) & 0x10 == 0x10)
	cpu.CFlag(hl > cpu.A)
}

func (cpu *CPU) CPA_u8() {
	u8 := cpu.bus.read(cpu.PC)
	cpu.ZFlag((cpu.A - u8) == 0)
	cpu.NFlag(true)
	cpu.HFlag(((cpu.A & 0xF) - (u8 & 0xF)) & 0x10 == 0x10)
	cpu.CFlag(u8 > cpu.A)
	cpu.PC++
}

func (cpu *CPU) DEC_r8(r8 *byte) {
	cpu.HFlag(((*r8 & 0xF) - (1 & 0xF)) & 0x10 == 0x10)
	*r8--
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(true)
}

func (cpu *CPU) DEC_memHL() {
	result := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.HFlag(((result & 0xF) - (1 & 0xF)) & 0x10 == 0x10)
	result--
	cpu.ZFlag(result == 0)
	cpu.NFlag(true)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), result)
}

func (cpu *CPU) INC_r8(r8 *byte) {
	cpu.HFlag(((*r8 & 0xF) + (1 & 0xF)) & 0x10 == 0x10)
	*r8++
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
}

func (cpu *CPU) INC_memHL() {
	result := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.HFlag(((result & 0xF) + (1 & 0xF)) & 0x10 == 0x10)
	result++
	cpu.ZFlag(result == 0)
	cpu.NFlag(true)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), result)
}

func (cpu *CPU) ORA_r8(r8 byte) {
	cpu.A |= r8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
}

func (cpu *CPU) ORA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.A |= hl
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
}

func (cpu *CPU) ORA_u8() {
	u8 := cpu.bus.read(cpu.PC)
	cpu.A |= u8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
	cpu.PC++
}

func (cpu *CPU) SBCA_r8(r8 byte) {
	cpu.HFlag(((cpu.A & 0xF) - (cpu.getCFlag() & 0xF) - (r8 & 0xF)) & 0x10 == 0x10)
	cpu.NFlag(true)
	cpu.CFlag(uint16(cpu.getCFlag()) + uint16(r8) > uint16(cpu.A))
	cpu.A -= (cpu.getCFlag() + r8)
	cpu.ZFlag(cpu.A == 0)
}

func (cpu *CPU) SBCA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.HFlag(((cpu.A & 0xF) - (cpu.getCFlag() & 0xF) - (hl & 0xF)) & 0x10 == 0x10)
	cpu.NFlag(true)
	cpu.CFlag(uint16(cpu.getCFlag()) + uint16(hl) > uint16(cpu.A))
	cpu.A -= (cpu.getCFlag() + hl)
	cpu.ZFlag(cpu.A == 0)
}

func (cpu *CPU) SBCA_u8() {
	u8 := cpu.bus.read(cpu.PC)
	cpu.HFlag(((cpu.A & 0xF) - (cpu.getCFlag() & 0xF) - (u8 & 0xF)) & 0x10 == 0x10)
	cpu.NFlag(true)
	cpu.CFlag(uint16(cpu.getCFlag()) + uint16(u8) > uint16(cpu.A))
	cpu.A -= (cpu.getCFlag() + u8)
	cpu.ZFlag(cpu.A == 0)
	cpu.PC++
}

func (cpu *CPU) SUBA_r8(r8 byte) {
	cpu.HFlag(((cpu.A & 0xF) - (r8 & 0xF)) & 0x10 == 0x10)
	cpu.CFlag(r8 > cpu.A)
	cpu.A -= r8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(true)
}

func (cpu *CPU) SUBA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.HFlag(((cpu.A & 0xF) - (hl & 0xF)) & 0x10 == 0x10)
	cpu.CFlag(hl > cpu.A)
	cpu.A -= hl
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(true)
}

func (cpu *CPU) SUBA_u8() {
	u8 := cpu.bus.read(cpu.PC)
	cpu.HFlag(((cpu.A & 0xF) - (u8 & 0xF)) & 0x10 == 0x10)
	cpu.CFlag(u8 > cpu.A)
	cpu.A -= u8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(true)
	cpu.PC++
}

func (cpu *CPU) XORA_r8(r8 byte) {
	cpu.A ^= r8
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
}

func (cpu *CPU) XORA_memHL() {
	cpu.A ^= cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
}

func (cpu *CPU) XORA_u8() {
	cpu.A ^= cpu.bus.read(cpu.PC)
	cpu.ZFlag(cpu.A == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
	cpu.PC++
}

// 16 bit arithmetic
func (cpu *CPU) ADDHL_r16(r1 byte, r2 byte) {
	hl := uint16(cpu.H) << 8 | uint16(cpu.L)
	r16 := uint16(r1) << 8 | uint16(r2)
	cpu.CFlag(uint32(hl) + uint32(r16) > 0xFFFF)
	cpu.HFlag(((hl & 0x0FFF) + (r16 & 0x0FFF)) & 0x1000 == 0x1000)
	result := hl + r16
	cpu.NFlag(false)
	cpu.H = byte((result >> 8) & 0xFF)
	cpu.L = byte(result)
}

func (cpu *CPU) DEC_r16(r1 *byte, r2 *byte) {
	*r2--
	if *r2 == 0xFF {
		*r1--
	}
	cpu.tick(4)
}

func (cpu *CPU) INC_r16(r1 *byte, r2 *byte) {
	*r2++
	if *r2 == 0 {
		*r1++
	}
	cpu.tick(4)
}

// bit operations
func (cpu *CPU) BITu3_r8(bit int, r8 byte) {
	cpu.ZFlag((r8 & (1 << bit)) >> bit == 0)
	cpu.NFlag(false)
	cpu.HFlag(true)
}

func (cpu *CPU) BITu3_memHL(bit int) {
	hl := ((cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))) & (1 << bit)) >> bit
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(true)
}

func (cpu *CPU) RESu3_r8(bit int, r8 *byte) {
	*r8 &= 0xFF ^ (1 << bit)
}

func (cpu *CPU) RESu3_memHL(bit int) {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	hl &= 0xFF ^ (1 << bit)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

func (cpu *CPU) SWAP_r8(r8 *byte) {
	*r8 = ((*r8 & 0x0F) << 4 | (*r8 >> 4))
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
}

func (cpu *CPU) SWAP_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	hl = ((hl & 0x0F) << 4 | (hl >> 4))
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.CFlag(false)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

// bit shift instructions
func (cpu *CPU) RL_r8(r8 *byte) {
	tempcarry := cpu.getCFlag()
	newcarry := *r8 >> 7
	*r8 = (*r8 << 1 | tempcarry)
	cpu.CFlag(newcarry == 1)
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RL_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	tempcarry := cpu.getCFlag()
	newcarry := hl >> 7
	hl = (hl << 1 | tempcarry)
	cpu.CFlag(newcarry == 1)
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

func (cpu *CPU) RLA() {
	tempcarry := cpu.getCFlag()
	newcarry := cpu.A >> 7
	cpu.A = (cpu.A << 1 | tempcarry)
	cpu.CFlag(newcarry == 1)
	cpu.ZFlag(false)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RLC_r8(r8 *byte) {
	cpu.CFlag((*r8 & 0x80) >> 7 == 1)
	*r8 = (*r8 << 1 | cpu.getCFlag())
	// TODO: check later
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RLC_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.CFlag((hl & 0x80) >> 7 == 1)
	hl = (hl << 1 | cpu.getCFlag())
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

func (cpu *CPU) RLCA() {
	tempcarry := cpu.A >> 7
	cpu.CFlag(tempcarry == 1)
	cpu.A = (cpu.A << 1 | tempcarry)
	cpu.ZFlag(false)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RR_r8(r8 *byte) {
	tempcarry := cpu.getCFlag()
	cpu.CFlag(*r8 & 0x1 == 1)
	*r8 = (tempcarry << 7 | *r8 >> 1)
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RR_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	tempcarry := cpu.getCFlag()
	cpu.CFlag(hl & 0x1 == 1)
	hl = (tempcarry << 7 | hl >> 1)
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

func (cpu *CPU) RRA() {
	tempcarry := cpu.getCFlag()
	newcarry := (cpu.A & 0x01)
	cpu.CFlag(newcarry == 1)
	cpu.A = (tempcarry << 7 | cpu.A >> 1)
	cpu.ZFlag(false)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RRC_r8(r8 *byte) {
	tempcarry := cpu.getCFlag()
	cpu.CFlag(*r8 & 0x1 == 1)
	// TODO: check later
	*r8 = (tempcarry << 7 | *r8 >> 1)
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) RRC_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	tempcarry := cpu.getCFlag()
	cpu.CFlag(hl & 0x1 == 1)
	// TODO: check later
	hl = (tempcarry << 7 | hl >> 1)
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

func (cpu *CPU) RRCA() {
	tempcarry := cpu.getCFlag()
	cpu.CFlag(cpu.A & 0x1 == 1)
	cpu.A = (tempcarry << 7 | cpu.A >> 1)
	cpu.ZFlag(false)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) SLA_r8(r8 *byte) {
	cpu.CFlag((*r8 & (1 << 7) >> 7) == 1)
	*r8 <<= 1
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) SLA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.CFlag((hl & (1 << 7) >> 7) == 1)
	hl <<= 1
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

func (cpu *CPU) SRA_r8(r8 *byte) {
	bit7 := *r8 >> 7
	cpu.CFlag(*r8 & 0x1 == 1)
	*r8 = (bit7 << 7 | *r8 >> 1)
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) SRA_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	bit7 := hl >> 7
	cpu.CFlag(hl & 0x1 == 1)
	hl = (bit7 << 7 | hl >> 1)
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

func (cpu *CPU) SRL_r8(r8 *byte) {
	cpu.CFlag(*r8 & 0x1 == 1)
	*r8 >>= 1
	cpu.ZFlag(*r8 == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) SRL_memHL() {
	hl := cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.CFlag(hl & 0x1 == 1)
	hl >>= 1
	cpu.ZFlag(hl == 0)
	cpu.NFlag(false)
	cpu.HFlag(false)
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), hl)
}

// load instructions
func (cpu *CPU) LDr8_r8(r1 *byte, r2 byte) {
	*r1 = r2
}

func (cpu *CPU) LDr8_u8(r8 *byte) {
	*r8 = cpu.bus.read(cpu.PC)
	cpu.PC++
}

func (cpu *CPU) LDr16_u16(r1 *byte, r2 *byte) {
	*r1 = cpu.bus.read(cpu.PC + 1)
	*r2 = cpu.bus.read(cpu.PC)
	cpu.PC += 2
}

func (cpu *CPU) LDmemHL_r8(r8 byte) {
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), r8)
}

func (cpu *CPU) LDmemHL_u8() {
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), cpu.bus.read(cpu.PC))
	cpu.PC++
}

func (cpu *CPU) LDr8_memHL(r8 *byte) {
	*r8 = cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
}

func (cpu *CPU) LDmemr16_A(r1 byte, r2 byte) {
	cpu.bus.write(uint16(r1) << 8 | uint16(r2), cpu.A)
}

func (cpu *CPU) LDmemu16_A() {
	cpu.bus.write(uint16(cpu.bus.read(cpu.PC + 1)) << 8 | uint16(cpu.PC), cpu.A)
	cpu.PC += 2
}

func (cpu *CPU) LDHmemu16_A() {
	cpu.bus.write((0xFF00 + uint16(cpu.bus.read(cpu.PC))), cpu.A)
	cpu.PC++
}

func (cpu *CPU) LDHmemC_A() {
	cpu.bus.write(uint16(cpu.C) + 0xFF00, cpu.A)
}

func (cpu *CPU) LDA_memr16(r1 byte, r2 byte) {
	cpu.A = cpu.bus.read(uint16(r1) << 8 | uint16(r2))
}

func (cpu *CPU) LDA_memu16() {
	cpu.A = cpu.bus.read(uint16(cpu.PC + 1) << 8 | uint16(cpu.PC))
	cpu.PC += 2
}

func (cpu *CPU) LDHA_memu16() {
	cpu.A = cpu.bus.read(0xFF00 + uint16(cpu.bus.read(cpu.PC)))
	cpu.PC++
}

func (cpu *CPU) LDHA_memC() {
	cpu.A = cpu.bus.read(uint16(cpu.C) + 0xFF00)
}

func (cpu *CPU) LDIHL_A() {
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), cpu.A)
	cpu.L++
	if cpu.L == 0x00 {
		cpu.H++
	}
}

func (cpu *CPU) LDDHL_A() {
	cpu.bus.write(uint16(cpu.H) << 8 | uint16(cpu.L), cpu.A)
	cpu.L--
	if cpu.L == 0xFF {
		cpu.H--
	}
}

func (cpu *CPU) LDA_IHL() {
	cpu.A = cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.L++
	if cpu.L == 0x00 {
		cpu.H++
	}
}

func (cpu *CPU) LDA_DHL() {
	cpu.A = cpu.bus.read(uint16(cpu.H) << 8 | uint16(cpu.L))
	cpu.L--
	if cpu.L == 0xFF {
		cpu.H--
	}
}

// jumps and subroutines
func (cpu *CPU) CALL_u16() {
	cpu.tick(4)
	// TODO: read u16 first
	hi := byte((cpu.PC + 2) >> 8)
	lo := byte(cpu.PC + 2)
	cpu.SP--
	cpu.bus.write(cpu.SP, hi)
	cpu.SP--
	cpu.bus.write(cpu.SP, lo)
	cpu.PC = cpu.bus.read16(cpu.PC)
}

func (cpu *CPU) CALLZ_u16() {
	// TODO: read u16 first
	if cpu.getZFlag() == 1 {
		cpu.tick(4)
		hi := byte((cpu.PC + 2) >> 8)
		lo := byte(((cpu.PC + 2) & 0xFF))
		cpu.SP--
		cpu.bus.write(cpu.SP, hi)
		cpu.SP--
		cpu.bus.write(cpu.SP, lo)
		cpu.PC = cpu.bus.read16(cpu.PC)
	} else {
		cpu.tick(8)
		cpu.PC += 2
	}
}

func (cpu *CPU) CALLNZ_u16() {
	if cpu.getZFlag() == 0 {
		cpu.tick(4)
		hi := byte((cpu.PC + 2) >> 8)
		lo := byte(((cpu.PC + 2) & 0xFF))
		cpu.SP--
		cpu.bus.write(cpu.SP, hi)
		cpu.SP--
		cpu.bus.write(cpu.SP, lo)
		cpu.PC = cpu.bus.read16(cpu.PC)
	} else {
		cpu.tick(8)
		cpu.PC += 2
	}
}

func (cpu *CPU) CALLC_u16() {
	if cpu.getCFlag() == 1 {
		cpu.tick(4)
		hi := byte((cpu.PC + 2) >> 8)
		lo := byte(((cpu.PC + 2) & 0xFF))
		cpu.SP--
		cpu.bus.write(cpu.SP, hi)
		cpu.SP--
		cpu.bus.write(cpu.SP, lo)
		cpu.PC = cpu.bus.read16(cpu.PC)
	} else {
		cpu.tick(8)
		cpu.PC += 2
	}
}

func (cpu *CPU) CALLNC_u16() {
	if cpu.getZFlag() == 0 {
		cpu.tick(4)
		hi := byte((cpu.PC + 2) >> 8)
		lo := byte(((cpu.PC + 2) & 0xFF))
		cpu.SP--
		cpu.bus.write(cpu.SP, hi)
		cpu.SP--
		cpu.bus.write(cpu.SP, lo)
		cpu.PC = cpu.bus.read16(cpu.PC)
	} else {
		cpu.tick(8)
		cpu.PC += 2
	}
}

func (cpu *CPU) JP_HL() {
	cpu.PC = uint16(cpu.H) << 8 | uint16(cpu.L)
}

func (cpu *CPU) JP_u16() {
	cpu.PC = cpu.bus.read16(cpu.PC)
	cpu.tick(4)
}

func (cpu *CPU) JPZ_u16() {
	u16 := cpu.bus.read16(cpu.PC)
	if cpu.getZFlag() == 1 {
		cpu.PC = u16
		cpu.tick(4)
	} else {
		cpu.PC += 2
	}
}

func (cpu *CPU) JPNZ_u16() {
	u16 := cpu.bus.read16(cpu.PC)
	if cpu.getZFlag() == 0 {
		cpu.PC = u16
		cpu.tick(4)
	} else {
		cpu.PC += 2
	}
}

func (cpu *CPU) JPC_u16() {
	u16 := cpu.bus.read16(cpu.PC)
	if cpu.getCFlag() == 1 {
		cpu.PC = u16
		cpu.tick(4)
	} else {
		cpu.PC += 2
	}
}

func (cpu *CPU) JPNC_u16() {
	u16 := cpu.bus.read16(cpu.PC)
	if cpu.getCFlag() == 0 {
		cpu.PC = u16
		cpu.tick(4)
	} else {
		cpu.PC += 2
	}
}

func (cpu *CPU) JR_i8() {
	offset := int8(cpu.bus.read(cpu.PC) + 1)
	cpu.PC += uint16(offset)
	cpu.tick(4)
}

func (cpu *CPU) JRZ_i8() {
	offset := int8(cpu.bus.read(cpu.PC) + 1)
	if cpu.getZFlag() == 1 {
		cpu.PC += uint16(offset)
		cpu.tick(4)
	} else {
		cpu.PC++
	}
}

func (cpu *CPU) JRNZ_i8() {
	offset := int8(cpu.bus.read(cpu.PC) + 1)
	if cpu.getZFlag() == 0 {
		cpu.PC += uint16(offset)
		cpu.tick(4)
	} else {
		cpu.PC++
	}
}

func (cpu *CPU) JRC_i8() {
	offset := int8(cpu.bus.read(cpu.PC) + 1)
	if cpu.getCFlag() == 1 {
		cpu.PC += uint16(offset)
		cpu.tick(4)
	} else {
		cpu.PC++
	}
}

func (cpu *CPU) JRNC_i8() {
	offset := int8(cpu.bus.read(cpu.PC) + 1)
	if cpu.getCFlag() == 0 {
		cpu.PC += uint16(offset)
		cpu.tick(4)
	} else {
		cpu.PC++
	}
}

func (cpu *CPU) RET_Z() {
	cpu.tick(4)
	if cpu.getZFlag() == 1 {
		cpu.PC = cpu.bus.read16(cpu.SP)
		cpu.SP += 2
		cpu.tick(4)
	}
}

func (cpu *CPU) RET_NZ() {
	cpu.tick(4)
	if cpu.getZFlag() == 0 {
		cpu.PC = cpu.bus.read16(cpu.SP)
		cpu.SP += 2
		cpu.tick(4)
	}
}

func (cpu *CPU) RET_C() {
	cpu.tick(4)
	if cpu.getCFlag() == 1 {
		cpu.PC = cpu.bus.read16(cpu.SP)
		cpu.SP += 2
		cpu.tick(4)
	}
}

func (cpu *CPU) RET_NC() {
	cpu.tick(4)
	if cpu.getCFlag() == 0 {
		cpu.PC = cpu.bus.read16(cpu.SP)
		cpu.SP += 2
		cpu.tick(4)
	}
}

func (cpu *CPU) RET() {
	cpu.PC = cpu.bus.read16(cpu.SP)
	cpu.SP += 2
	cpu.tick(4)
}

func (cpu *CPU) RETI() {
	cpu.PC = cpu.bus.read16(cpu.SP)
	cpu.SP += 2
	cpu.tick(2)
	cpu.bus.interrupt.IME = 1
}

func (cpu *CPU) RST_vec(vec uint16) {
	hi := byte((cpu.PC) >> 8)
	lo := byte(((cpu.PC) & 0xFF)) // inefficient change later
	cpu.tick(4)
	cpu.SP--
	cpu.bus.write(cpu.SP, hi)
	cpu.SP--
	cpu.bus.write(cpu.SP, lo)
	cpu.PC = vec
}

// stack operations instructions
func (cpu *CPU) ADDHL_SP() {
	hl := uint16(cpu.H) << 8 | uint16(cpu.L)
	cpu.CFlag(uint32(hl) + uint32(cpu.SP) > 0xFFFF)
	cpu.HFlag(((hl & 0x0FFF) + (cpu.SP & 0x0FFF)) & 0x1000 == 0x1000)
	result := hl + cpu.SP
	cpu.NFlag(false)
	cpu.H = byte((result >> 8) & 0xFF)
	cpu.L = byte(result)
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

func (cpu *CPU) DEC_SP() {
	cpu.SP--
	cpu.tick(4)
}

func (cpu *CPU) INC_SP() {
	cpu.SP++
	cpu.tick(4)
}

func (cpu *CPU) LDSP_u16() {
	cpu.SP = cpu.bus.read16(cpu.PC)
	cpu.PC += 2
}

func (cpu *CPU) LDmemu16_SP() {
	u16 := cpu.bus.read16(cpu.PC)
	cpu.bus.write(u16, byte(cpu.SP))
	cpu.bus.write(u16 + 1, byte(cpu.SP >> 8))
	cpu.PC += 2
}

func (cpu *CPU) LDHL_SPi8() {
	i8 := int8(cpu.bus.read(cpu.PC))
	result := uint16(cpu.SP) + uint16(i8)
	cpu.HFlag(((cpu.SP & 0xF) + ((uint16(i8)) & 0xF)) >= 0x10)
	cpu.CFlag((uint16(i8) & 0xFF) + (cpu.SP & 0xFF) > 0xFF)
	cpu.H = byte(result >> 8)
	cpu.L = byte(result & 0xFF)
	cpu.ZFlag(false)
	cpu.NFlag(false)
	cpu.PC++
}

func (cpu *CPU) LDSP_HL() {
	cpu.SP = uint16(cpu.H) << 8 | uint16(cpu.L)
}

func (cpu *CPU) POP_AF() {
	cpu.A = cpu.bus.read(cpu.SP) & 0xF0
	cpu.SP++
	cpu.F = cpu.bus.read(cpu.SP)
	cpu.SP++
}

func (cpu *CPU) POP_r16(r1 *byte, r2 *byte) {
	*r2 = cpu.bus.read(cpu.SP)
	cpu.SP++
	*r1 = cpu.bus.read(cpu.SP)
	cpu.SP++
}

func (cpu *CPU) PUSH_AF() {
	cpu.SP--
	cpu.bus.write(cpu.SP, cpu.F)
	cpu.SP--
	cpu.bus.write(cpu.SP, cpu.A)
}

func (cpu *CPU) PUSH_r16(r1 byte, r2 byte) {
	cpu.tick(4)
	cpu.SP--
	cpu.bus.write(cpu.SP, r1)
	cpu.SP--
	cpu.bus.write(cpu.SP, r2)
}

// miscellaneous instructions
func (cpu *CPU) CCF() {
	carry := cpu.getCFlag()
	cpu.CFlag(carry == 1)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) CPL() {
	cpu.A = ^cpu.A
	cpu.NFlag(true)
	cpu.HFlag(true)
}

func (cpu *CPU) DAA() {
	// addition
	if cpu.getNFlag() == 0 {
		if cpu.A > 0x99 || cpu.getCFlag() == 1 {
			cpu.A += 0x60
			cpu.CFlag(true)
		}
		if cpu.getHFlag() == 1 || (cpu.A & 0xF) > 0x09 {
			cpu.A += 0x06
		}
	} else {
		// subtraction
		if cpu.getCFlag() == 1 {
			cpu.A -= 0x60
		}
		if cpu.getHFlag() == 1 {
			cpu.A -= 0x06
		}
	}
	cpu.ZFlag(cpu.A == 0)
	cpu.HFlag(false)
}

func (cpu *CPU) DI() {
	// todo interrupt delay or whatever
	cpu.bus.interrupt.IME = 0
}

func (cpu *CPU) EI() {
	cpu.bus.interrupt.IMEDelay = true
}

func (cpu *CPU) HALT() {
	cpu.halt = true
}

func (cpu *CPU) NOP() {

}

func (cpu *CPU) SCF() {
	cpu.CFlag(true)
	cpu.NFlag(false)
	cpu.HFlag(false)
}

func (cpu *CPU) STOP() {
	if cpu.bus.read(cpu.PC) == 0x00 {
		cpu.PC++
	}
	// needs work
}

func (cpu *CPU) CB() {
	// get the cb opcode
	cpu.opcode = cpu.bus.read(cpu.PC)
	cbopcodes[cpu.opcode].exec(cpu)
	cpu.PC++ // since none of the cb opcodes change PC we increment after the instruction is done
}

// undefined opcodes
func (cpu *CPU) UND_OP() {
	fmt.Println("undefined opcode! 0x" + strconv.FormatUint(uint64(cpu.opcode), 16))
	os.Exit(0)
}

func (cpu *CPU) UND_CBOP() {
	fmt.Println("undefined cb opcode! 0x" + strconv.FormatUint(uint64(cpu.opcode), 16))
	os.Exit(0)
}

var opcodes = [256]Opcode {
	Opcode{"NOP", func(cpu *CPU) { cpu.NOP() }}, Opcode{"LD BC, u16", func(cpu *CPU) { cpu.LDr16_u16(&cpu.B, &cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC BC", func(cpu *CPU) { cpu.INC_r16(&cpu.B, &cpu.C) }}, Opcode{"INC B", func(cpu *CPU) { cpu.INC_r8(&cpu.B) }}, Opcode{"DEC B", func(cpu *CPU) { cpu.DEC_r8(&cpu.B) }}, Opcode{"LD B, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.B) }}, Opcode{"RLCA", func(cpu *CPU) { cpu.RLCA() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"ADD HL, BC", func(cpu *CPU) { cpu.ADDHL_r16(cpu.B, cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC C", func(cpu *CPU) { cpu.INC_r8(&cpu.C) }}, Opcode{"DEC C", func(cpu *CPU) { cpu.DEC_r8(&cpu.C) }}, Opcode{"LD C, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD DE, u16", func(cpu *CPU) { cpu.LDr16_u16(&cpu.D, &cpu.E) }}, Opcode{"LD (DE), A", func(cpu *CPU) { cpu.LDmemr16_A(cpu.D, cpu.E) }}, Opcode{"INC DE", func(cpu *CPU) { cpu.INC_r16(&cpu.D, &cpu.E) }}, Opcode{"INC D", func(cpu *CPU) { cpu.INC_r8(&cpu.D) }}, Opcode{"DEC D", func(cpu *CPU) { cpu.DEC_r8(&cpu.D) }}, Opcode{"LD D, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.D) }}, Opcode{"RLA", func(cpu *CPU) { cpu.RLA() }}, Opcode{"JR i8", func(cpu *CPU) { cpu.JR_i8() }}, Opcode{"ADD HL, DE", func(cpu *CPU) { cpu.ADDHL_r16(cpu.D, cpu.E) }}, Opcode{"LD A, (DE)", func(cpu *CPU) { cpu.LDA_memr16(cpu.D, cpu.E) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC E", func(cpu *CPU) { cpu.INC_r8(&cpu.E) }}, Opcode{"DEC E", func(cpu *CPU) { cpu.DEC_r8(&cpu.E) }}, Opcode{"LD E, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.E) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"JR NZ, i8", func(cpu *CPU) { cpu.JRNZ_i8() }}, Opcode{"LD HL, u16", func(cpu *CPU) { cpu.LDr16_u16(&cpu.H, &cpu.L) }}, Opcode{"LD (HL+), A", func(cpu *CPU) { cpu.LDIHL_A() }}, Opcode{"INC HL", func(cpu *CPU) { cpu.INC_r16(&cpu.H, &cpu.L) }}, Opcode{"INC H", func(cpu *CPU) { cpu.INC_r8(&cpu.H) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD H, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.H) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"JR Z, i8", func(cpu *CPU) { cpu.JRZ_i8() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD A, (HL+)", func(cpu *CPU) { cpu.LDA_IHL() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD L, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.L) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD SP, u16", func(cpu *CPU) { cpu.LDSP_u16() }}, Opcode{"LD (HL-), A", func(cpu *CPU) { cpu.LDDHL_A() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"DEC (HL)", func(cpu *CPU) { cpu.DEC_memHL() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"INC A", func(cpu *CPU) { cpu.INC_r8(&cpu.A) }}, Opcode{"DEC A", func(cpu *CPU) { cpu.DEC_r8(&cpu.A) }}, Opcode{"LD A, u8", func(cpu *CPU) { cpu.LDr8_u8(&cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"LD B, B", func(cpu *CPU) { cpu.LDr8_r8(&cpu.B, cpu.B) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD B, A", func(cpu *CPU) { cpu.LDr8_r8(&cpu.B, cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD C, A", func(cpu *CPU) { cpu.LDr8_r8(&cpu.C, cpu.A) }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD D, A", func(cpu *CPU) { cpu.LDr8_r8(&cpu.D, cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD H, A", func(cpu *CPU) { cpu.LDr8_r8(&cpu.H, cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD (HL), A", func(cpu *CPU) { cpu.LDmemr16_A(cpu.H, cpu.L) }}, Opcode{"LD A, B", func(cpu *CPU) { cpu.LDr8_r8(&cpu.A, cpu.B) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD A, E", func(cpu *CPU) { cpu.LDr8_r8(&cpu.A, cpu.E) }}, Opcode{"LD A, H", func(cpu *CPU) { cpu.LDr8_r8(&cpu.A, cpu.H) }}, Opcode{"LD A, L", func(cpu *CPU) { cpu.LDr8_r8(&cpu.A, cpu.L) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"ADD A, B", func(cpu *CPU) { cpu.ADDA_r8(cpu.B) }}, Opcode{"ADD A, C", func(cpu *CPU) { cpu.ADDA_r8(cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"ADD A, (HL)", func(cpu *CPU) { cpu.ADDA_memHL() }}, Opcode{"ADD A, A", func(cpu *CPU) { cpu.ADDA_r8(cpu.A) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"SUB A, B", func(cpu *CPU) { cpu.SUBA_r8(cpu.B) }}, Opcode{"SUB A, C", func(cpu *CPU) { cpu.SUBA_r8(cpu.C) }}, Opcode{"SUB A, D", func(cpu *CPU) { cpu.SUBA_r8(cpu.D) }}, Opcode{"SUB A, E", func(cpu *CPU) { cpu.SUBA_r8(cpu.E) }}, Opcode{"SUB A, H", func(cpu *CPU) { cpu.SUBA_r8(cpu.H) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"AND A, A", func(cpu *CPU) { cpu.ANDA_r8(cpu.A) }}, Opcode{"XOR A, B", func(cpu *CPU) { cpu.XORA_r8(cpu.B) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"XOR A, A", func(cpu *CPU) { cpu.XORA_r8(cpu.A) }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CP A, H", func(cpu *CPU) { cpu.CPA_r8(cpu.H) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CP A, (HL)", func(cpu *CPU) { cpu.CPA_memHL() }}, Opcode{"CP A, A", func(cpu *CPU) { cpu.CPA_r8(cpu.A) }}, 
	Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"POP BC", func(cpu *CPU) { cpu.POP_r16(&cpu.B, &cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"JP u16", func(cpu *CPU) { cpu.JP_u16() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"PUSH BC", func(cpu *CPU) { cpu.PUSH_r16(cpu.B, cpu.C) }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"RET", func(cpu *CPU) { cpu.RET() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CB", func(cpu *CPU) { cpu.CB() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CALL u16", func(cpu *CPU) { cpu.CALL_u16() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"RET NC", func(cpu *CPU) { cpu.RET_NC() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"LD (FF00+u8), A", func(cpu *CPU) { cpu.LDHmemu16_A() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD (FF00 + C), A", func(cpu *CPU) { cpu.LDHmemC_A() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"ADD SP, i8", func(cpu *CPU) { cpu.ADDSP_i8() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"LD (u16), A", func(cpu *CPU) { cpu.LDmemu16_A() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, 
	Opcode{"LD A, (FF00+u8)", func(cpu *CPU) { cpu.LDHA_memu16() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"DI", func(cpu *CPU) { cpu.DI() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"NOP", func(cpu *CPU) { cpu.UND_OP() }}, Opcode{"CP A, u8", func(cpu *CPU) { cpu.CPA_u8() }}, Opcode{"RST 0x38", func(cpu *CPU) { cpu.RST_vec(0x0038) }}, 
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