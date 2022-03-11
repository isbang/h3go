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

// MAX_CELL_BNDRY_VERTS is maximum number of cell boundary vertices.
// Worst case is pentagon: 5 original verts + 5 edge crossings
const MAX_CELL_BNDRY_VERTS = 10

// GeoBoundary is cell boundary in latitude/longitude
type GeoBoundary struct {
	numVerts int                            // number of vertices
	verts    [MAX_CELL_BNDRY_VERTS]GeoCoord // vertices in ccw order
}
