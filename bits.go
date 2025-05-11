/*
 * Copyright 2025 Nathanne Isip
 * This file is part of Risbee (https://github.com/nthnn/risbee)
 * This code is licensed under MIT license (see LICENSE for details)
 */

package risbee

// uint16LittleEndian reads a 2-byte slice b in
// little-endian order and returns the corresponding
// uint16.
func uint16LittleEndian(b []byte) uint16 {
	if len(b) < 2 {
		return 0
	}

	return uint16(b[0]) |
		uint16(b[1])<<8
}

// uint32LittleEndian reads a 4-byte slice b in
// little-endian order and returns the corresponding
// uint32.
func uint32LittleEndian(b []byte) uint32 {
	if len(b) < 4 {
		return 0
	}

	return uint32(b[0]) |
		uint32(b[1])<<8 |
		uint32(b[2])<<16 |
		uint32(b[3])<<24
}

// uint64LittleEndian reads an 8-byte slice b
// in little-endian order and returns the
// corresponding uint64.
func uint64LittleEndian(b []byte) uint64 {
	if len(b) < 8 {
		return 0
	}

	return uint64(b[0]) |
		uint64(b[1])<<8 |
		uint64(b[2])<<16 |
		uint64(b[3])<<24 |
		uint64(b[4])<<32 |
		uint64(b[5])<<40 |
		uint64(b[6])<<48 |
		uint64(b[7])<<56
}

// putUint16 writes the uint16 v into b in
// little-endian order.
func putUint16(b []byte, v uint16) {
	if len(b) < 2 {
		return
	}

	b[0] = byte(v)
	b[1] = byte(v >> 8)
}

// putUint32 writes the uint32 v into b
// in little-endian order.
func putUint32(b []byte, v uint32) {
	if len(b) < 4 {
		return
	}

	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
}

// putUint64 writes the uint64 v into b
// in little-endian order.
func putUint64(b []byte, v uint64) {
	if len(b) < 8 {
		return
	}

	b[0] = byte(v)
	b[1] = byte(v >> 8)
	b[2] = byte(v >> 16)
	b[3] = byte(v >> 24)
	b[4] = byte(v >> 32)
	b[5] = byte(v >> 40)
	b[6] = byte(v >> 48)
	b[7] = byte(v >> 56)
}
