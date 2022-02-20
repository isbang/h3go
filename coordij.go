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

// CoordIJ is IJ hexagon coordinates
// Each axis is spaced 120 degrees apart.
type CoordIJ struct {
	i int // i component
	j int // j component
}

// ToIJK transforms coordinates from the IJ coordinate system to the IJK+
// coordinate system.
func (ij *CoordIJ) ToIJK() CoordIJK {
	ijk := CoordIJK{
		i: ij.i,
		j: ij.j,
		k: 0,
	}

	_ijkNormalize(&ijk)
	return ijk
}
