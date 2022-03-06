// Copyright 2022  Il Sub Bang
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package h3go

import (
	"math"
	"strconv"
)

type H3Index uint64

// define's of constants for bitwise manipulation of H3Index's.
const (
	// The number of bits in an H3 index.
	H3_NUM_BITS = 64

	// The bit offset of the max resolution digit in an H3 index.
	H3_MAX_OFFSET = 63

	// The bit offset of the mode in an H3 index.
	H3_MODE_OFFSET = 59

	// The bit offset of the base cell in an H3 index.
	H3_BC_OFFSET = 45

	// The bit offset of the resolution in an H3 index.
	H3_RES_OFFSET = 52

	// The bit offset of the reserved bits in an H3 index.
	H3_RESERVED_OFFSET = 56

	// The number of bits in a single H3 resolution digit.
	H3_PER_DIGIT_OFFSET = 3

	// 1 in the highest bit, 0's everywhere else.
	H3_HIGH_BIT_MASK = uint64(1) << H3_MAX_OFFSET

	// 0 in the highest bit, 1's everywhere else.
	H3_HIGH_BIT_MASK_NEGATIVE = ^H3_HIGH_BIT_MASK

	// 1's in the 4 mode bits, 0's everywhere else.
	H3_MODE_MASK = uint64(15) << H3_MODE_OFFSET

	// 0's in the 4 mode bits, 1's everywhere else.
	H3_MODE_MASK_NEGATIVE = ^H3_MODE_MASK

	// 1's in the 7 base cell bits, 0's everywhere else.
	H3_BC_MASK = uint64(127) << H3_BC_OFFSET

	// 0's in the 7 base cell bits, 1's everywhere else.
	H3_BC_MASK_NEGATIVE = ^H3_BC_MASK

	// 1's in the 4 resolution bits, 0's everywhere else.
	H3_RES_MASK = uint64(15) << H3_RES_OFFSET

	// 0's in the 4 resolution bits, 1's everywhere else.
	H3_RES_MASK_NEGATIVE = ^H3_RES_MASK

	// 1's in the 3 reserved bits, 0's everywhere else.
	H3_RESERVED_MASK = uint64(7) << H3_RESERVED_OFFSET

	// 0's in the 3 reserved bits, 1's everywhere else.
	H3_RESERVED_MASK_NEGATIVE = ^H3_RESERVED_MASK

	// 1's in the 3 bits of res 15 digit bits, 0's everywhere else.
	H3_DIGIT_MASK = uint64(7)

	// 0's in the 7 base cell bits, 1's everywhere else.
	H3_DIGIT_MASK_NEGATIVE = ^H3_DIGIT_MASK
)

// H3 index with mode 0, res 0, base cell 0, and 7 for all index digits.
// Typically used to initialize the creation of an H3 cell index, which
// expects all direction digits to be 7 beyond the cell's resolution.
const H3_INIT = H3Index(35184372088831)

// Invalid index used to indicate an error from geoToH3 and related functions
// or missing data in arrays of h3 indices. Analogous to NaN in floating point.
const H3_NULL = H3Index(0)

// H3_GET_HIGH_BIT gets the highest bit of the H3 index.
//
// Deprecated: Use (H3Index).GetHighBit instead.
func H3_GET_HIGH_BIT(h3 H3Index) int {
	return int((uint64(h3) & H3_HIGH_BIT_MASK) >> H3_MAX_OFFSET)
}

/* ========================================================================== */

// GetHighBit gets the highest bit of the H3 index.
func (h3 H3Index) GetHighBit() int {
	return int((uint64(h3) & H3_HIGH_BIT_MASK) >> H3_MAX_OFFSET)
}

// H3_SET_HIGH_BIT sets the highest bit of the h3 to v.
//
// Deprecated: Use (*H3Index).SetHighBit instead.
func H3_SET_HIGH_BIT(h3 *H3Index, v int) {
	*h3 = H3Index((uint64(*h3) & H3_HIGH_BIT_MASK_NEGATIVE) | ((uint64(v)) << H3_MAX_OFFSET))
}

// SetHighBit sets the highest bit of the h3 to v.
func (h3 *H3Index) SetHighBit(v int) {
	*h3 = H3Index((uint64(*h3) & H3_HIGH_BIT_MASK_NEGATIVE) | ((uint64(v)) << H3_MAX_OFFSET))
}

// H3_GET_MODE gets the integer mode of h3.
//
// Deprecated: Use (H3Index).GetMode instead.
func H3_GET_MODE(h3 H3Index) int {
	return int((uint64(h3) & H3_MODE_MASK) >> H3_MODE_OFFSET)
}

// GetMode gets the integer mode of h3.
func (h3 H3Index) GetMode() int {
	return int((uint64(h3) & H3_MODE_MASK) >> H3_MODE_OFFSET)
}

// H3_SET_MODE sets the integer mode of h3 to v.
//
// Deprecated: Use (*H3Index).SetMode instead.
func H3_SET_MODE(h3 *H3Index, v int) {
	*h3 = H3Index((uint64(*h3) & H3_MODE_MASK_NEGATIVE) | (uint64(v) << H3_MODE_OFFSET))
}

// SetMode sets the integer mode of h3 to v.
func (h3 *H3Index) SetMode(v int) {
	*h3 = H3Index((uint64(*h3) & H3_MODE_MASK_NEGATIVE) | (uint64(v) << H3_MODE_OFFSET))
}

// H3_GET_BASE_CELL gets the integer base cell of h3.
//
// Deprecated: Use (H3Index).GetBaseCell instead.
func H3_GET_BASE_CELL(h3 H3Index) int {
	return int((uint64(h3) & H3_BC_MASK) >> H3_BC_OFFSET)
}

// GetBaseCell gets the integer base cell of h3.
func (h3 H3Index) GetBaseCell() int {
	return int((uint64(h3) & H3_BC_MASK) >> H3_BC_OFFSET)
}

// H3_SET_BASE_CELL sets the integer base cell of h3 to bc.
//
// Deprecated: Use (*H3Index).SetBaseCell instead.
func H3_SET_BASE_CELL(h3 *H3Index, bc int) {
	*h3 = H3Index((uint64(*h3) & H3_BC_MASK_NEGATIVE) | (uint64(bc) << H3_BC_OFFSET))
}

// SetBaseCell sets the integer base cell of h3 to bc.
func (h3 *H3Index) SetBaseCell(bc int) {
	*h3 = H3Index((uint64(*h3) & H3_BC_MASK_NEGATIVE) | (uint64(bc) << H3_BC_OFFSET))
}

// H3_GET_RESOLUTION gets the integer resolution of h3.
//
// Deprecated: Use (H3Index).GetResolution instead.
func H3_GET_RESOLUTION(h3 H3Index) int {
	return int((uint64(h3) & H3_RES_MASK) >> H3_RES_OFFSET)
}

// GetResolution gets the integer resolution of h3.
func (h3 H3Index) GetResolution() int {
	return int((uint64(h3) & H3_RES_MASK) >> H3_RES_OFFSET)
}

// H3_SET_RESOLUTION sets the integer resolution of h3.
//
// Deprecated: Use (*H3Index).SetResolution instead.
func H3_SET_RESOLUTION(h3 *H3Index, res int) {
	*h3 = H3Index((uint64(*h3) & H3_RES_MASK_NEGATIVE) | (uint64(res) << H3_RES_OFFSET))
}

// SetResolution sets the integer resolution of h3.
func (h3 *H3Index) SetResolution(res int) {
	*h3 = H3Index((uint64(*h3) & H3_RES_MASK_NEGATIVE) | (uint64(res) << H3_RES_OFFSET))
}

// H3_GET_RESERVED_BITS gets a value in the reserved space. Should always be zero for valid indexes.
//
// Deprecated: Use (H3Index).GetReservedBits instead.
func H3_GET_RESERVED_BITS(h3 H3Index) int {
	return int((uint64(h3) & H3_RESERVED_MASK) >> H3_RESERVED_OFFSET)
}

// GetReservedBits gets a value in the reserved space. Should always be zero for valid indexes.
func (h3 H3Index) GetReservedBits() int {
	return int((uint64(h3) & H3_RESERVED_MASK) >> H3_RESERVED_OFFSET)
}

// H3_SET_RESERVED_BITS sets a value in the reserved space. Setting to non-zero
// may produce invalid indexes.
//
// Deprecated: Use (*H3Index).SetReservedBits instead.
func H3_SET_RESERVED_BITS(h3 *H3Index, v int) {
	*h3 = H3Index((uint64(*h3) & H3_RESERVED_MASK_NEGATIVE) | (uint64(v) << H3_RESERVED_OFFSET))
}

// SetReservedBits sets a value in the reserved space. Setting to non-zero
// may produce invalid indexes.
func (h3 *H3Index) SetReservedBits(v int) {
	*h3 = H3Index((uint64(*h3) & H3_RESERVED_MASK_NEGATIVE) | (uint64(v) << H3_RESERVED_OFFSET))
}

// H3_GET_INDEX_DIGIT gets the resolution res integer digit (0-7) of h3.
//
// Deprecated: Use (H3Index).GetIndexDigit instead.
func H3_GET_INDEX_DIGIT(h3 H3Index, res int) Direction {
	resDigit := (MAX_H3_RES - res) * H3_PER_DIGIT_OFFSET

	return Direction((uint64(h3) >> resDigit) & H3_DIGIT_MASK)
}

// GetIndexDigit gets the resolution res integer digit (0-7) of h3.
func (h3 H3Index) GetIndexDigit(res int) Direction {
	resDigit := (MAX_H3_RES - res) * H3_PER_DIGIT_OFFSET

	return Direction((uint64(h3) >> resDigit) & H3_DIGIT_MASK)
}

// H3_SET_INDEX_DIGIT sets the resolution res digit of h3 to the integer digit (0-7)
//
// Deprecated: Use (*H3Index).SetIndexDigit instead.
func H3_SET_INDEX_DIGIT(h3 *H3Index, res int, digit Direction) {
	resDigit := (MAX_H3_RES - res) * H3_PER_DIGIT_OFFSET

	*h3 = H3Index((uint64(*h3) & ^(H3_DIGIT_MASK << resDigit)) |
		(uint64(digit) << resDigit))
}

// SetIndexDigit sets the resolution res digit of h3 to the integer digit (0-7)
func (h3 *H3Index) SetIndexDigit(res int, digit Direction) {
	resDigit := (MAX_H3_RES - res) * H3_PER_DIGIT_OFFSET

	*h3 = H3Index((uint64(*h3) & ^(H3_DIGIT_MASK << resDigit)) |
		(uint64(digit) << resDigit))
}

// Return codes for compact
const (
	COMPACT_SUCCESS       = 0
	COMPACT_LOOP_EXCEEDED = -1
	COMPACT_DUPLICATE     = -2
	COMPACT_ALLOC_FAILED  = -3
)

// H3GetResolution returns the H3 resolution of an H3 index.
//
// Deprecated: Use (H3Index).GetResolution instead.
func H3GetResolution(h H3Index) int { return H3_GET_RESOLUTION(h) }

// H3GetBaseCell returns the H3 base cell "number" of an H3 cell (hexagon or pentagon).
//
// Note: Technically works on H3 edges, but will return base cell of the
// origin cell.
//
// Deprecated: Use (H3Index).GetBaseCell instead.
func H3GetBaseCell(h H3Index) int { return H3_GET_BASE_CELL(h) }

// StringToH3 converts a string representation of an H3 index into an H3 index.
//
// Return The H3 index corresponding to the string argument, or H3_NULL if
// invalid.
func StringToH3(str string) H3Index {
	// If failed, h will be unmodified and we should return H3_NULL anyways.
	u64, err := strconv.ParseUint(str, 16, 64)
	if err != nil {
		return H3_NULL
	}
	return H3Index(u64)
}

// H3ToString converts an H3 index into a string representation.
//
// Deprecated: Use (H3Index).String instead.
func H3ToString(h H3Index) string {
	return strconv.FormatUint(uint64(h), 16)
}

// String converts an H3 index into a string representation.
func (h3 H3Index) String() string {
	return strconv.FormatUint(uint64(h3), 16)
}

// H3IsValid returns whether or not an H3 index is a valid cell (hexagon or
// pentagon).
//
// Return true if the H3 index if valid, and false if it is not.
//
// Deprecated: Use (H3Index).IsValid instead.
func H3IsValid(h H3Index) bool {
	if H3_GET_HIGH_BIT(h) != 0 {
		return false
	}

	if H3_GET_MODE(h) != H3_HEXAGON_MODE {
		return false
	}

	if H3_GET_RESERVED_BITS(h) != 0 {
		return false
	}

	baseCell := H3_GET_BASE_CELL(h)
	if baseCell < 0 || baseCell >= NUM_BASE_CELLS {
		return false
	}

	res := H3_GET_RESOLUTION(h)
	if res < 0 || res > MAX_H3_RES {
		return false
	}

	foundFirstNonZeroDigit := false
	for r := 1; r <= res; r++ {
		digit := H3_GET_INDEX_DIGIT(h, r)

		if !foundFirstNonZeroDigit && digit != CENTER_DIGIT {
			foundFirstNonZeroDigit = true
			if _isBaseCellPentagon(baseCell) && digit == K_AXES_DIGIT {
				return false
			}
		}

		if digit < CENTER_DIGIT || digit >= Direction(NUM_DIGITS) {
			return false
		}
	}

	for r := res + 1; r <= MAX_H3_RES; r++ {
		digit := H3_GET_INDEX_DIGIT(h, r)
		if digit != INVALID_DIGIT {
			return false
		}
	}

	return true
}

// IsValid returns whether or not an H3 index is a valid cell (hexagon or
// pentagon).
//
// Return true if the H3 index if valid, and false if it is not.
func (h3 H3Index) IsValid() bool {
	if H3_GET_HIGH_BIT(h3) != 0 {
		return false
	}

	if H3_GET_MODE(h3) != H3_HEXAGON_MODE {
		return false
	}

	if H3_GET_RESERVED_BITS(h3) != 0 {
		return false
	}

	baseCell := H3_GET_BASE_CELL(h3)
	if baseCell < 0 || baseCell >= NUM_BASE_CELLS {
		return false
	}

	res := H3_GET_RESOLUTION(h3)
	if res < 0 || res > MAX_H3_RES {
		return false
	}

	foundFirstNonZeroDigit := false
	for r := 1; r <= res; r++ {
		digit := H3_GET_INDEX_DIGIT(h3, r)

		if !foundFirstNonZeroDigit && digit != CENTER_DIGIT {
			foundFirstNonZeroDigit = true
			if _isBaseCellPentagon(baseCell) && digit == K_AXES_DIGIT {
				return false
			}
		}

		if digit < CENTER_DIGIT || digit >= Direction(NUM_DIGITS) {
			return false
		}
	}

	for r := res + 1; r <= MAX_H3_RES; r++ {
		digit := H3_GET_INDEX_DIGIT(h3, r)
		if digit != INVALID_DIGIT {
			return false
		}
	}

	return true
}

// setH3Index initializes an H3 index.
//
// Deprecated: Use _setH3Index instead.
func setH3Index(hp *H3Index, res int, baseCell int, initDigit Direction) {
	h := H3_INIT
	H3_SET_MODE(&h, H3_HEXAGON_MODE)
	H3_SET_RESOLUTION(&h, res)
	H3_SET_BASE_CELL(&h, baseCell)
	for r := 1; r <= res; r++ {
		H3_SET_INDEX_DIGIT(&h, r, initDigit)
	}
	*hp = h
}

// _setH3Index initializes an H3 index.
func _setH3Index(res int, baseCell int, initDigit Direction) H3Index {
	h := H3_INIT
	H3_SET_MODE(&h, H3_HEXAGON_MODE)
	H3_SET_RESOLUTION(&h, res)
	H3_SET_BASE_CELL(&h, baseCell)
	for r := 1; r <= res; r++ {
		H3_SET_INDEX_DIGIT(&h, r, initDigit)
	}
	return h
}

// H3ToParent produces the parent index for a given H3 index
//
// Return H3Index of the parent, or H3_NULL if you actually asked for a child
//
// Deprecated: Use (H3Index).ToParent instead.
func H3ToParent(h H3Index, parentRes int) H3Index {
	childRes := H3_GET_RESOLUTION(h)
	if parentRes > childRes {
		return H3_NULL
	} else if parentRes == childRes {
		return h
	} else if parentRes < 0 || parentRes > MAX_H3_RES {
		return H3_NULL
	}

	parentH := h
	H3_SET_RESOLUTION(&parentH, parentRes)
	for i := parentRes + 1; i <= childRes; i++ {
		H3_SET_INDEX_DIGIT(&parentH, i, Direction(H3_DIGIT_MASK))
	}
	return parentH
}

// ToParent produces the parent index for a given H3 index
//
// Return H3Index of the parent, or H3_NULL if you actually asked for a child
func (h3 H3Index) ToParent(parentRes int) H3Index {
	childRes := H3_GET_RESOLUTION(h3)
	if parentRes > childRes {
		return H3_NULL
	} else if parentRes == childRes {
		return h3
	} else if parentRes < 0 || parentRes > MAX_H3_RES {
		return H3_NULL
	}

	parentH := h3
	H3_SET_RESOLUTION(&parentH, parentRes)
	for i := parentRes + 1; i <= childRes; i++ {
		H3_SET_INDEX_DIGIT(&parentH, i, Direction(H3_DIGIT_MASK))
	}
	return parentH
}

// _isValidChildRes determines whether one resolution is a valid child
// resolution of another. Each resolution is considered a valid child resolution
// of itself.
//
// Return The validity of the child resolution.
func _isValidChildRes(parentRes int, childRes int) bool {
	if childRes < parentRes || childRes > MAX_H3_RES {
		return false
	}
	return true
}

// MaxH3ToChildrenSize returns the maximum number of children possible for a
// given child level.
//
// Return int count of maximum number of children (equal for hexagons, less for
// pentagons.
func MaxH3ToChildrenSize(h H3Index, childRes int) int {
	parentRes := H3_GET_RESOLUTION(h)
	if !_isValidChildRes(parentRes, childRes) {
		return 0
	}
	return _ipow(7, childRes-parentRes)
}

// makeDirectChild takes an index and immediately returns the immediate child
// index based on the specified cell number. Bit operations only, could generate
// invalid indexes if not careful (deleted cell under a pentagon).
//
// Return The new H3Index for the child.
func makeDirectChild(h H3Index, cellNumber Direction) H3Index {
	childRes := H3_GET_RESOLUTION(h) + 1

	childH := h
	H3_SET_RESOLUTION(&childH, childRes)
	H3_SET_INDEX_DIGIT(&childH, childRes, cellNumber)
	return childH
}

// H3ToChildren takes the given hexagon id and generates all of the children
// at the specified resolution storing them into the provided memory pointer.
// It's assumed that maxH3ToChildrenSize was used to determine the allocation.
//
// Deprecated: Use (H3Index).ToChildren instead.
func H3ToChildren(h H3Index, childRes int, children *[]H3Index) {
	parentRes := H3_GET_RESOLUTION(h)
	if !_isValidChildRes(parentRes, childRes) {
		return
	} else if parentRes == childRes {
		*children = append(*children, h)
		return
	}

	isAPentagon := H3IsPentagon(h)
	for i := CENTER_DIGIT; i < 7; i++ {
		if isAPentagon && i == K_AXES_DIGIT {
			continue
		}

		H3ToChildren(makeDirectChild(h, i), childRes, children)
	}
}

// ToChildren takes the given hexagon id and generates all of the children
// at the specified resolution.
//
// TODO: enhance algorithm
func (h3 H3Index) ToChildren(childRes int) []H3Index {
	buffer := make([]H3Index, 0, MaxH3ToChildrenSize(h3, childRes))
	H3ToChildren(h3, childRes, &buffer)
	return buffer
}

// H3ToCenterChild produces the center child index for a given H3 index at
// the specified resolution.
//
// Return H3Index of the center child, or H3_NULL if you actually asked for a
// parent.
//
// Deprecated: Use (H3Index).ToCenterChild instead.
func H3ToCenterChild(h H3Index, childRes int) H3Index {
	parentRes := H3_GET_RESOLUTION(h)
	if !_isValidChildRes(parentRes, childRes) {
		return H3_NULL
	} else if childRes == parentRes {
		return h
	}

	child := h
	H3_SET_RESOLUTION(&child, childRes)
	for i := parentRes + 1; i <= childRes; i++ {
		H3_SET_INDEX_DIGIT(&child, i, 0)
	}
	return child
}

// ToCenterChild produces the center child index for a given H3 index at
// the specified resolution.
//
// Return H3Index of the center child, or H3_NULL if you actually asked for a
// parent.
func (h3 H3Index) ToCenterChild(childRes int) H3Index {
	parentRes := H3_GET_RESOLUTION(h3)
	if !_isValidChildRes(parentRes, childRes) {
		return H3_NULL
	} else if childRes == parentRes {
		return h3
	}

	child := h3
	H3_SET_RESOLUTION(&child, childRes)
	for i := parentRes + 1; i <= childRes; i++ {
		H3_SET_INDEX_DIGIT(&child, i, 0)
	}
	return child
}

// Compact takes a set of hexagons all at the same resolution and compresses
// them by pruning full child branches to the parent level. This is also done
// for all parents recursively to get the minimum number of hex addresses that
// perfectly cover the defined space.
//
// Return an error code on bad input data.
func Compact(h3Set []H3Index) ([]H3Index, error) {
	if len(h3Set) == 0 {
		return nil, nil
	}

	res := H3_GET_RESOLUTION(h3Set[0])
	if res == 0 {
		compacted := make([]H3Index, len(h3Set))
		copy(compacted, h3Set)
		return compacted, nil
	}

	result := make([]H3Index, 0, len(h3Set))
	remaining := make([]H3Index, len(h3Set))
	copy(remaining, h3Set)

	for len(remaining) > 0 {
		if len(remaining) < 6 {
			// cannot compact more. append and break
			result = append(result, remaining...)
			break
		}

		// map[cell]count
		compactable := make(map[H3Index]int, len(remaining))

		res := H3_GET_RESOLUTION(remaining[0])
		parentRes := res - 1

		// count parent cells
		for _, cell := range remaining {
			parent := H3ToParent(cell, parentRes)
			isPentagon := H3IsPentagon(parent)
			if _, ok := compactable[parent]; ok {
				compactable[parent]++
				if compactable[parent] > 7 {
					return nil, ErrCompactDuplicate
				}
			} else if isPentagon {
				// set 2 if cell is pentagon. it helps checking if dragonball is completed.
				compactable[parent] = 2
			} else {
				compactable[parent] = 1
			}
		}

		// append uncompactable cells into result and cleanup remaining
		for i, cell := range remaining {
			parent := H3ToParent(cell, parentRes)
			if compactable[parent] < 7 {
				result = append(result, cell)
			}
			remaining[i] = 0
		}
		remaining = remaining[:0]

		// move compactable cells to remaining
		for cell, count := range compactable {
			if count == 7 {
				remaining = append(remaining, cell)
			}
		}
	}

	return result, nil
}

// Uncompact takes a compressed set of hexagons and expands back to the original
// set of hexagons.
//
// Return ErrUncompactResExceeded if any hexagon is smaller than the output
// resolution.
func Uncompact(compactedSet []H3Index, res int) ([]H3Index, error) {
	maxSize, err := MaxUncompactSize(compactedSet, res)
	if err != nil {
		return nil, err
	}

	h3Set := make([]H3Index, 0, maxSize)

	for _, cell := range compactedSet {
		if cell == 0 {
			continue
		}

		if cell.GetResolution() == res {
			h3Set = append(h3Set, cell)
		} else {
			h3Set = append(h3Set, cell.ToChildren(res)...)
		}
	}

	return h3Set, nil
}

// MaxUncompactSize takes a compacted set of hexagons are provides an
// upper-bound estimate of the size of the uncompacted set of hexagons.
//
// Return The number of hexagons to allocate memory for, or a negative number
// if an error occurs.
func MaxUncompactSize(compactedSet []H3Index, res int) (int, error) {
	maxNumHexagons := 0
	for i := 0; i < len(compactedSet); i++ {
		if compactedSet[i] == 0 {
			continue
		}
		currentRes := H3_GET_RESOLUTION(compactedSet[i])
		if !_isValidChildRes(currentRes, res) {
			// Nonsensical. Abort.
			return 0, ErrUncompactResExceeded
		}
		if currentRes == res {
			maxNumHexagons++
		} else {
			// Bigger hexagon to reduce in size
			maxNumHexagons += MaxH3ToChildrenSize(compactedSet[i], res)
		}
	}
	return maxNumHexagons, nil
}

// H3IsResClassIII takes a hexagon ID and determines if it is in a Class III
// resolution (rotated versus the icosahedron and subject to shape distortion
// adding extra points on icosahedron edges, making them not true hexagons).
//
// Return true if the hexagon is class III, otherwise 0.
//
// Deprecated: Use (H3Index).IsResClassIII instead.
func H3IsResClassIII(h H3Index) bool {
	return H3_GET_RESOLUTION(h)%2 == 1
}

// IsResClassIII takes a hexagon ID and determines if it is in a Class III
// resolution (rotated versus the icosahedron and subject to shape distortion
// adding extra points on icosahedron edges, making them not true hexagons).
//
// Return true if the hexagon is class III, otherwise false.
func (h3 H3Index) IsResClassIII() bool {
	return H3_GET_RESOLUTION(h3)%2 == 1
}

// H3IsPentagon takes an H3Index and determines if it is actually a
// pentagon.
//
// Return true if it is a pentagon, otherwise false.
//
// Deprecated: Use (H3Index).IsPentagon instead.
func H3IsPentagon(h H3Index) bool {
	return _isBaseCellPentagon(H3_GET_BASE_CELL(h)) &&
		_h3LeadingNonZeroDigit(h) == CENTER_DIGIT
}

// IsPentagon takes an H3Index and determines if it is actually a
// pentagon.
//
// Return true if it is a pentagon, otherwise false.
func (h3 H3Index) IsPentagon() bool {
	return _isBaseCellPentagon(H3_GET_BASE_CELL(h3)) &&
		_h3LeadingNonZeroDigit(h3) == CENTER_DIGIT
}

// _h3LeadingNonZeroDigit returns the highest resolution non-zero digit in an
// H3Index.
func _h3LeadingNonZeroDigit(h H3Index) Direction {
	for r := 1; r <= H3_GET_RESOLUTION(h); r++ {
		if H3_GET_INDEX_DIGIT(h, r) > 1 {
			return H3_GET_INDEX_DIGIT(h, r)
		}
	}

	// if we're here it's all 0's
	return CENTER_DIGIT
}

// _h3RotatePent60ccw rotate an H3Index 60 degrees counter-clockwise about a
// pentagonal center.
func _h3RotatePent60ccw(h H3Index) H3Index {
	// rotate in place; skips any leading 1 digits (k-axis)

	foundFirstNonZeroDigit := false
	for r, res := 1, H3_GET_RESOLUTION(h); r <= res; r++ {
		// rotate this digit
		H3_SET_INDEX_DIGIT(&h, r, _rotate60ccw(H3_GET_INDEX_DIGIT(h, r)))

		// look for the first non-zero digit so we
		// can adjust for deleted k-axes sequence
		// if necessary
		if !foundFirstNonZeroDigit && H3_GET_INDEX_DIGIT(h, r) != 0 {
			foundFirstNonZeroDigit = true

			// adjust for deleted k-axes sequence
			if _h3LeadingNonZeroDigit(h) == K_AXES_DIGIT {
				h = _h3Rotate60ccw(h)
			}
		}
	}
	return h
}

// _h3RotatePent60cw rotate an H3Index 60 degrees clockwise about a pentagonal
// center.
func _h3RotatePent60cw(h H3Index) H3Index {
	// rotate in place; skips any leading 1 digits (k-axis)

	foundFirstNonZeroDigit := false
	for r, res := 1, H3_GET_RESOLUTION(h); r <= res; r++ {
		// rotate this digit
		H3_SET_INDEX_DIGIT(&h, r, _rotate60cw(H3_GET_INDEX_DIGIT(h, r)))

		// look for the first non-zero digit so we
		// can adjust for deleted k-axes sequence
		// if necessary
		if !foundFirstNonZeroDigit && H3_GET_INDEX_DIGIT(h, r) != 0 {
			foundFirstNonZeroDigit = true

			// adjust for deleted k-axes sequence
			if _h3LeadingNonZeroDigit(h) == K_AXES_DIGIT {
				h = _h3Rotate60cw(h)
			}
		}
	}
	return h
}

// _h3Rotate60ccw rotate an H3Index 60 degrees counter-clockwise.
func _h3Rotate60ccw(h H3Index) H3Index {
	for r, res := 1, H3_GET_RESOLUTION(h); r <= res; r++ {
		oldDigit := H3_GET_INDEX_DIGIT(h, r)
		H3_SET_INDEX_DIGIT(&h, r, _rotate60ccw(oldDigit))
	}

	return h
}

// _h3Rotate60cw rotate an H3Index 60 degrees clockwise.
func _h3Rotate60cw(h H3Index) H3Index {
	for r, res := 1, H3_GET_RESOLUTION(h); r <= res; r++ {
		H3_SET_INDEX_DIGIT(&h, r, _rotate60cw(H3_GET_INDEX_DIGIT(h, r)))
	}

	return h
}

// _faceIjkToH3 convert an FaceIJK address to the corresponding H3Index.
//
// Return The encoded H3Index (or H3_NULL on failure).
func _faceIjkToH3(fijk *FaceIJK, res int) H3Index {
	// initialize the index
	h := H3_INIT
	H3_SET_MODE(&h, H3_HEXAGON_MODE)
	H3_SET_RESOLUTION(&h, res)

	// check for res 0/base cell
	if res == 0 {
		if fijk.coord.i > MAX_FACE_COORD ||
			fijk.coord.j > MAX_FACE_COORD ||
			fijk.coord.k > MAX_FACE_COORD {
			// out of range input
			return H3_NULL
		}

		H3_SET_BASE_CELL(&h, _faceIjkToBaseCell(fijk))
		return h
	}

	// we need to find the correct base cell FaceIJK for this H3 index;
	// start with the passed in face and resolution res ijk coordinates
	// in that face's coordinate system
	fijkBC := *fijk

	// build the H3Index from finest res up
	// adjust r for the fact that the res 0 base cell offsets the indexing
	// digits
	ijk := &fijkBC.coord
	for r := res - 1; r >= 0; r-- {
		lastIJK := *ijk
		var lastCenter CoordIJK
		if isResClassIII(r + 1) {
			// rotate ccw
			_upAp7(ijk)
			lastCenter = *ijk
			_downAp7(&lastCenter)
		} else {
			// rotate cw
			_upAp7r(ijk)
			lastCenter = *ijk
			_downAp7r(&lastCenter)
		}

		var diff CoordIJK
		_ijkSub(&lastIJK, &lastCenter, &diff)
		_ijkNormalize(&diff)

		H3_SET_INDEX_DIGIT(&h, r+1, _unitIjkToDigit(&diff))
	}

	// fijkBC should now hold the IJK of the base cell in the
	// coordinate system of the current face

	if fijkBC.coord.i > MAX_FACE_COORD ||
		fijkBC.coord.j > MAX_FACE_COORD ||
		fijkBC.coord.k > MAX_FACE_COORD {
		// out of range input
		return H3_NULL
	}

	// lookup the correct base cell
	baseCell := _faceIjkToBaseCell(&fijkBC)
	H3_SET_BASE_CELL(&h, baseCell)

	// rotate if necessary to get canonical base cell orientation
	// for this base cell
	numRots := _faceIjkToBaseCellCCWrot60(&fijkBC)
	if _isBaseCellPentagon(baseCell) {
		// force rotation out of missing k-axes sub-sequence
		if _h3LeadingNonZeroDigit(h) == K_AXES_DIGIT {
			// check for a cw/ccw offset face; default is ccw
			if _baseCellIsCwOffset(baseCell, fijkBC.face) {
				h = _h3Rotate60cw(h)
			} else {
				h = _h3Rotate60ccw(h)
			}
		}

		for i := 0; i < numRots; i++ {
			h = _h3RotatePent60ccw(h)
		}
	} else {
		for i := 0; i < numRots; i++ {
			h = _h3Rotate60ccw(h)
		}
	}

	return h
}

// GeoToH3 encodes a coordinate on the sphere to the H3 index of the containing cell at
// the specified resolution.
//
// Return The encoded H3Index (or H3_NULL on failure).
func GeoToH3(g *GeoCoord, res int) H3Index {
	if res < 0 || res > MAX_H3_RES {
		return H3_NULL
	}

	if !math.IsInf(g.lat, 0) || !math.IsInf(g.lon, 0) {
		return H3_NULL
	}

	var fijk FaceIJK
	_geoToFaceIjk(g, res, &fijk)
	return _faceIjkToH3(&fijk, res)
}

// _h3ToFaceIjkWithInitializedFijk convert an H3Index to the FaceIJK address on
// a specified icosahedral face.
//
// Return true if the possibility of overage exists, otherwise false.
func _h3ToFaceIjkWithInitializedFijk(h H3Index, fijk *FaceIJK) bool {
	ijk := &fijk.coord
	res := H3_GET_RESOLUTION(h)

	// center base cell hierarchy is entirely on this face
	possibleOverage := true
	if !_isBaseCellPentagon(H3_GET_BASE_CELL(h)) &&
		(res == 0 ||
			(fijk.coord.i == 0 && fijk.coord.j == 0 && fijk.coord.k == 0)) {
		possibleOverage = false
	}

	for r := 1; r <= res; r++ {
		if isResClassIII(r) {
			// Class III == rotate ccw
			_downAp7(ijk)
		} else {
			// Class II == rotate cw
			_downAp7r(ijk)
		}

		_neighbor(ijk, H3_GET_INDEX_DIGIT(h, r))
	}

	return possibleOverage
}

// _h3ToFaceIjk convert an H3Index to a FaceIJK address.
func _h3ToFaceIjk(h H3Index, fijk *FaceIJK) {
	baseCell := H3_GET_BASE_CELL(h)
	// adjust for the pentagonal missing sequence; all of sub-sequence 5 needs
	// to be adjusted (and some of sub-sequence 4 below)
	if _isBaseCellPentagon(baseCell) && _h3LeadingNonZeroDigit(h) == 5 {
		h = _h3Rotate60cw(h)
	}

	// start with the "home" face and ijk+ coordinates for the base cell of c
	*fijk = baseCellData[baseCell].homeFijk
	if !_h3ToFaceIjkWithInitializedFijk(h, fijk) {
		return // no overage is possible; h lies on this face
	}

	// if we're here we have the potential for an "overage"; i.e., it is
	// possible that c lies on an adjacent face

	origIJK := fijk.coord

	// if we're in Class III, drop into the next finer Class II grid
	res := H3_GET_RESOLUTION(h)
	if isResClassIII(res) {
		// Class III
		_downAp7r(&fijk.coord)
		res++
	}

	// adjust for overage if needed
	// a pentagon base cell with a leading 4 digit requires special handling
	pentLeading4 := (_isBaseCellPentagon(baseCell) && _h3LeadingNonZeroDigit(h) == 4)
	if _adjustOverageClassII(fijk, res, pentLeading4, false) != NO_OVERAGE {
		// if the base cell is a pentagon we have the potential for secondary
		// overages
		if _isBaseCellPentagon(baseCell) {
			for _adjustOverageClassII(fijk, res, false, false) != NO_OVERAGE {
				continue
			}
		}

		if res != H3_GET_RESOLUTION(h) {
			_upAp7r(&fijk.coord)
		}
	} else if res != H3_GET_RESOLUTION(h) {
		fijk.coord = origIJK
	}
}

// H3ToGeo determines the spherical coordinates of the center point of an
// H3Index.
func H3ToGeo(h3 H3Index, g *GeoCoord) {
	var fijk FaceIJK
	_h3ToFaceIjk(h3, &fijk)
	_faceIjkToGeo(&fijk, H3_GET_RESOLUTION(h3), g)
}

// H3ToGeoBoundary determines the cell boundary in spherical coordinates for an H3 index.
func H3ToGeoBoundary(h3 H3Index, gb *GeoBoundary) {
	var fijk FaceIJK
	_h3ToFaceIjk(h3, &fijk)
	if H3IsPentagon(h3) {
		_faceIjkPentToGeoBoundary(&fijk, H3_GET_RESOLUTION(h3), 0,
			NUM_PENT_VERTS, gb)
	} else {
		_faceIjkToGeoBoundary(&fijk, H3_GET_RESOLUTION(h3), 0, NUM_HEX_VERTS,
			gb)
	}
}

// MaxFaceCount returns the max number of possible icosahedron faces an H3 index
// may intersect.
func MaxFaceCount(h3 H3Index) int {
	// a pentagon always intersects 5 faces, a hexagon never intersects more
	// than 2 (but may only intersect 1)
	if H3IsPentagon(h3) {
		return 5
	}
	return 2
}

// H3GetFaces find all icosahedron faces intersected by a given H3 index,
// represented as integers from 0-19. The array is sparse; since 0 is a valid
// value, invalid array values are represented as -1. It is the responsibility
// of the caller to filter out invalid values.
//
// @param out Output array. Must be of size maxFaceCount(h3).
func H3GetFaces(h3 H3Index, out *[]int) {
	res := H3_GET_RESOLUTION(h3)
	isPentagon := H3IsPentagon(h3)

	// We can't use the vertex-based approach here for class II pentagons,
	// because all their vertices are on the icosahedron edges. Their
	// direct child pentagons cross the same faces, so use those instead.
	if isPentagon && !isResClassIII(res) {
		// Note that this would not work for res 15, but this is only run on
		// Class II pentagons, it should never be invoked for a res 15 index.
		childPentagon := makeDirectChild(h3, 0)
		H3GetFaces(childPentagon, out)
		return
	}

	// convert to FaceIJK
	var fijk FaceIJK
	_h3ToFaceIjk(h3, &fijk)

	// Get all vertices as FaceIJK addresses. For simplicity, always
	// initialize the array with 6 verts, ignoring the last one for pentagons
	var fijkVerts []FaceIJK
	var vertexCount int

	if isPentagon {
		vertexCount = NUM_PENT_VERTS
		fijkVerts = faceIjkPentToVerts(&fijk, &res)
	} else {
		vertexCount = NUM_HEX_VERTS
		fijkVerts = faceIjkToVerts(&fijk, &res)
	}

	// We may not use all of the slots in the output array,
	// so fill with invalid values to indicate unused slots
	faceCount := MaxFaceCount(h3)
	for i := 0; i < faceCount; i++ {
		(*out)[i] = INVALID_FACE
	}

	// add each vertex face, using the output array as a hash set
	for i := 0; i < vertexCount; i++ {
		vert := &fijkVerts[i]

		// Adjust overage, determining whether this vertex is
		// on another face
		if isPentagon {
			_adjustPentVertOverage(vert, res)
		} else {
			_adjustOverageClassII(vert, res, false, true)
		}

		// Save the face to the output array
		face := vert.face
		pos := 0
		// Find the first empty output position, or the first position
		// matching the current face
		for (*out)[pos] != INVALID_FACE && (*out)[pos] != face {
			pos++
		}
		(*out)[pos] = face
	}
}

// PentagonIndexCount returns the number of pentagons (same at any resolution)
func PentagonIndexCount() int {
	return NUM_PENTAGONS
}

// GetPentagonIndexes generates all pentagons at the specified resolution.
func GetPentagonIndexes(res int, out *[]H3Index) {
	i := 0
	for bc := 0; bc < NUM_BASE_CELLS; bc++ {
		if _isBaseCellPentagon(bc) {
			var pentagon H3Index
			setH3Index(&pentagon, res, bc, 0)
			(*out)[i] = pentagon
			i++
		}
	}
}

// isResClassIII returns whether or not a resolution is a Class III grid. Note
// that odd resolutions are Class III and even resolutions are Class II.
//
// Return true if the resolution is a Class III grid, and false if the
// resolution is a Class II grid.
func isResClassIII(res int) bool {
	return res%2 == 1
}
