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

// H3IndexesAreNeighbors returns whether or not the provided H3Indexes are
// neighbors.
func H3IndexesAreNeighbors(origin H3Index, destination H3Index) bool {
	// Make sure they're hexagon indexes
	if H3_GET_MODE(origin) != H3_HEXAGON_MODE ||
		H3_GET_MODE(destination) != H3_HEXAGON_MODE {
		return false
	}

	// Hexagons cannot be neighbors with themselves
	if origin == destination {
		return false
	}

	// Only hexagons in the same resolution can be neighbors
	if H3_GET_RESOLUTION(origin) != H3_GET_RESOLUTION(destination) {
		return false
	}

	// H3 Indexes that share the same parent are very likely to be neighbors
	// Child 0 is neighbor with all of its parent's 'offspring', the other
	// children are neighbors with 3 of the 7 children. So a simple comparison
	// of origin and destination parents and then a lookup table of the children
	// is a super-cheap way to possibly determine they are neighbors.
	parentRes := H3_GET_RESOLUTION(origin) - 1
	if parentRes > 0 && (H3ToParent(origin, parentRes) == H3ToParent(destination, parentRes)) {
		originResDigit := H3_GET_INDEX_DIGIT(origin, parentRes+1)
		destinationResDigit := H3_GET_INDEX_DIGIT(destination, parentRes+1)
		if originResDigit == CENTER_DIGIT || destinationResDigit == CENTER_DIGIT {
			return true
		}
		// These sets are the relevant neighbors in the clockwise
		// and counter-clockwise
		var neighborSetClockwise = []Direction{
			CENTER_DIGIT, JK_AXES_DIGIT, IJ_AXES_DIGIT, J_AXES_DIGIT,
			IK_AXES_DIGIT, K_AXES_DIGIT, I_AXES_DIGIT,
		}
		var neighborSetCounterclockwise = []Direction{
			CENTER_DIGIT, IK_AXES_DIGIT, JK_AXES_DIGIT, K_AXES_DIGIT,
			IJ_AXES_DIGIT, I_AXES_DIGIT, J_AXES_DIGIT,
		}
		if neighborSetClockwise[originResDigit] == destinationResDigit ||
			neighborSetCounterclockwise[originResDigit] == destinationResDigit {
			return true
		}
	}

	// Otherwise, we have to determine the neighbor relationship the "hard" way.
	neighborRing := KRing(origin, 1)
	for i := 0; i < 7; i++ {
		if neighborRing[i] == destination {
			return true
		}
	}

	// Made it here, they definitely aren't neighbors
	return false
}

// GetH3UnidirectionalEdge returns a unidirectional edge H3 index based on the
// provided origin and destination.
func GetH3UnidirectionalEdge(origin H3Index, destination H3Index) H3Index {
	// Short-circuit and return an invalid index value if they are not neighbors
	if H3IndexesAreNeighbors(origin, destination) == false {
		return H3_NULL
	}

	// Otherwise, determine the IJK direction from the origin to the destination
	output := origin
	H3_SET_MODE(&output, H3_UNIEDGE_MODE)

	isPentagon := H3IsPentagon(origin)

	// Checks each neighbor, in order, to determine which direction the
	// destination neighbor is located. Skips CENTER_DIGIT since that
	// would be this index.
	var neighbor H3Index
	// Excluding from branch coverage as we never hit the end condition
	// LCOV_EXCL_BR_START
	direction := K_AXES_DIGIT
	if isPentagon {
		direction = J_AXES_DIGIT
	}

	for ; direction < Direction(NUM_DIGITS); direction++ {
		// LCOV_EXCL_BR_STOP
		rotations := 0
		neighbor = h3NeighborRotations(origin, direction, &rotations)
		if neighbor == destination {
			H3_SET_RESERVED_BITS(&output, int(direction))
			return output
		}
	}

	// This should be impossible, return H3_NULL in this case;
	return H3_NULL // LCOV_EXCL_LINE
}

// GetOriginH3IndexFromUnidirectionalEdge returns the origin hexagon from the
// unidirectional edge H3Index.
func GetOriginH3IndexFromUnidirectionalEdge(edge H3Index) H3Index {
	if H3_GET_MODE(edge) != H3_UNIEDGE_MODE {
		return H3_NULL
	}
	origin := edge
	H3_SET_MODE(&origin, H3_HEXAGON_MODE)
	H3_SET_RESERVED_BITS(&origin, 0)
	return origin
}

// GetDestinationH3IndexFromUnidirectionalEdge returns the destination hexagon
// from the unidirectional edge H3Index.
func GetDestinationH3IndexFromUnidirectionalEdge(edge H3Index) H3Index {
	if H3_GET_MODE(edge) != H3_UNIEDGE_MODE {
		return H3_NULL
	}
	direction := H3_GET_RESERVED_BITS(edge)
	rotations := 0
	destination := h3NeighborRotations(
		GetOriginH3IndexFromUnidirectionalEdge(edge), Direction(direction), &rotations)
	return destination
}

// H3UnidirectionalEdgeIsValid determines if the provided H3Index is a valid
// unidirectional edge index.
func H3UnidirectionalEdgeIsValid(edge H3Index) bool {
	if H3_GET_MODE(edge) != H3_UNIEDGE_MODE {
		return false
	}

	neighborDirection := H3_GET_RESERVED_BITS(edge)
	if neighborDirection <= int(CENTER_DIGIT) || neighborDirection >= NUM_DIGITS {
		return false
	}

	origin := GetOriginH3IndexFromUnidirectionalEdge(edge)
	if H3IsPentagon(origin) && neighborDirection == int(K_AXES_DIGIT) {
		return false
	}

	return H3IsValid(origin)
}

// GetH3IndexesFromUnidirectionalEdge returns the origin, destination pair of
// hexagon IDs for the given edge ID.
func GetH3IndexesFromUnidirectionalEdge(edge H3Index, originDestination *[]H3Index) {
	(*originDestination)[0] = GetOriginH3IndexFromUnidirectionalEdge(edge)
	(*originDestination)[1] = GetDestinationH3IndexFromUnidirectionalEdge(edge)
}

// GetH3UnidirectionalEdgesFromHexagon provides all of the unidirectional edges
// from the current H3Index.
func GetH3UnidirectionalEdgesFromHexagon(origin H3Index, edges *[]H3Index) {
	// Determine if the origin is a pentagon and special treatment needed.
	isPentagon := H3IsPentagon(origin)

	// This is actually quite simple. Just modify the bits of the origin
	// slightly for each direction, except the 'k' direction in pentagons,
	// which is zeroed.
	for i := 0; i < 6; i++ {
		if isPentagon && i == 0 {
			(*edges)[i] = H3_NULL
		} else {
			(*edges)[i] = origin
			H3_SET_MODE(&(*edges)[i], H3_UNIEDGE_MODE)
			H3_SET_RESERVED_BITS(&(*edges)[i], i+1)
		}
	}
}

// GetH3UnidirectionalEdgeBoundary provides the coordinates defining the
// unidirectional edge.
func GetH3UnidirectionalEdgeBoundary(edge H3Index, gb *GeoBoundary) {
	// Get the origin and neighbor direction from the edge
	direction := H3_GET_RESERVED_BITS(edge)
	origin := GetOriginH3IndexFromUnidirectionalEdge(edge)

	// Get the start vertex for the edge
	startVertex := vertexNumForDirection(origin, direction)
	if startVertex == INVALID_VERTEX_NUM {
		// This is not actually an edge (i.e. no valid direction),
		// so return no vertices.
		gb.numVerts = 0
		return
	}

	// Get the geo boundary for the appropriate vertexes of the origin. Note
	// that while there are always 2 topological vertexes per edge, the
	// resulting edge boundary may have an additional distortion vertex if it
	// crosses an edge of the icosahedron.
	var fijk FaceIJK
	_h3ToFaceIjk(origin, &fijk)
	res := H3_GET_RESOLUTION(origin)
	isPentagon := H3IsPentagon(origin)

	if isPentagon {
		_faceIjkPentToGeoBoundary(&fijk, res, startVertex, 2, gb)
	} else {
		_faceIjkToGeoBoundary(&fijk, res, startVertex, 2, gb)
	}
}
