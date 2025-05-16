<p align="center">
    <img src="assets/risbee-logo.png" width="175" />
</p>
<h1 align="center">Risbee</h1>

![Build CI](https://github.com/nthnn/risbee/actions/workflows/build_ci.yml/badge.svg)
[![License](https://img.shields.io/badge/license-MIT-blue)](https://opensource.org/license/mit)

Risbee is a small, self-contained virtual machine that draws inspiration from the RISC-V instruction set. Built in Go, it offers a simple and approachable way to experiment with low-level concepts like registers, memory management, and instruction decoding. With just few kilobytes of byte-addressable memory and 32 general-purpose registers, Risbee keeps its footprint minimal while still supporting a rich subset of operations—everything from basic loads and stores to arithmetic, branching, and even custom syscalls.

You don’t need to be a seasoned compiler engineer or hardware expert to get started. Risbee’s design emphasizes clarity: the VM initializes with a clear default state, lets you load a binary or raw byte slice into memory at a fixed offset, and then runs a straightforward fetch-decode-execute loop. If you want to print text, handle I/O, or integrate with your Go application in other ways, you simply register a Go function as a syscall handler, and Risbee will invoke it whenever your code calls an environment call instruction.

All you need is a Go workspace and a compiled RISC-V binary (or a simple byte array), and you’re ready to explore how a tiny VM brings machine instructions to life.

## Supported Features

- **Fetch-Decode-Execute Loop**: Continuously fetches 32-bit little-endian instructions, decodes them by opcode/function-codes, and executes until halted.
- **Instruction Support**
    - **Loads** (`LB`, `LH`, `LW`, `LD`, and unsigned variants)
    - **Stores** (`SB`, `SH`, `SW`, and `SD`)
    - **Immediate ALU** (`ADDI`, `SLTI`, `XORI`, `ORI`, `ANDI`, and shifts)
    - **Register-Register ALU** (32/64-bit adds, subs, shifts, multiplies, divides, remainders)
    - **Control Flow** (`BEQ`, `BNE`, `BLT`, `BGE`, `BLTU`, `BGEU`, `JAL`, `JALR`)
    - **Fences** (no-op placeholder)
    - **Syscalls** (via `CALL`/`ECALL`)
- **Syscall API**
    - Register handlers with `SetSystemCall(code, fn)`.
    - Retrieve string and pointer parameters with `GetStringPointer` and `GetPointerParam`.
    - Built-in exit syscall (`code 0` uses R10 for status).
- **Memory & Registers**
    - Dynamic kilobytes contiguous memory (`[]byte`)
    - 32 × 64-bit registers (R0 read-only zero)
    - Program Counter initialized to `0x1000`
    - Stack Pointer (`R2`) auto-set to top of memory on load
- **Error Handling**: Invalid instructions or syscalls trigger `panic()`, printing an error, setting exit code to `-1`, and halting.

## Installation

To incorporate Risbee into your Go project, use go get:

```go
# Install via go get
go get github.com/nthnn/risbee
```

## Quick Start

This program accepts a RISC-V binary filename as a command-line argument, sets up a simple print syscall (code 1), loads the binary into the VM’s memory, and runs it. It prints usage instructions and exits with code 1 if no filename is provided, or an error message if loading fails.

```go
import (
	"fmt"
	"os"

	"github.com/nthnn/risbee"
)

func main() {
	// Ensure a filename is provided as an argument.
	if len(os.Args) < 2 {
		fmt.Println("Usage: risbee <filename>")
		os.Exit(1)
	}

	// Create and initialize the VM.
	vm := &risbee.RisbeeVm{}
	vm.Initialize()

	// Register a simple "print" syscall at code 1.
	// When invoked, it reads a string pointer from a0,
	// prints the string, and returns.
	vm.SetSystemCall(1, func(vm *risbee.RisbeeVm) uint64 {
		ptr := vm.GetPointerParam(0)

		fmt.Print(vm.GetStringPointer(ptr))
		return ptr
	})

	// Load the RISC-V binary into VM memory; exit on failure.
	if !vm.LoadFile(os.Args[1]) {
		fmt.Println("Failed to load file:", os.Args[1])
		os.Exit(1)
	}

	// Execute the loaded program.
	vm.Run()
}

```

### Key Methods

- `Initialize()`: Reset PC, exit code, running flag, and syscall table.
- `LoadFromBytes(data []byte) bool`: Load raw bytes at 0x1000 into VM memory.
- `SetSystemCall(code uint64, fn RisbeeVmSyscallFn)`: Register a syscall handler.
- `GetPointerParam(idx uint64) uint64`: Read syscall argument from a0+idx.
- `GetStringPointer(ptr uint64) string`: Read null-terminated string from VM memory.
- `Run()`: Enter the fetch-decode-execute loop.
- `Stop()`: Halt execution.
- `GetExitCode() int`: Retrieve VM exit status.

## Memory Layout

The VM reserves the first 4 KiB (0x0000–0x0FFF) as a “reserved” region that you can use for static data, heap, or simply leave untouched. At address 0x1000, the VM begins loading your program image, and everything from 0x1000 up to the end of the space is available for code, global variables, stack, and heap allocations. By convention, the stack pointer (register R2) is initialized to the very top of memory (0x10000), allowing your program to grow the stack downward into the unused upper region.

```
0x0000 ─────────────────────────── Reserved (data/heap)
   │
0x0FFF ────────────────────────────
0x1000 ─────────────────────────── Program image load address
   │   • Text segment (.text)
   │   • Read-only data (.rodata)
   │   • Initialized data (.data)
   │   • Uninitialized data (.bss)
   │   • Heap (grows upward)
   │   • Stack (grows downward, SP=R2 starts at 0x10000)
  ...  ─────────────────────────── End of VM memory
```

- **Load Offset (`0x1000`)**: All binaries and raw byte slices are copied here.
- **Stack Pointer (`R2`)**: Set to `0x10000` on load, so your stack grows downward into fresh memory.
- **Heap**: Begins immediately after any static data; you can manage it entirely in software.

This fixed, contiguous layout keeps things simple and predictable, letting you focus on instruction semantics rather than complex memory mapping.

## Custom SDKs

If you’re building applications that rely on your own syscall conventions, writing a small “SDK” (set of helper functions or libraries) can make life much easier. A custom SDK provides idiomatic wrappers around raw environment calls, handles parameter marshalling, and can expose higher-level APIs to your Go or C programs. Below is a guide to developing SDKs for Risbee syscalls.

On the example usage above, the system call with `0x01` address which is for printing string, is being declared as such below on the source file of binary (see next section).

```c
static inline long prints(const char* str) {
    register long a0 asm("a0") = (long) str;
    register long scid asm("a7") = 1;

    asm volatile ("scall" : "+r"(a0) : "r"(scid));
    return a0;
}
```

Hence, it can be used as shown below:

```c
int main() {
    prints("Hello, world!\r\n");
    return 0;
}
```

## Building Program Binaries

To turn your C (or C++) and assembly sources into a raw RISC-V binary suitable for the Risbee VM, follow these steps:

1. **Compile and Link**: Use the RISC-V GCC toolchain targeting RV64IM, no standard library, and your custom linker script.

    ```sh
    riscv64-unknown-elf-g++     \
        -march=rv64im           \
        -mabi=lp64              \
        -nostdlib               \
        -Wl,-T,scripts/link.ld  \
        -O2 -o main.out         \
        scripts/launcher.s      \
        main.c
    ```

    - `-march=rv64im` selects 64-bit integer RISC-V with multiplication.
    - `-mabi=lp64` chooses the 64-bit ABI.
    - `-nostdlib` prevents linking against the host’s C runtime.
    - `-T link.ld` tells the linker to use your memory layout.

2. **Extract the Raw Image**: Convert the ELF output into a flat binary blob that Risbee can load.

    ```sh
    riscv64-unknown-elf-objcopy -O binary main.out main.bin
    ```

    This strips headers and relocations, leaving just the machine code and data laid out at the offsets defined by `link.ld`.

## License

```
Copyright 2025 Nathanne Isip

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```
