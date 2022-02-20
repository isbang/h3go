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

// Vec2d is 2D floating-point vector
type Vec2d struct {
	x float64 // x component
	y float64 // y component
}

func (v2d *Vec2d) Magnitude() float64 {
	return math.Sqrt(v2d.x*v2d.x + v2d.y*v2d.y)
}

// _v2dMag calculates the magnitude of a 2D cartesian vector.
//
// Deprecated: Use (*Vec2d).Magnitude instead.
func _v2dMag(v *Vec2d) float64 {
	return v.Magnitude()
}

// v2dIntersect finds the intersection between two lines. Assumes that the
// lines intersect and that the intersection is not at an endpoint of either
// line.
func v2dIntersect(p0, p1, p2, p3 *Vec2d) *Vec2d {
	var s1, s2 Vec2d
	s1.x = p1.x - p0.x
	s1.y = p1.y - p0.y
	s2.x = p3.x - p2.x
	s2.y = p3.y - p2.y

	t := (s2.x*(p0.y-p2.y) - s2.y*(p0.x-p2.x)) / (-s2.x*s1.y + s1.x*s2.y)

	return &Vec2d{
		x: p0.x + (t * s1.x),
		y: p0.y + (t * s1.y),
	}
}

// _v2dIntersect finds the intersection between two lines. Assumes that the
// lines intersect and that the intersection is not at an endpoint of either
// line.
//
// Deprecated: Use v2dIntersect instead.
func _v2dIntersect(p0, p1, p2, p3 *Vec2d, inter *Vec2d) {
	var s1, s2 Vec2d
	s1.x = p1.x - p0.x
	s1.y = p1.y - p0.y
	s2.x = p3.x - p2.x
	s2.y = p3.y - p2.y

	t := (s2.x*(p0.y-p2.y) - s2.y*(p0.x-p2.x)) / (-s2.x*s1.y + s1.x*s2.y)

	inter.x = p0.x + (t * s1.x)
	inter.y = p0.y + (t * s1.y)
}

// _v2dEquals checks whether two 2D vectors are equal. Does not consider
// possible false negatives due to floating-point errors.
func _v2dEquals(v1, v2 *Vec2d) bool {
	return v1.x == v2.x && v1.y == v2.y
}
