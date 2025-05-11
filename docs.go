/*
Package risbee implements a lightweight RISC-V inspired virtual machine (VM),
designed for educational purposes, embedded experimentation, and syscall
integration. The VM exposes a fixed-size byte-addressable memory arena,
public general-purpose registers, program counter tracking, and a simple
syscall dispatch mechanism. It supports a subset of RISC-V instruction
formats including loads, stores, immediate arithmetic, register-register
operations (both 32- and 64-bit variants), control flow (branches, jumps),
and environment calls (syscalls).

Key Features:
  - Fixed Memory: 64 KiB of memory for program code and data, with a
default code load offset at 0x1000 (4096).
  - 32 General-Purpose Registers: 64-bit registers R0–R31, with R0
hardwired to zero (writes ignored).
  - Program Counter (PC): 64-bit PC initialized to 0x1000; increments
automatically or is modified by control-transfer instructions.
  - Syscall Integration: Register-based syscall interface via ECALL
instructions, supporting custom registration of handlers.
  - Exit Handling: Built-in exit code propagation and graceful shutdown.

Usage Overview:
  1. Instantiate RisbeeVm and call Initialize() to set up PC, exit code,
and syscall table.
  2. Load a RISC-V binary into VM memory at the load offset (4096) using
LoadFile(fileName); on success, the stack pointer (R2) is set to memory top.
  3. Optionally register custom syscalls: SetSystemCall(addr, handler).
  4. Execute the program with Run(), which fetches and executes
instructions until StopVM() is called or an exit syscall occurs.
  5. Retrieve the exit status via GetExitCode().

Memory Layout:
  - 0x0000–0x0FFF: Reserved or available for data.
  - 0x1000–0xFFFF: Program image, heap, and stack.

Syscall Mechanism:
  - ECALL encoded in the CALL instruction family; a0 (R10) holds the syscall code.
  - Syscalls are looked up in SysCalls map; code 0 triggers VM exit with status in a0.

Error Handling:
  - Invalid instructions or memory accesses invoke panic(), printing an error,
stopping the VM, and setting exit code to -1.

Constants and Types:
  - RISBEE_STACK_SIZE: Total VM memory size (64 KiB).
  - RisbeeVmSyscallFn: Callback signature for syscall handlers.
  - RisbeeVm: Core struct encapsulating VM state, memory, registers, PC, and syscalls.
*/

package risbee
