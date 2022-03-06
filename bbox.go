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

// BBox is geographic bounding box with coordinates defined in radians
type BBox struct {
	north float64 // north latitude
	south float64 // south latitude
	east  float64 // east longitude
	west  float64 // west longitude
}

/**
 * Whether the given bounding box crosses the antimeridian
 * @param  bbox Bounding box to inspect
 * @return      is transmeridian
 */
func bboxIsTransmeridian(bbox *BBox) bool {
	return bbox.east < bbox.west
}

/**
 * Get the center of a bounding box
 * @param bbox   Input bounding box
 * @param center Output center coordinate
 */
func bboxCenter(bbox *BBox, center *GeoCoord) {
	center.lat = (bbox.north + bbox.south) / 2.0
	// If the bbox crosses the antimeridian, shift east 360 degrees
	east := bbox.east
	if bboxIsTransmeridian(bbox) {
		east = bbox.east + M_2PI
	}
	center.lon = constrainLng((east + bbox.west) / 2.0)
}

/**
 * Whether the bounding box contains a given point
 * @param  bbox  Bounding box
 * @param  point Point to test
 * @return       Whether the point is contained
 */
func bboxContains(bbox *BBox, point *GeoCoord) bool {
	if bboxIsTransmeridian(bbox) {
		return point.lat >= bbox.south && point.lat <= bbox.north &&
			(point.lon >= bbox.west || point.lon <= bbox.east)
	}
	return point.lat >= bbox.south && point.lat <= bbox.north &&
		(point.lon >= bbox.west && point.lon <= bbox.east)
}

/**
 * Whether two bounding boxes are strictly equal
 * @param  b1 Bounding box 1
 * @param  b2 Bounding box 2
 * @return    Whether the boxes are equal
 */
func bboxEquals(b1, b2 *BBox) bool {
	return b1.north == b2.north && b1.south == b2.south &&
		b1.east == b2.east && b1.west == b2.west
}

/**
 * _hexRadiusKm returns the radius of a given hexagon in Km
 *
 * @param h3Index the index of the hexagon
 * @return the radius of the hexagon in Km
 */
func _hexRadiusKm(h3Index H3Index) float64 {
	// There is probably a cheaper way to determine the radius of a
	// hexagon, but this way is conceptually simple
	var h3Center GeoCoord
	var h3Boundary GeoBoundary
	H3ToGeo(h3Index, &h3Center)
	H3ToGeoBoundary(h3Index, &h3Boundary)
	return PointDistKm(&h3Center, &h3Boundary.verts[0])
}

/**
 * bboxHexEstimate returns an estimated number of hexagons that fit
 *                 within the cartesian-projected bounding box
 *
 * @param bbox the bounding box to estimate the hexagon fill level
 * @param res the resolution of the H3 hexagons to fill the bounding box
 * @return the estimated number of hexagons to fill the bounding box
 */
func bboxHexEstimate(bbox *BBox, res int) int {
	// Get the area of the pentagon as the maximally-distorted area possible
	pentagons := make([]H3Index, 12)
	GetPentagonIndexes(res, &pentagons)
	pentagonRadiusKm := _hexRadiusKm(pentagons[0])
	// Area of a regular hexagon is 3/2*sqrt(3) * r * r
	// The pentagon has the most distortion (smallest edges) and shares its
	// edges with hexagons, so the most-distorted hexagons have this area,
	// shrunk by 20% off chance that the bounding box perfectly bounds a
	// pentagon.
	pentagonAreaKm2 := 0.8 * (2.59807621135 * pentagonRadiusKm * pentagonRadiusKm)

	// Then get the area of the bounding box of the geofence in question
	var p1, p2 GeoCoord
	p1.lat = bbox.north
	p1.lon = bbox.east
	p2.lat = bbox.south
	p2.lon = bbox.west
	d := PointDistKm(&p1, &p2)
	// Derived constant based on: https://math.stackexchange.com/a/1921940
	// Clamped to 3 as higher values tend to rapidly drag the estimate to zero.
	a := d * d / math.Min(3.0, math.Abs((p1.lon-p2.lon)/(p1.lat-p2.lat)))

	// Divide the two to get an estimate of the number of hexagons needed
	estimate := int(math.Ceil(a / pentagonAreaKm2))
	if estimate == 0 {
		estimate = 1
	}
	return estimate
}

/**
 * lineHexEstimate returns an estimated number of hexagons that trace
 *                 the cartesian-projected line
 *
 *  @param origin the origin coordinates
 *  @param destination the destination coordinates
 *  @param res the resolution of the H3 hexagons to trace the line
 *  @return the estimated number of hexagons required to trace the line
 */
func lineHexEstimate(origin *GeoCoord, destination *GeoCoord, res int) int {
	// Get the area of the pentagon as the maximally-distorted area possible
	pentagons := make([]H3Index, 12)
	GetPentagonIndexes(res, &pentagons)
	pentagonRadiusKm := _hexRadiusKm(pentagons[0])

	dist := PointDistKm(origin, destination)
	estimate := int(math.Ceil(dist / (2 * pentagonRadiusKm)))
	if estimate == 0 {
		estimate = 1
	}
	return estimate
}
