/*
 * Copyright 2025 Nathanne Isip
 * This file is part of Risbee (https://github.com/nthnn/risbee)
 * This code is licensed under MIT license (see LICENSE for details)
 */

package risbee

// Instruction decoding constants:
// Primary opcodes defining the major RISC-V instruction formats.
const (
	// RISBEE_OPINST_LOAD is the opcode for load instructions (e.g., LB, LW).
	RISBEE_OPINST_LOAD = 3
	// RISBEE_OPINST_STORE is the opcode for store instructions (e.g., SB, SW).
	RISBEE_OPINST_STORE = 35
	// RISBEE_OPINST_IMM is the opcode for immediate ALU operations (I-type).
	RISBEE_OPINST_IMM = 19
	// RISBEE_OPINST_IALU is the opcode for 64-bit immediate ALU operations.
	RISBEE_OPINST_IALU = 27
	// RISBEE_OPINST_RT64 is the opcode for 64-bit register–register operations.
	RISBEE_OPINST_RT64 = 51
	// RISBEE_OPINST_RT32 is the opcode for 32-bit register–register word ops.
	RISBEE_OPINST_RT32 = 59
	// RISBEE_OPINST_LUI is the opcode for Load Upper Immediate (U-type).
	RISBEE_OPINST_LUI = 55
	// RISBEE_OPINST_AUIPC is the opcode for Add Upper Immediate to PC (U-type).
	RISBEE_OPINST_AUIPC = 23
	// RISBEE_OPINST_JAL is the opcode for Jump and Link (J-type).
	RISBEE_OPINST_JAL = 111
	// RISBEE_OPINST_JALR is the opcode for Jump and Link Register (I-type).
	RISBEE_OPINST_JALR = 103
	// RISBEE_OPINST_BRANCH is the opcode for conditional branches (B-type).
	RISBEE_OPINST_BRANCH = 99
	// RISBEE_OPINST_FENCE is the opcode for memory ordering fences.
	RISBEE_OPINST_FENCE = 15
	// RISBEE_OPINST_CALL is the opcode for environment calls / syscalls.
	RISBEE_OPINST_CALL = 115
)

// Function3 codes for load instruction variants (determines width and sign).
const (
	RISBEE_FC3_LB   = 0 // Load Byte (signed)
	RISBEE_FC3_LHW  = 1 // Load Halfword (signed 16-bit)
	RISBEE_FC3_LW   = 2 // Load Word (signed 32-bit)
	RISBEE_FC3_LDW  = 3 // Load Doubleword (64-bit)
	RISBEE_FC3_LBU  = 4 // Load Byte Unsigned
	RISBEE_FC3_LHU  = 5 // Load Halfword Unsigned
	RISBEE_FC3_LRES = 6 // Reserved or custom load
)

// Function3 codes for store instruction variants (determines width).
const (
	RISBEE_FC3_SB  = 0 // Store Byte
	RISBEE_FC3_SHW = 1 // Store Halfword
	RISBEE_FC3_SW  = 2 // Store Word
	RISBEE_FC3_SDW = 3 // Store Doubleword
)

// Function3 codes for I-type arithmetic instructions.
const (
	RISBEE_FC3_ADDI  = 0 // Add Immediate
	RISBEE_FC3_SLLI  = 1 // Shift Left Logical Immediate
	RISBEE_FC3_SLTI  = 2 // Set Less Than Immediate (signed)
	RISBEE_FC3_SLTIU = 3 // Set Less Than Immediate (unsigned)
	RISBEE_FC3_XORI  = 4 // XOR Immediate
	RISBEE_FC3_SRLI  = 5 // Shift Right Logical Immediate or Arithmetic (distinguished by funct7)
	RISBEE_FC3_ORI   = 6 // OR Immediate
	RISBEE_FC3_ANDI  = 7 // AND Immediate
)

// Function3 codes for 64-bit and word-width IALU operations.
const (
	RISBEE_FC3_SLLIW  = 0 // Shift Left Logical Immediate for 32-bit
	RISBEE_FC3_SRLIW  = 1 // Shift Right Logical Immediate for 32-bit
	RISBEE_FC3_SRAIW  = 5 // Shift Right Arithmetic Immediate for 32-bit
	RISBEE_FC3_SLLI64 = 6 // Shift Left Logical Immediate for 64-bit
	RISBEE_FC3_SRLI64 = 7 // Shift Right Logical Immediate for 64-bit
	RISBEE_FC3_SRAI64 = 8 // Shift Right Arithmetic Immediate for 64-bit
)

// Function3 codes for conditional branch types.
const (
	RISBEE_FC3_BEQ  = 0 // Branch if Equal
	RISBEE_FC3_BNE  = 1 // Branch if Not Equal
	RISBEE_FC3_BLT  = 4 // Branch if Less Than (signed)
	RISBEE_FC3_BGE  = 5 // Branch if Greater or Equal (signed)
	RISBEE_FC3_BLTU = 6 // Branch if Less Than (unsigned)
	RISBEE_FC3_BGEU = 7 // Branch if Greater or Equal (unsigned)
)

// Combined funct7 and funct3 codes for RT32 (32-bit register) operations.
const (
	RISBEE_OPINST_RT32_ADDW  = 0x0   // ADDW: add word
	RISBEE_OPINST_RT32_SUBW  = 0x100 // SUBW: subtract word
	RISBEE_OPINST_RT32_SLLW  = 0x1   // SLLW: shift left logical word
	RISBEE_OPINST_RT32_SRLW  = 0x5   // SRLW: shift right logical word
	RISBEE_OPINST_RT32_SRAW  = 0x105 // SRAW: shift right arithmetic word
	RISBEE_OPINST_RT32_MULW  = 0x8   // MULW: multiply word
	RISBEE_OPINST_RT32_DIVW  = 0xC   // DIVW: divide word (signed)
	RISBEE_OPINST_RT32_DIVUW = 0xD   // DIVUW: divide word (unsigned)
	RISBEE_OPINST_RT32_REMW  = 0xE   // REMW: remainder word (signed)
	RISBEE_OPINST_RT32_REMUW = 0xF   // REMUW: remainder word (unsigned)
)

// Combined funct7 and funct3 codes for RT64 (64-bit register) operations.
const (
	RISBEE_OPINST_RT64_ADD    = 0x0   // ADD: add
	RISBEE_OPINST_RT64_SUB    = 0x100 // SUB: subtract
	RISBEE_OPINST_RT64_SLL    = 0x1   // SLL: shift left logical
	RISBEE_OPINST_RT64_SLT    = 0x2   // SLT: set less than (signed)
	RISBEE_OPINST_RT64_SLTU   = 0x3   // SLTU: set less than (unsigned)
	RISBEE_OPINST_RT64_XOR    = 0x4   // XOR: exclusive OR
	RISBEE_OPINST_RT64_SRL    = 0x5   // SRL: shift right logical
	RISBEE_OPINST_RT64_SRA    = 0x105 // SRA: shift right arithmetic
	RISBEE_OPINST_RT64_OR     = 0x6   // OR: bitwise OR
	RISBEE_OPINST_RT64_AND    = 0x7   // AND: bitwise AND
	RISBEE_OPINST_RT64_MUL    = 0x8   // MUL: multiply
	RISBEE_OPINST_RT64_MULH   = 0x9   // MULH: multiply high signed
	RISBEE_OPINST_RT64_MULHSU = 0xA   // MULHSU: multiply high signed×unsigned
	RISBEE_OPINST_RT64_MULHU  = 0xB   // MULHU: multiply high unsigned
	RISBEE_OPINST_RT64_DIV    = 0xC   // DIV: divide (signed)
	RISBEE_OPINST_RT64_DIVU   = 0xD   // DIVU: divide (unsigned)
	RISBEE_OPINST_RT64_REM    = 0xE   // REM: remainder (signed)
	RISBEE_OPINST_RT64_REMU   = 0xF   // REMU: remainder (unsigned)
)
