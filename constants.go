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

const (
	// pi
	M_PI = math.Pi // 3.14159265358979323846

	// pi / 2.0
	M_PI_2 = math.Pi / 2.0 // 1.5707963267948966

	// 2.0 * pi
	M_2PI = 2.0 * math.Pi // 6.28318530717958647692528676655900576839433

	// pi / 180
	M_PI_180 = math.Pi / 180 // 0.0174532925199432957692369076848861271111
	// pi * 180
	M_180_PI = math.Pi * 180 // 57.29577951308232087679815481410517033240547

	// threshold epsilon
	EPSILON = 0.0000000000000001
	// sqrt(3) / 2.0
	M_SQRT3_2 = 0.8660254037844386467637231707529361834714
	// sin(60')
	M_SIN60 = M_SQRT3_2

	// rotation angle between Class II and Class III resolution axes
	// (asin(sqrt(3.0 / 28.0)))
	M_AP7_ROT_RADS = 0.333473172251832115336090755351601070065900389

	// sin(M_AP7_ROT_RADS)
	M_SIN_AP7_ROT = 0.3273268353539885718950318

	// cos(M_AP7_ROT_RADS)
	M_COS_AP7_ROT = 0.9449111825230680680167902

	// earth radius in kilometers using WGS84 authalic radius
	EARTH_RADIUS_KM = 6371.007180918475

	// scaling factor from hex2d resolution 0 unit length (or distance between
	// adjacent cell center points on the plane) to gnomonic unit length.
	RES0_U_GNOMONIC = 0.38196601125010500003

	// max H3 resolution; H3 version 1 has 16 resolutions, numbered 0 through 15
	MAX_H3_RES = 15

	// The number of faces on an icosahedron
	NUM_ICOSA_FACES = 20
	// The number of H3 base cells
	NUM_BASE_CELLS = 122
	// The number of vertices in a hexagon
	NUM_HEX_VERTS = 6
	// The number of vertices in a pentagon
	NUM_PENT_VERTS = 5
	// The number of pentagons per resolution
	NUM_PENTAGONS = 12

	// H3 index modes
	H3_HEXAGON_MODE = 1
	H3_UNIEDGE_MODE = 2
)
