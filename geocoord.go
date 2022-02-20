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
	// epsilon of ~0.1mm in degrees
	EPSILON_DEG = .000000001

	// epsilon of ~0.1mm in radians
	EPSILON_RAD = EPSILON_DEG * M_PI_180
)

// GeoCoord has latitude/longitude in radians
type GeoCoord struct {
	lat float64 // latitude in radians
	lon float64 // longitude in radians
}

// _posAngleRads normalizes radians to a value between 0.0 and two PI.
//
// Return The normalized radians value.
func _posAngleRads(rads float64) float64 {
	tmp := rads
	if rads < 0.0 {
		tmp = rads + M_2PI
	}
	if rads >= M_2PI {
		tmp -= M_2PI
	}
	return tmp
}

// geoAlmostEqualThreshold determines if the components of two spherical
// coordinates are within some threshold distance of each other.
//
// Return whether or not the two coordinates are within the threshold distance
// of each other.
func geoAlmostEqualThreshold(p1, p2 *GeoCoord, threshold float64) bool {
	return math.Abs(p1.lat-p2.lat) < threshold &&
		math.Abs(p1.lon-p2.lon) < threshold
}

// geoAlmostEqual determines if the components of two spherical coordinates are
// within our standard epsilon distance of each other.
//
// Return whether or not the two coordinates are within the epsilon distance of
// each other.
func geoAlmostEqual(p1, p2 *GeoCoord) bool {
	return geoAlmostEqualThreshold(p1, p2, EPSILON_RAD)
}

// setGeoDegs set the components of spherical coordinates in decimal degrees.
//
// Deprecated: Use (*GeoCoord).setGeoDegs instead.
func setGeoDegs(p *GeoCoord, latDegs float64, lonDegs float64) {
	_setGeoRads(p, DegsToRads(latDegs), DegsToRads(lonDegs))
}

// setGeoDegs set the components of spherical coordinates in decimal degrees.
func (p *GeoCoord) setGeoDegs(latDegs float64, lonDegs float64) {
	_setGeoRads(p, DegsToRads(latDegs), DegsToRads(lonDegs))
}

// _setGeoRads set the components of spherical coordinates in radians.
//
// Deprecated: Use (*GeoCoord).setGeoRads instead.
func _setGeoRads(p *GeoCoord, latRads float64, lonRads float64) {
	p.lat = latRads
	p.lon = lonRads
}

// setGeoRads set the components of spherical coordinates in radians.
func (p *GeoCoord) setGeoRads(latRads float64, lonRads float64) {
	p.lat = latRads
	p.lon = lonRads
}

// DegsToRads convert from decimal degrees to radians.
//
// Return the corresponding radians.
func DegsToRads(degrees float64) float64 {
	return degrees * M_PI_180
}

// RadsToDegs convert from radians to decimal degrees.
//
// Return the corresponding decimal degrees.
func RadsToDegs(radians float64) float64 {
	return radians * M_180_PI
}

// constrainLat makes sure latitudes are in the proper bounds
//
// Return the corrected lat value
func constrainLat(lat float64) float64 {
	for lat > M_PI_2 {
		lat = lat - M_PI
	}
	return lat
}

// constrainLng makes sure longitudes are in the proper bounds
//
// Return the corrected lng value
func constrainLng(lng float64) float64 {
	for lng > M_PI {
		lng = lng - (2 * M_PI)
	}
	for lng < -M_PI {
		lng = lng + (2 * M_PI)
	}
	return lng
}

// PointDistRads calculates the great circle distance in radians between two
// spherical coordinates.
//
// This function uses the Haversine formula.
// For math details, see:
//     https://en.wikipedia.org/wiki/Haversine_formula
//     https://www.movable-type.co.uk/scripts/latlong.html
//
// Return the great circle distance in radians between a and b
//
func PointDistRads(a, b *GeoCoord) float64 {
	sinLat := math.Sin((b.lat - a.lat) / 2.0)
	sinLng := math.Sin((b.lon - a.lon) / 2.0)

	A := sinLat*sinLat + math.Cos(a.lat)*math.Cos(b.lat)*sinLng*sinLng

	return 2 * math.Atan2(math.Sqrt(A), math.Sqrt(1-A))
}

// PointDistKm calculates the great circle distance in kilometers between two
// spherical coordinates.
func PointDistKm(a, b *GeoCoord) float64 {
	return PointDistRads(a, b) * EARTH_RADIUS_KM
}

// PointDistM calculates the great circle distance in meters between two
// spherical coordinates.
func PointDistM(a, b *GeoCoord) float64 {
	return PointDistKm(a, b) * 1000
}

// _geoAzimuthRads determines the azimuth to p2 from p1 in radians.
//
// Return the azimuth in radians from p1 to p2.
func _geoAzimuthRads(p1, p2 *GeoCoord) float64 {
	return math.Atan2(
		math.Cos(p2.lat)*math.Sin(p2.lon-p1.lon),
		math.Cos(p1.lat)*math.Sin(p2.lat)-
			math.Sin(p1.lat)*math.Cos(p2.lat)*math.Cos(p2.lon-p1.lon),
	)
}

// _geoAzDistanceRads computes the point on the sphere a specified azimuth and
// distance from another point.
//
// Deprecated: Use (*GeoCoord).geoAzDistanceRads instead.
func _geoAzDistanceRads(p1 *GeoCoord, az float64, distance float64, p2 *GeoCoord) {
	if distance < EPSILON {
		*p2 = *p1
		return
	}

	var sinlat, sinlon, coslon float64

	az = _posAngleRads(az)

	// check for due north/south azimuth
	if az < EPSILON || math.Abs(az-M_PI) < EPSILON {
		if az < EPSILON { // due north
			p2.lat = p1.lat + distance
		} else { // due south
			p2.lat = p1.lat - distance
		}

		if math.Abs(p2.lat-M_PI_2) < EPSILON { // north pole
			p2.lat = M_PI_2
			p2.lon = 0.0
		} else if math.Abs(p2.lat+M_PI_2) < EPSILON { // south pole
			p2.lat = -M_PI_2
			p2.lon = 0.0
		} else {
			p2.lon = constrainLng(p1.lon)
		}
	} else { // not due north or south
		sinlat = math.Sin(p1.lat)*math.Cos(distance) +
			math.Cos(p1.lat)*math.Sin(distance)*math.Cos(az)
		if sinlat > 1.0 {
			sinlat = 1.0
		}
		if sinlat < -1.0 {
			sinlat = -1.0
		}
		p2.lat = math.Asin(sinlat)
		if math.Abs(p2.lat-M_PI_2) < EPSILON { // north pole
			p2.lat = M_PI_2
			p2.lon = 0.0
		} else if math.Abs(p2.lat+M_PI_2) < EPSILON { // south pole
			p2.lat = -M_PI_2
			p2.lon = 0.0
		} else {
			sinlon = math.Sin(az) * math.Sin(distance) / math.Cos(p2.lat)
			coslon = (math.Cos(distance) - math.Sin(p1.lat)*math.Sin(p2.lat)) /
				math.Cos(p1.lat) / math.Cos(p2.lat)
			if sinlon > 1.0 {
				sinlon = 1.0
			}
			if sinlon < -1.0 {
				sinlon = -1.0
			}
			if coslon > 1.0 {
				coslon = 1.0
			}
			if coslon < -1.0 {
				coslon = -1.0
			}
			p2.lon = constrainLng(p1.lon + math.Atan2(sinlon, coslon))
		}
	}
}

// geoAzDistanceRads computes the point on the sphere a specified azimuth and
// distance from another point.
func (p *GeoCoord) geoAzDistanceRads(az float64, distance float64) GeoCoord {
	if distance < EPSILON {
		return *p
	}

	var p2 GeoCoord

	var sinlat, sinlon, coslon float64

	az = _posAngleRads(az)

	// check for due north/south azimuth
	if az < EPSILON || math.Abs(az-M_PI) < EPSILON {
		if az < EPSILON { // due north
			p2.lat = p.lat + distance
		} else { // due south
			p2.lat = p.lat - distance
		}

		if math.Abs(p2.lat-M_PI_2) < EPSILON { // north pole
			p2.lat = M_PI_2
			p2.lon = 0.0
		} else if math.Abs(p2.lat+M_PI_2) < EPSILON { // south pole
			p2.lat = -M_PI_2
			p2.lon = 0.0
		} else {
			p2.lon = constrainLng(p.lon)
		}
	} else { // not due north or south
		sinlat = math.Sin(p.lat)*math.Cos(distance) +
			math.Cos(p.lat)*math.Sin(distance)*math.Cos(az)
		if sinlat > 1.0 {
			sinlat = 1.0
		}
		if sinlat < -1.0 {
			sinlat = -1.0
		}
		p2.lat = math.Asin(sinlat)
		if math.Abs(p2.lat-M_PI_2) < EPSILON { // north pole
			p2.lat = M_PI_2
			p2.lon = 0.0
		} else if math.Abs(p2.lat+M_PI_2) < EPSILON { // south pole
			p2.lat = -M_PI_2
			p2.lon = 0.0
		} else {
			sinlon = math.Sin(az) * math.Sin(distance) / math.Cos(p2.lat)
			coslon = (math.Cos(distance) - math.Sin(p.lat)*math.Sin(p2.lat)) /
				math.Cos(p.lat) / math.Cos(p2.lat)
			if sinlon > 1.0 {
				sinlon = 1.0
			}
			if sinlon < -1.0 {
				sinlon = -1.0
			}
			if coslon > 1.0 {
				coslon = 1.0
			}
			if coslon < -1.0 {
				coslon = -1.0
			}
			p2.lon = constrainLng(p.lon + math.Atan2(sinlon, coslon))
		}
	}

	return p2
}

/*
 * The following functions provide meta information about the H3 hexagons at
 * each zoom level. Since there are only 16 total levels, these are current
 * handled with hardwired static values, but it may be worthwhile to put these
 * static values into another file that can be autogenerated by source code in
 * the future.
 */

func HexAreaKm2(res int) float64 {
	var areas = [...]float64{
		4250546.848, 607220.9782, 86745.85403, 12392.26486,
		1770.323552, 252.9033645, 36.1290521, 5.1612932,
		0.7373276, 0.1053325, 0.0150475, 0.0021496,
		0.0003071, 0.0000439, 0.0000063, 0.0000009,
	}
	return areas[res]
}

func HexAreaM2(res int) float64 {
	var areas = [...]float64{
		4.25055e+12, 6.07221e+11, 86745854035, 12392264862,
		1770323552, 252903364.5, 36129052.1, 5161293.2,
		737327.6, 105332.5, 15047.5, 2149.6,
		307.1, 43.9, 6.3, 0.9,
	}
	return areas[res]
}

func EdgeLengthKm(res int) float64 {
	var lens = [...]float64{
		1107.712591, 418.6760055, 158.2446558, 59.81085794,
		22.6063794, 8.544408276, 3.229482772, 1.220629759,
		0.461354684, 0.174375668, 0.065907807, 0.024910561,
		0.009415526, 0.003559893, 0.001348575, 0.000509713,
	}
	return lens[res]
}

func EdgeLengthM(res int) float64 {
	var lens = [...]float64{
		1107712.591, 418676.0055, 158244.6558, 59810.85794,
		22606.3794, 8544.408276, 3229.482772, 1220.629759,
		461.3546837, 174.3756681, 65.90780749, 24.9105614,
		9.415526211, 3.559893033, 1.348574562, 0.509713273,
	}
	return lens[res]
}

// NumHexagons returns number of unique valid H3Indexes at given resolution.
func NumHexagons(res int) int64 {
	/**
	 * Note: this *actually* returns the number of *cells*
	 * (which includes the 12 pentagons) at each resolution.
	 *
	 * This table comes from the recurrence:
	 *
	 *  num_cells(0) = 122
	 *  num_cells(i+1) = (num_cells(i)-12)*7 + 12*6
	 *
	 */
	var nums = [...]int64{
		122,
		842,
		5882,
		41162,
		288122,
		2016842,
		14117882,
		98825162,
		691776122,
		4842432842,
		33897029882,
		237279209162,
		1660954464122,
		11626681248842,
		81386768741882,
		569707381193162,
	}
	return nums[res]
}

// triangleEdgeLengthsToArea calculates surface area in radians^2 of spherical
// triangle on unit sphere.
//
// For the math, see:
//   https://en.wikipedia.org/wiki/Spherical_trigonometry#Area_and_spherical_excess
//
// Return area in radians^2 of triangle on unit sphere
func triangleEdgeLengthsToArea(a, b, c float64) float64 {
	s := (a + b + c) / 2

	a = (s - a) / 2
	b = (s - b) / 2
	c = (s - c) / 2
	s = s / 2

	return 4 * math.Atan(math.Sqrt(math.Tan(s)*math.Tan(a)*math.Tan(b)*math.Tan(c)))
}

// triangleArea computes area in radians^2 of a spherical triangle, given its
// vertices.
//
// Return area of triangle on unit sphere, in radians^2
func triangleArea(a, b, c *GeoCoord) float64 {
	return triangleEdgeLengthsToArea(
		PointDistRads(a, b),
		PointDistRads(b, c),
		PointDistRads(c, a),
	)
}

// CellAreaRads2 computes area of H3 cell in radians^2.
//
// The area is calculated by breaking the cell into spherical triangles and
// summing up their areas. Note that some H3 cells (hexagons and pentagons)
// are irregular, and have more than 6 or 5 sides.
//
// TODO: optimize the computation by re-using the edges shared between triangles
//
// Return cell area in radians^2
func CellAreaRads2(cell H3Index) float64 {
	var c GeoCoord
	var gb GeoBoundary
	H3ToGeo(cell, &c)
	H3ToGeoBoundary(cell, &gb)

	area := 0.0
	for i := 0; i < gb.numVerts; i++ {
		j := (i + 1) % gb.numVerts
		area += triangleArea(&gb.verts[i], &gb.verts[j], &c)
	}

	return area
}

// CellAreaKm2 computes area of H3 cell in kilometers^2.
func CellAreaKm2(h H3Index) float64 {
	return CellAreaRads2(h) * EARTH_RADIUS_KM * EARTH_RADIUS_KM
}

// CellAreaM2 computes area of H3 cell in meters^2.
func CellAreaM2(h H3Index) float64 {
	return CellAreaKm2(h) * 1000 * 1000
}

// ExactEdgeLengthRads computes length of a unidirectional edge in radians.
//
// Return length in radians
func ExactEdgeLengthRads(edge H3Index) float64 {
	var gb GeoBoundary

	GetH3UnidirectionalEdgeBoundary(edge, &gb)

	length := 0.0
	for i := 0; i < gb.numVerts-1; i++ {
		length += PointDistRads(&gb.verts[i], &gb.verts[i+1])
	}

	return length
}

// ExactEdgeLengthKm computes length of a unidirectional edge in kilometers.
func ExactEdgeLengthKm(edge H3Index) float64 {
	return ExactEdgeLengthRads(edge) * EARTH_RADIUS_KM
}

// ExactEdgeLengthM computes length of a unidirectional edge in meters.
func ExactEdgeLengthM(edge H3Index) float64 {
	return ExactEdgeLengthKm(edge) * 1000
}
