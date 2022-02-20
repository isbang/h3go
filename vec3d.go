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

// Vec3d is 3D floating point structure
type Vec3d struct {
	x float64 // x component
	y float64 // y component
	z float64 // z component
}

// _square square of a number
func _square(x float64) float64 { return x * x }

// _pointSquareDist calculates the square of the distance between two 3D
// coordinates.
func _pointSquareDist(v1, v2 *Vec3d) float64 {
	return _square(v1.x-v2.x) + _square(v1.y-v2.y) + _square(v1.z-v2.z)
}

// _geoToVec3d calculate the 3D coordinate on unit sphere from the latitude and
// longitude.
func _geoToVec3d(geo *GeoCoord) *Vec3d {
	r := math.Cos(geo.lat)

	return &Vec3d{
		x: math.Sin(geo.lat),
		y: math.Cos(geo.lon) * r,
		z: math.Sin(geo.lon) * r,
	}
}
