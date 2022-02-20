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

import "math"

// CoordIJK is IJK hexagon coordinates
//
// Each axis is spaced 120 degrees apart.
type CoordIJK struct {
	i int // i component
	j int // j component
	k int // k component
}

// UNIT_VECS is CoordIJK unit vectors corresponding to the 7 H3 digits.
var UNIT_VECS = [...]CoordIJK{
	{0, 0, 0}, // direction 0
	{0, 0, 1}, // direction 1
	{0, 1, 0}, // direction 2
	{0, 1, 1}, // direction 3
	{1, 0, 0}, // direction 4
	{1, 0, 1}, // direction 5
	{1, 1, 0}, // direction 6
}

// SetIJK sets an IJK coordinate to the specified component values.
func (ijk *CoordIJK) SetIJK(i, j, k int) {
	ijk.i = i
	ijk.j = j
	ijk.k = k
}

// ToHex2d finds the center point in 2D cartesian coordinates of a hex.
func (ijk *CoordIJK) ToHex2d() *Vec2d {
	i := ijk.i - ijk.k
	j := ijk.j - ijk.k

	return &Vec2d{
		x: float64(i) - 0.5*float64(j),
		y: float64(j) * M_SQRT3_2,
	}
}

// Scale uniformly scales ijk coordinates by a scalar. Works in place.
func (ijk *CoordIJK) Scale(factor int) {
	ijk.i *= factor
	ijk.j *= factor
	ijk.k *= factor
}

// Normalize normalizes ijk coordinates by setting the components to the
// smallest possible values. Works in place.
func (ijk *CoordIJK) Normalize() {
	// remove any negative values
	if ijk.i < 0 {
		ijk.j -= ijk.i
		ijk.k -= ijk.i
		ijk.i = 0
	}

	if ijk.j < 0 {
		ijk.i -= ijk.j
		ijk.k -= ijk.j
		ijk.j = 0
	}

	if ijk.k < 0 {
		ijk.i -= ijk.k
		ijk.j -= ijk.k
		ijk.k = 0
	}

	// remove the min value if needed
	min := ijk.i

	if ijk.j < min {
		min = ijk.j
	}

	if ijk.k < min {
		min = ijk.k
	}

	if min > 0 {
		ijk.i -= min
		ijk.j -= min
		ijk.k -= min
	}
}

// UnitToDigit determines the H3 digit corresponding to a unit vector in ijk
// coordinates.
//
// Return the H3 digit (0-6) corresponding to the ijk unit vector, or
// INVALID_DIGIT on failure.
func (ijk *CoordIJK) UnitToDigit() Direction {
	c := *ijk
	_ijkNormalize(&c)

	digit := INVALID_DIGIT
	for i := CENTER_DIGIT; i < Direction(NUM_DIGITS); i++ {
		if _ijkMatches(&c, &UNIT_VECS[i]) {
			digit = i
			break
		}
	}

	return digit
}

// upAp7 finds the normalized ijk coordinates of the indexing parent of a cell
// in a counter-clockwise aperture 7 grid. Works in place.
func (ijk *CoordIJK) upAp7() {
	// convert to CoordIJ
	i := ijk.i - ijk.k
	j := ijk.j - ijk.k

	ijk.i = int(math.Round(float64((3*i - j) / 7.0)))
	ijk.j = int(math.Round(float64((i + 2*j) / 7.0)))
	ijk.k = 0
	_ijkNormalize(ijk)
}

// upAp7r finds the normalized ijk coordinates of the indexing parent of a cell in a clockwise aperture 7 grid. Works in place.
func (ijk *CoordIJK) upAp7r() {
	// convert to CoordIJ
	i := ijk.i - ijk.k
	j := ijk.j - ijk.k

	ijk.i = int(math.Round(float64((2*i + j) / 7.0)))
	ijk.j = int(math.Round(float64((3*j - i) / 7.0)))
	ijk.k = 0
	_ijkNormalize(ijk)
}

// downAp7 finds the normalized ijk coordinates of the hex centered on the
// indicated hex at the next finer aperture 7 counter-clockwise resolution.
// Works in place.
func (ijk *CoordIJK) downAp7() {
	// res r unit vectors in res r+1
	iVec := CoordIJK{3, 0, 1}
	jVec := CoordIJK{1, 3, 0}
	kVec := CoordIJK{0, 1, 3}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// downAp7r finds the normalized ijk coordinates of the hex centered on the
// indicated hex at the next finer aperture 7 clockwise resolution.
// Works in place.
func (ijk *CoordIJK) downAp7r() {
	// res r unit vectors in res r+1
	iVec := CoordIJK{3, 1, 0}
	jVec := CoordIJK{0, 3, 1}
	kVec := CoordIJK{1, 0, 3}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// neighbor finds the normalized ijk coordinates of the hex in the specified
// digit direction from the specified ijk coordinates. Works in place.
func (ijk *CoordIJK) neighbor(digit Direction) {
	if digit > CENTER_DIGIT && digit < Direction(NUM_DIGITS) {
		_ijkAdd(ijk, &UNIT_VECS[digit], ijk)
		_ijkNormalize(ijk)
	}
}

// Rotate60ccw rotates ijk coordinates 60 degrees counter-clockwise.
// Works in place.
func (ijk *CoordIJK) Rotate60ccw() {
	// unit vector rotations
	iVec := CoordIJK{1, 1, 0}
	jVec := CoordIJK{0, 1, 1}
	kVec := CoordIJK{1, 0, 1}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// Rotate60cw rotates ijk coordinates 60 degrees clockwise. Works in place.
func (ijk *CoordIJK) Rotate60cw() {
	// unit vector rotations
	iVec := CoordIJK{1, 0, 1}
	jVec := CoordIJK{1, 1, 0}
	kVec := CoordIJK{0, 1, 1}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// downAp3 finds the normalized ijk coordinates of the hex centered on the
// indicated hex at the next finer aperture 3 counter-clockwise resolution.
// Works in place.
func (ijk *CoordIJK) downAp3() {
	// res r unit vectors in res r+1
	iVec := CoordIJK{2, 0, 1}
	jVec := CoordIJK{1, 2, 0}
	kVec := CoordIJK{0, 1, 2}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// downAp3r finds the normalized ijk coordinates of the hex centered on the
// indicated hex at the next finer aperture 3 clockwise resolution.
// Works in place.
func (ijk *CoordIJK) downAp3r() {
	// res r unit vectors in res r+1
	iVec := CoordIJK{2, 1, 0}
	jVec := CoordIJK{0, 2, 1}
	kVec := CoordIJK{1, 0, 2}

	_ijkScale(&iVec, ijk.i)
	_ijkScale(&jVec, ijk.j)
	_ijkScale(&kVec, ijk.k)

	_ijkAdd(&iVec, &jVec, ijk)
	_ijkAdd(ijk, &kVec, ijk)

	_ijkNormalize(ijk)
}

// ToCube convert IJK coordinates to cube coordinates, in place
func (ijk *CoordIJK) ToCube() {
	ijk.i = -ijk.i + ijk.k
	ijk.j = ijk.j - ijk.k
	ijk.k = -ijk.i - ijk.j
}

// _setIJK sets an IJK coordinate to the specified component values.
//
// Deprecated: Use (*CoordIJK).SetIJK instead.
func _setIJK(ijk *CoordIJK, i, j, k int) {
	ijk.SetIJK(i, j, k)
}

// _hex2dToCoordIJK determine the containing hex in ijk+ coordinates for a 2D
// cartesian coordinate vector (from DGGRID).
func _hex2dToCoordIJK(v *Vec2d, h *CoordIJK) {
	var a1, a2 float64
	var x1, x2 float64
	var m1, m2 int
	var r1, r2 float64

	// quantize into the ij system and then normalize
	h.k = 0

	a1 = math.Abs(v.x)
	a2 = math.Abs(v.y)

	// first do a reverse conversion
	x2 = a2 / M_SIN60
	x1 = a1 + x2/2.0

	// check if we have the center of a hex
	m1 = int(x1)
	m2 = int(x2)

	// otherwise round correctly
	r1 = x1 - float64(m1)
	r2 = x2 - float64(m2)

	if r1 < 0.5 {
		if r1 < 1.0/3.0 {
			if r2 < (1.0+r1)/2.0 {
				h.i = m1
				h.j = m2
			} else {
				h.i = m1
				h.j = m2 + 1
			}
		} else {
			if r2 < (1.0 - r1) {
				h.j = m2
			} else {
				h.j = m2 + 1
			}

			if (1.0-r1) <= r2 && r2 < (2.0*r1) {
				h.i = m1 + 1
			} else {
				h.i = m1
			}
		}
	} else {
		if r1 < 2.0/3.0 {
			if r2 < (1.0 - r1) {
				h.j = m2
			} else {
				h.j = m2 + 1
			}

			if (2.0*r1-1.0) < r2 && r2 < (1.0-r1) {
				h.i = m1
			} else {
				h.i = m1 + 1
			}
		} else {
			if r2 < (r1 / 2.0) {
				h.i = m1 + 1
				h.j = m2
			} else {
				h.i = m1 + 1
				h.j = m2 + 1
			}
		}
	}

	// now fold across the axes if necessary

	if v.x < 0.0 {
		if (h.j % 2) == 0 { // even
			axisi := int64(h.j) / int64(2)
			diff := int64(h.i) - axisi
			h.i = int(int64(h.i) - 2*diff)
		} else {
			axisi := int64(h.j+1) / 2
			diff := int64(h.i) - axisi
			h.i = int(int64(h.i) - (2*diff + 1))
		}
	}

	if v.y < 0.0 {
		h.i = h.i - (2*h.j+1)/2
		h.j = -1 * h.j
	}

	_ijkNormalize(h)
}

// _ijkToHex2d finds the center point in 2D cartesian coordinates of a hex.
func _ijkToHex2d(h *CoordIJK, v *Vec2d) {
	i := h.i - h.k
	j := h.j - h.k

	v.x = float64(i) - 0.5*float64(j)
	v.y = float64(j) * M_SQRT3_2
}

// _ijkMatches returns whether or not two ijk coordinates contain exactly the
// same component values.
func _ijkMatches(c1, c2 *CoordIJK) bool {
	return c1.i == c2.i && c1.j == c2.j && c1.k == c2.k
}

// _ijkAdd adds two ijk coordinates.
func _ijkAdd(h1, h2 *CoordIJK, sum *CoordIJK) {
	sum.i = h1.i + h2.i
	sum.j = h1.j + h2.j
	sum.k = h1.k + h2.k
}

// _ijkSub subtracts two ijk coordinates.
func _ijkSub(h1, h2 *CoordIJK, diff *CoordIJK) {
	diff.i = h1.i - h2.i
	diff.j = h1.j - h2.j
	diff.k = h1.k - h2.k
}

// _ijkScale uniformly scales ijk coordinates by a scalar. Works in place.
//
// Deprecated: Use (*CoordIJK).Scale instead.
func _ijkScale(c *CoordIJK, factor int) {
	c.Scale(factor)
}

// _ijkNormalize normalizes ijk coordinates by setting the components to the
// smallest possible values. Works in place.
//
// Deprecated: Use (*CoordIJK).Normalize instead.
func _ijkNormalize(c *CoordIJK) {
	c.Normalize()
}

// _unitIjkToDigit determines the H3 digit corresponding to a unit vector in ijk
// coordinates.
//
// Return the H3 digit (0-6) corresponding to the ijk unit vector, or
// INVALID_DIGIT on failure.
//
// Deprecated: Use (*CoordIJK).UnitToDigit instead.
func _unitIjkToDigit(ijk *CoordIJK) Direction {
	return ijk.UnitToDigit()
}

// _upAp7 finds the normalized ijk coordinates of the indexing parent of a cell
// in a counter-clockwise aperture 7 grid. Works in place.
//
// Deprecated: Use (*CoordIJK).upAp7 instead.
func _upAp7(ijk *CoordIJK) {
	ijk.upAp7()
}

// _upAp7r finds the normalized ijk coordinates of the indexing parent of a cell
// in a clockwise aperture 7 grid. Works in place.
//
// Deprecated: Use (*CoordIJK).upAp7r instead.
func _upAp7r(ijk *CoordIJK) {
	ijk.upAp7r()
}

// _downAp7 finds the normalized ijk coordinates of the hex centered on the
// indicated hex at the next finer aperture 7 counter-clockwise resolution.
// Works in place.
//
// Deprecated: Use (*CoordIJK).downAp7 instead.
func _downAp7(ijk *CoordIJK) {
	ijk.downAp7()
}

// _downAp7r finds the normalized ijk coordinates of the hex centered on the
// indicated hex at the next finer aperture 7 clockwise resolution.
// Works in place.
//
// Deprecated: Use (*CoordIJK).downAp7r instead.
func _downAp7r(ijk *CoordIJK) {
	ijk.downAp7r()
}

// _neighbor finds the normalized ijk coordinates of the hex in the specified
// digit direction from the specified ijk coordinates. Works in place.
//
// Deprecated: Use (*CoordIJK).neighbor instead.
func _neighbor(ijk *CoordIJK, digit Direction) {
	ijk.neighbor(digit)
}

// _ijkRotate60ccw rotates ijk coordinates 60 degrees counter-clockwise.
// Works in place.
//
// Deprecated: Use (*CoordIJK).Rotate60ccw instead.
func _ijkRotate60ccw(ijk *CoordIJK) {
	ijk.Rotate60ccw()
}

// _ijkRotate60cw rotates ijk coordinates 60 degrees clockwise. Works in place.
//
// Deprecated: Use (*CoordIJK).Rotate60cw instead.
func _ijkRotate60cw(ijk *CoordIJK) {
	ijk.Rotate60cw()
}

// _downAp3 finds the normalized ijk coordinates of the hex centered on the
// indicated hex at the next finer aperture 3 counter-clockwise resolution.
// Works in place.
//
// Deprecated: Use (*CoordIJK).downAp3 instead.
func _downAp3(ijk *CoordIJK) {
	ijk.downAp3()
}

// _downAp3r finds the normalized ijk coordinates of the hex centered on the
// indicated hex at the next finer aperture 3 clockwise resolution.
// Works in place.
//
// Deprecated: Use (*CoordIJK).downAp3r instead.
func _downAp3r(ijk *CoordIJK) {
	ijk.downAp3r()
}

// ijkDistance finds the distance between the two coordinates. Returns result.
func ijkDistance(c1, c2 *CoordIJK) int {
	var diff CoordIJK
	_ijkSub(c1, c2, &diff)
	_ijkNormalize(&diff)
	absDiff := CoordIJK{abs(diff.i), abs(diff.j), abs(diff.k)}
	return max(absDiff.i, max(absDiff.j, absDiff.k))
}

// ijkToIj transforms coordinates from the IJK+ coordinate system to the IJ
// coordinate system.
func ijkToIj(ijk *CoordIJK, ij *CoordIJ) {
	ij.i = ijk.i - ijk.k
	ij.j = ijk.j - ijk.k
}

// ijToIjk transforms coordinates from the IJ coordinate system to the IJK+
// coordinate system.
func ijToIjk(ij *CoordIJ, ijk *CoordIJK) {
	ijk.i = ij.i
	ijk.j = ij.j
	ijk.k = 0

	_ijkNormalize(ijk)
}

// ijkToCube convert IJK coordinates to cube coordinates, in place.
//
// Deprecated: Use (*CoordIJK).ToCube instead
func ijkToCube(ijk *CoordIJK) {
	ijk.ToCube()
}

// cubeToIjk convert cube coordinates to IJK coordinates, in place.
func cubeToIjk(ijk *CoordIJK) {
	ijk.i = -ijk.i
	ijk.k = 0
	_ijkNormalize(ijk)
}
