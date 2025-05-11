/*
 * Copyright 2025 Nathanne Isip
 * This file is part of Risbee (https://github.com/nthnn/risbee)
 * This code is licensed under MIT license (see LICENSE for details)
 */

package risbee

// Performs left shift operation on a 64-bit signed integer.
//
// This function performs a left shift operation on the 64-bit signed integer
// `a` by the number of bits specified by `b`. It returns the result of the
// left shift operation.
//
// Parameters:
// - a The 64-bit signed integer value to be shifted.
// - b The number of bits to shift `a` by.
//
// Returns the result of the left shift operation.
func shiftLeftInt64(
	x int64,
	y int64,
) int64 {
	if y >= 0 && y < 64 {
		return int64(uint64(x) << y)
	} else if y < 0 && y > -64 {
		return int64(uint64(x) >> -y)
	}

	return 0
}

// Performs right shift operation on a 64-bit signed integer.
//
// This function performs a right shift operation on the 64-bit signed integer
// a by the number of bits specified by b. It returns the result of the
// right shift operation.
//
// Parameters:
// - a The 64-bit signed integer value to be shifted.
// - b The number of bits to shift a by.
//
// Returns the result of the right shift operation.
func shiftRightInt64(
	x int64,
	y int64,
) int64 {
	if y >= 0 && y < 64 {
		return int64(uint64(x) >> y)
	} else if y < 0 && y > -64 {
		return int64(uint64(x) << -y)
	}

	return 0
}

// Performs right shift operation on a 128-bit signed integer.
//
// This function performs a right shift operation on the lower 64 bits of the
// 128-bit signed integer a by the number of bits specified by b. It
// returns the result of the right shift operation.
//
// Parameters:
// - a The 128-bit signed integer value to be shifted.
// - b The number of bits to shift a by.
//
// Returns the result of the right shift operation.
func shiftRightInt128(
	x int64,
	y int64,
) int64 {
	if y >= 0 && y < 128 {
		return int64(uint64(x) >> y)
	} else if y < 0 && y > -128 {
		return int64(uint64(x) << -y)
	}

	return 0
}

// Performs arithmetic right shift operation on a 64-bit signed integer.
//
// This function performs an arithmetic right shift operation on the 64-bit
// signed integer a by the number of bits specified by b. It returns the
// result of the arithmetic right shift operation.
//
// Parameters:
// - a The 64-bit signed integer value to be shifted.
// - b The number of bits to shift a by.
//
// Returns the result of the arithmetic right shift operation.
func arithShiftRightInt64(
	x int64,
	y int64,
) int64 {
	if y >= 0 && y < 64 {
		return x >> y
	} else if y >= 64 {
		if x < 0 {
			return -1
		}

		return 0
	} else if y < 0 && y > -64 {
		return x << -y
	}

	return 0
}
