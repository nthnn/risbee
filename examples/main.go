/*
 * Copyright 2025 Nathanne Isip
 * This file is part of Risbee (https://github.com/nthnn/risbee)
 * This code is licensed under MIT license (see LICENSE for details)
 */

// Example main.go demonstrating VM initialization, syscall
// registration, file loading, and execution.
//
// This program accepts a RISC-V binary filename as a command-line
// argument, sets up a simple print syscall (code 1), loads the
// binary into the VM’s memory, and runs it. It prints usage
// instructions and exits with code 1 if no filename is provided,
// or an error message if loading fails.
//
// Usage:
//
//	go run main.go <riscv-binary>
//
// Breakdown:
//  1. Argument Check: Ensures a filename argument is passed;
//     otherwise, prints usage and exits.
//  2. VM Creation: Instantiates a new RisbeeVm and initializes
//     its state (memory, PC, syscalls).
//  3. Syscall Registration: Registers syscall code 1 to print
//     a null-terminated string from VM memory.
//     - Retrieves string pointer via GetPointerParam(0), reads
//     string with GetStringPointer, prints it.
//  4. File Loading: Uses LoadFile to read the binary into VM
//     memory at offset 0x1000;
//     on failure, prints an error and exits.
//  5. Execution: Calls Run(), entering the fetch-execute loop
//     until the program calls exit.
package main

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
