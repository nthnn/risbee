/*
 * Copyright 2025 Nathanne Isip
 * This file is part of Risbee (https://github.com/nthnn/risbee)
 * This code is licensed under MIT license (see LICENSE for details)
 */

package risbee

const (
	// RISBEE_STACK_SIZE defines the total size of the VM memory (64 KiB).
	RISBEE_STACK_SIZE = 65536
)

// RisbeeVmSyscallFn represents the signature of a syscall handler function.
type RisbeeVmSyscallFn func(vm *RisbeeVm) uint64

// RisbeeVm encapsulates the state of the RISC-V inspired virtual machine.
// It includes memory, registers, program counter, exit code, running status,
// and a map of registered syscall handlers.
type RisbeeVm struct {
	Memory        [RISBEE_STACK_SIZE]byte      // Byte-addressable VM memory
	Registers     [32]uint64                   // General-purpose registers R0–R31
	Pc            uint64                       // Program counter
	ExitCode      int                          // Exit code of the VM
	Running       bool                         // VM running status
	SysCalls      map[uint64]RisbeeVmSyscallFn // Registered syscalls
	ExitCallback  func(uint64)                 // Exit system call callback function
	PanicCallback func(string)                 // Panic callback function
}

// This function initializes the Risbee virtual machine
// instance. It sets up the initial state of the virtual machine.
//
// Parameters:
//   - exitCallback Callback triggered when system
//     call for exit is invoked
//   - panicCallback Callback for encountered panic errors
func (vm *RisbeeVm) Initialize(
	exitCallback func(uint64),
	panicCallback func(string),
) {
	vm.Pc = 4096
	vm.ExitCode = 0
	vm.Running = false
	vm.SysCalls = map[uint64]RisbeeVmSyscallFn{}

	vm.ExitCallback = exitCallback
	vm.PanicCallback = panicCallback
}

// Stops the execution of the virtual machine.
// This method halts the execution of the virtual machine.
func (vm *RisbeeVm) Stop() {
	vm.Running = false
}

// Gets the exit code of the virtual machine.
//
// Returns the exit code.
func (vm *RisbeeVm) GetExitCode() int {
	return vm.ExitCode
}

// It copies the contents of `data` into the VM’s internal
// memory starting at the fixed load offset (4096 bytes),
// mirroring the behavior of LoadFile but without disk I/O.
func (vm *RisbeeVm) LoadFromBytes(Data []byte) bool {
	const loadOffset = 4096
	size := len(Data)

	if size == 0 ||
		uint64(size) > RISBEE_STACK_SIZE-loadOffset {
		return false
	}

	copy(vm.Memory[loadOffset:loadOffset+size], Data)
	vm.Registers[2] = RISBEE_STACK_SIZE

	return true
}

// This function starts the execution of the Risbee virtual machine
// instance and executes the loaded program, if any, and handles any
// system calls or instructions encountered during execution until
// the program exits or an error occurs.
func (vm *RisbeeVm) Run() {
	vm.Running = true
	for vm.Running {
		inst := vm.fetch()
		vm.execute(inst)
	}
}

// This method returns a boolean value indicating whether the
// virtual machine is currently running or not.
func (vm *RisbeeVm) IsRunning() bool {
	return vm.Running
}

// GetPointerParam retrieves the syscall
// parameter from register a0+aIndex.
func (vm *RisbeeVm) GetPointerParam(Index uint64) uint64 {
	return vm.Registers[10+Index]
}

// GetStringPointer reads a null-terminated string
// from VM memory at the given pointer.
//
// Returns "(null)" if the pointer is zero, or the
// extracted string otherwise.
func (vm *RisbeeVm) GetStringPointer(Pointer uint64) string {
	var str string
	if Pointer == 0 {
		str = "(null)"
	} else {
		mem := vm.Memory
		start := Pointer
		end := Pointer

		for end < uint64(len(mem)) && mem[end] != 0 {
			end++
		}

		str = string(mem[start:end])
	}

	return str
}

// SetSystemCall registers a syscall handler
// function at the given address (code).
func (vm *RisbeeVm) SetSystemCall(
	Address uint64,
	Callback RisbeeVmSyscallFn,
) {
	vm.SysCalls[Address] = Callback
}

// GetSystemCall retrieves a registered syscall
// handler by address (code).
//
// Returns the handler and a boolean indicating existence.
func (vm *RisbeeVm) GetSystemCall(
	Address uint64,
) (RisbeeVmSyscallFn, bool) {
	callback, ok := vm.SysCalls[Address]
	return callback, ok
}

// Sets the exit code of the virtual machine.
//
// Parameters:
// - exitCode The exit code.
func (vm *RisbeeVm) setExitCode(exitCode int) {
	vm.ExitCode = exitCode
}

// Handle a panic situation in the virtual machine.
//
// This function is called to handle a panic situation
// in the Risbee virtual machine. It prints the panic
// message and performs any necessary cleanup before
// terminating the program.
func (vm *RisbeeVm) panic(message string) {
	if vm.PanicCallback != nil {
		vm.PanicCallback(message)
	}

	vm.Stop()
	vm.setExitCode(-1)
}

// Fetches the next instruction to be executed in a virtual machine.
//
// This function fetches the next instruction from the program counter of
// the specified Risbee virtual machine instance vm. It returns the fetched
// instruction for execution by the virtual machine.
//
// Returns the next instruction to be executed.
func (vm *RisbeeVm) fetch() uint32 {
	return uint32LittleEndian(
		vm.Memory[vm.Pc : vm.Pc+4],
	)
}

// Handles a system call in a Risbee virtual machine instance.
//
// This function handles a system call specified by code within the
// specified Risbee virtual machine instance vm. It executes the
// corresponding system call routine and returns the result.
//
// Parameters:
// - code The system call code to be handled.
//
// Returns the result of the system call execution.
func (vm *RisbeeVm) handleSyscall(code uint64) uint64 {
	if code == 0 {
		exitCode := int(vm.GetPointerParam(0))
		vm.setExitCode(exitCode)

		if vm.ExitCallback != nil {
			vm.ExitCallback(uint64(exitCode))
			vm.Stop()
		}

		return uint64(exitCode)
	} else if fn, ok := vm.SysCalls[code]; ok {
		return fn(vm)
	} else {
		vm.panic("Invalid system call.")
	}

	return 0
}

// Executes the given instruction.
//
// Parameters:
// - inst The instruction to execute.
func (vm *RisbeeVm) execute(inst uint32) {
	opcode := inst & 0x7F

	rd := (inst >> 7) & 0x1F
	rs1 := (inst >> 15) & 0x1F
	rs2 := (inst >> 20) & 0x1F

	switch opcode {
	case RISBEE_OPINST_LOAD:
		functionCode3 := (inst >> 12) & 0x7
		immediate := int64(int32(inst&0xFFF00000) >> 20)
		addr := vm.Registers[rs1] + uint64(immediate)

		var val int64
		switch functionCode3 {
		case RISBEE_FC3_LB:
			val = int64(int8(vm.Memory[addr]))

		case RISBEE_FC3_LHW:
			val = int64(int16(uint16LittleEndian(
				vm.Memory[addr:],
			)))

		case RISBEE_FC3_LW:
			val = int64(int32(uint32LittleEndian(
				vm.Memory[addr:],
			)))

		case RISBEE_FC3_LDW:
			val = int64(uint64LittleEndian(
				vm.Memory[addr:],
			))

		case RISBEE_FC3_LBU:
			val = int64(vm.Memory[addr])

		case RISBEE_FC3_LHU:
			val = int64(uint16LittleEndian(
				vm.Memory[addr:],
			))

		case RISBEE_FC3_LRES:
			val = int64(uint32LittleEndian(
				vm.Memory[addr:],
			))

		default:
			vm.panic("Invalid load instruction.")
		}

		if rd != 0 {
			vm.Registers[rd] = uint64(val)
		}

	case RISBEE_OPINST_STORE:
		functionCode3 := (inst >> 12) & 0x7

		imm11_5 := (inst >> 20) & 0xFE0
		imm4_0 := (inst >> 7) & 0x1F
		immediate := int64(
			int32((imm11_5|imm4_0)<<20) >> 20,
		)

		addr := vm.Registers[rs1] + uint64(immediate)
		val := vm.Registers[rs2]

		switch functionCode3 {
		case RISBEE_FC3_SB:
			vm.Memory[addr] = byte(val)

		case RISBEE_FC3_SHW:
			putUint16(
				vm.Memory[addr:],
				uint16(val),
			)

		case RISBEE_FC3_SW:
			putUint32(
				vm.Memory[addr:],
				uint32(val),
			)

		case RISBEE_FC3_SDW:
			putUint64(
				vm.Memory[addr:],
				val,
			)

		default:
			vm.panic("Invalid store instruction.")
		}

	case RISBEE_OPINST_IMM:
		functionCode3 := (inst >> 12) & 0x7
		immediate := int64(int32(inst&0xFFF00000) >> 20)
		shiftAmount := (inst >> 20) & 0x3F

		val := int64(vm.Registers[rs1])
		switch functionCode3 {
		case RISBEE_FC3_ADDI:
			val = val + immediate

		case RISBEE_FC3_SLLI:
			val = shiftLeftInt64(val, int64(shiftAmount))

		case RISBEE_FC3_SLTI:
			if val < immediate {
				val = 1
			} else {
				val = 0
			}

		case RISBEE_FC3_SLTIU:
			if uint64(val) < uint64(immediate) {
				val = 1
			} else {
				val = 0
			}

		case RISBEE_FC3_XORI:
			val = val ^ immediate

		case RISBEE_FC3_SRLI:
			functionCode6 := (inst >> 26) & 0x3F

			switch functionCode6 >> 4 {
			case 0x0:
				val = shiftRightInt64(
					val,
					int64(shiftAmount),
				)

			case 0x1:
				val = arithShiftRightInt64(
					val,
					int64(shiftAmount),
				)

			default:
				vm.panic("Invalid immediate shift instruction.")
			}

		case RISBEE_FC3_ORI:
			val = val | immediate

		case RISBEE_FC3_ANDI:
			val = val & immediate

		default:
			vm.panic("Invalid immediate instruction.")
		}

		if rd != 0 {
			vm.Registers[rd] = uint64(val)
		}

	case RISBEE_OPINST_IALU:
		functionCode3 := (inst >> 12) & 0x7
		immediate := int64(int32(inst&0xFFF00000) >> 20)

		val := int64(vm.Registers[rs1])
		switch functionCode3 {
		case RISBEE_FC3_SLLIW:
			val = int64(val + immediate)

		case RISBEE_FC3_SRLIW:
			val = shiftLeftInt64(val, immediate)

		case RISBEE_FC3_SRAIW:
			shiftAmount := rs2
			functionCode7 := (inst >> 25) & 0x7F

			switch functionCode7 >> 5 {
			case 0x0:
				val = shiftRightInt64(
					val,
					int64(shiftAmount),
				)

			case 0x1:
				val = arithShiftRightInt64(
					val,
					int64(shiftAmount),
				)

			default:
				vm.panic("Invalid immediate shift instruction.")
			}

		case RISBEE_FC3_SLLI64:
			val = shiftLeftInt64(
				val,
				immediate&0x3F,
			)

		case RISBEE_FC3_SRLI64:
			val = shiftRightInt64(
				val,
				immediate&0x3F,
			)

		case RISBEE_FC3_SRAI64:
			val = arithShiftRightInt64(
				val,
				immediate&0x3F,
			)

		default:
			vm.panic("Invalid immediate instruction.")
		}

		if rd != 0 {
			vm.Registers[rd] = uint64(val)
		}

	case RISBEE_OPINST_RT64:
		functionCode3 := (inst >> 12) & 0x7
		functionCode7 := (inst >> 25) & 0x7F

		val1 := int64(vm.Registers[rs1])
		val2 := int64(vm.Registers[rs2])

		var val int64
		switch (functionCode7 << 3) | functionCode3 {
		case RISBEE_OPINST_RT64_ADD:
			val = val1 + val2

		case RISBEE_OPINST_RT64_SUB:
			val = val1 - val2

		case RISBEE_OPINST_RT64_SLL:
			val = shiftLeftInt64(val1, val2&0x1F)

		case RISBEE_OPINST_RT64_SLT:
			if val1 < val2 {
				val = 1
			} else {
				val = 0
			}

		case RISBEE_OPINST_RT64_SLTU:
			if uint64(val1) < uint64(val2) {
				val = 1
			} else {
				val = 0
			}

		case RISBEE_OPINST_RT64_XOR:
			val = val1 ^ val2

		case RISBEE_OPINST_RT64_SRL:
			val = shiftRightInt64(val1, val2&0x1F)

		case RISBEE_OPINST_RT64_SRA:
			val = arithShiftRightInt64(
				val1,
				val2&0x1F,
			)

		case RISBEE_OPINST_RT64_OR:
			val = val1 | val2

		case RISBEE_OPINST_RT64_AND:
			val = val1 & val2

		case RISBEE_OPINST_RT64_MUL:
			val = val1 * val2

		case RISBEE_OPINST_RT64_MULH:
			val = shiftRightInt128(val1*val2, 64)

		case RISBEE_OPINST_RT64_MULHSU:
			val = int64(uint64(shiftRightInt128(
				val1*int64(uint64(val2)),
				64,
			)))

		case RISBEE_OPINST_RT64_MULHU:
			val = int64(uint64(shiftRightInt128(
				int64(uint64(val1))*int64(uint64(val2)),
				64,
			)))

		case RISBEE_OPINST_RT64_DIV:
			dividend, divisor := val1, val2
			if dividend == (-9223372036854775807-1) &&
				divisor == -1 {
				val = -9223372036854775807 - 1
			} else if divisor == 0 {
				val = -1
			} else {
				val = dividend / divisor
			}

		case RISBEE_OPINST_RT64_DIVU:
			dividend, divisor := uint64(val1), uint64(val2)
			if divisor == 0 {
				val = -1
			} else {
				val = int64(dividend / divisor)
			}

		case RISBEE_OPINST_RT64_REM:
			dividend, divisor := val1, val2
			if dividend == (-9223372036854775807-1) &&
				divisor == -1 {
				val = 0
			} else if divisor == 0 {
				val = dividend
			} else {
				val = dividend % divisor
			}

		case RISBEE_OPINST_RT64_REMU:
			dividend := uint64(val1)
			divisor := uint64(val2)

			if divisor == 0 {
				val = int64(dividend)
			} else {
				val = int64(dividend % divisor)
			}

		default:
			vm.panic("Invalid arith instruction.")
		}

		if rd != 0 {
			vm.Registers[rd] = uint64(val)
		}

	case RISBEE_OPINST_RT32:
		functionCode3 := (inst >> 12) & 0x7
		functionCode7 := (inst >> 25) & 0x7F

		val1 := int64(vm.Registers[rs1])
		val2 := int64(vm.Registers[rs2])

		var val int64
		switch (functionCode7 << 3) | functionCode3 {
		case RISBEE_OPINST_RT32_ADDW:
			val = int64(val1 + val2)

		case RISBEE_OPINST_RT32_SUBW:
			val = int64(val1 - val2)

		case RISBEE_OPINST_RT32_SLLW:
			val = shiftLeftInt64(val1, val2&0x1F)

		case RISBEE_OPINST_RT32_SRLW:
			val = shiftRightInt64(val1, val2&0x1F)

		case RISBEE_OPINST_RT32_SRAW:
			val = int64(arithShiftRightInt64(
				val1,
				val2&0x1F,
			))

		case RISBEE_OPINST_RT32_MULW:
			val = int64(int32(val1) * int32(val2))

		case RISBEE_OPINST_RT32_DIVW:
			dividend := int32(val1)
			divisor := int32(val2)

			if dividend == (-2147483647-1) &&
				divisor == -1 {
				val = -2147483648
			} else if divisor == 0 {
				val = -1
			} else {
				val = int64(dividend / divisor)
			}

		case RISBEE_OPINST_RT32_DIVUW:
			dividend := uint32(val1)
			divisor := uint32(val2)

			if divisor == 0 {
				val = -1
			} else {
				val = int64(dividend / divisor)
			}

		case RISBEE_OPINST_RT32_REMW:
			dividend := int32(val1)
			divisor := int32(val2)

			if dividend == (-2147483647-1) && divisor == -1 {
				val = 0
			} else if divisor == 0 {
				val = int64(dividend)
			} else {
				val = int64(dividend % divisor)
			}

		case RISBEE_OPINST_RT32_REMUW:
			dividend := uint32(val1)
			divisor := uint32(val2)

			if divisor == 0 {
				val = int64(dividend)
			} else {
				val = int64(dividend % divisor)
			}

		default:
			vm.panic("Invalid store doubleword instruction.")
		}

		if rd != 0 {
			vm.Registers[rd] = uint64(val)
		}

	case RISBEE_OPINST_LUI:
		immediate := int64(int32(inst & 0xFFFFF000))
		if rd != 0 {
			vm.Registers[rd] = uint64(immediate)
		}

	case RISBEE_OPINST_AUIPC:
		immediate := int64(int32(inst & 0xFFFFF000))
		if rd != 0 {
			vm.Registers[rd] = uint64(
				vm.Pc + uint64(immediate),
			)
		}

	case RISBEE_OPINST_JAL:
		imm20 := (inst >> 11) & 0x100000
		imm10_1 := (inst >> 20) & 0x7FE
		imm11 := (inst >> 9) & 0x800
		imm19_12 := inst & 0xFF000

		immediate := int64(int32((imm20|
			imm10_1|
			imm11|
			imm19_12)<<11) >> 11)

		if rd != 0 {
			vm.Registers[rd] = vm.Pc + 4
		}

		vm.Pc = vm.Pc + uint64(immediate)
		return

	case RISBEE_OPINST_JALR:
		immediate := int64(int32(inst&0xFFF00000) >> 20)
		pc := vm.Pc + 4

		vm.Pc = uint64(int64(
			vm.Registers[rs1]+uint64(immediate)) & -2,
		)

		if rd != 0 {
			vm.Registers[rd] = uint64(pc)
		}
		return

	case RISBEE_OPINST_BRANCH:
		functionCode3 := (inst >> 12) & 0x7

		imm12 := (inst >> 19) & 0x1000
		imm10_5 := (inst >> 20) & 0x7E0
		imm4_1 := (inst >> 7) & 0x1E
		imm11 := (inst << 4) & 0x800

		immediate := int64(int32((imm12|
			imm10_5|
			imm4_1|
			imm11)<<19) >> 19)

		val1 := vm.Registers[rs1]
		val2 := vm.Registers[rs2]

		var condition bool
		switch functionCode3 {
		case RISBEE_FC3_BEQ:
			condition = (val1 == val2)

		case RISBEE_FC3_BNE:
			condition = (val1 != val2)

		case RISBEE_FC3_BLT:
			condition = (int64(val1) < int64(val2))

		case RISBEE_FC3_BGE:
			condition = (int64(val1) >= int64(val2))

		case RISBEE_FC3_BLTU:
			condition = (val1 < val2)

		case RISBEE_FC3_BGEU:
			condition = (val1 >= val2)

		default:
			vm.panic("Invalid branch instruction.")
		}

		if condition {
			vm.Pc = vm.Pc + uint64(immediate)
			return
		}

	case RISBEE_OPINST_FENCE:
		// No-op for now (memory ordering not needed temporarily)

	case RISBEE_OPINST_CALL:
		functionCode11 := (inst >> 20) & 0xFFF

		switch functionCode11 {
		case 0x0:
			code := vm.Registers[17]
			vm.Registers[10] = vm.handleSyscall(code)

		case 0x1:
			vm.ExitCode = -1
			vm.Running = false

		default:
			vm.panic("Invalid system instruction.")
		}

	default:
		vm.panic("Invalid opcode instruction.")
	}

	vm.Pc += 4
}
