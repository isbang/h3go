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

// VertexNode is a single node in a vertex graph, part of a linked list.
type VertexNode struct {
	from GeoCoord
	to   GeoCoord
	next *VertexNode
}

// VertexGraph is a data structure to store a graph of vertices
type VertexGraph struct {
	buckets    []*VertexNode
	numBuckets int
	size       int
	res        int
}

/**
 * Initialize a new VertexGraph
 * @param graph       Graph to initialize
 * @param  numBuckets Number of buckets to include in the graph
 * @param  res        Resolution of the hexagons whose vertices we're storing
 */
func initVertexGraph(graph *VertexGraph, numBuckets int, res int) {
	if numBuckets > 0 {
		graph.buckets = make([]*VertexNode, numBuckets)
	} else {
		graph.buckets = nil
	}

	graph.numBuckets = numBuckets
	graph.size = 0
	graph.res = res
}

/**
 * Destroy a VertexGraph's sub-objects, freeing their memory. The caller is
 * responsible for freeing memory allocated to the VertexGraph struct itself.
 * @param graph Graph to destroy
 */
func destroyVertexGraph(graph *VertexGraph) {
	for {
		node := firstVertexNode(graph)
		if node == nil {
			break
		}
		removeVertexNode(graph, node)
	}
	graph.buckets = nil
}

/**
 * Get an integer hash for a lat/lon point, at a precision determined
 * by the current hexagon resolution.
 * TODO: Light testing suggests this might not be sufficient at resolutions
 * finer than 10. Design a better hash function if performance and collisions
 * seem to be an issue here.
 * @param  vertex     Lat/lon vertex to hash
 * @param  res        Resolution of the hexagon the vertex belongs to
 * @param  numBuckets Number of buckets in the graph
 * @return            Integer hash
 */
func _hashVertex(vertex *GeoCoord, res int, numBuckets int) uint32 {
	// Simple hash: Take the sum of the lat and lon with a precision level
	// determined by the resolution, converted to int, modulo bucket count.
	return uint32(
		math.Mod(
			math.Abs(
				(vertex.lat+vertex.lon)*math.Pow(10, float64(15-res)),
			),
			float64(numBuckets),
		),
	)
}

func _initVertexNode(fromVtx *GeoCoord, toVtx *GeoCoord) *VertexNode {
	return &VertexNode{
		from: *fromVtx,
		to:   *toVtx,
		next: nil,
	}
}

/**
 * Add a edge to the graph
 * @param graph   Graph to add node to
 * @param fromVtx Start vertex
 * @param toVtx   End vertex
 * @return        Pointer to the new node
 */
func addVertexNode(graph *VertexGraph, fromVtx *GeoCoord, toVtx *GeoCoord) *VertexNode {
	// Make the new node
	node := _initVertexNode(fromVtx, toVtx)
	// Determine location
	index := _hashVertex(fromVtx, graph.res, graph.numBuckets)

	// Check whether there's an existing node in that spot
	currentNode := graph.buckets[index]
	if currentNode == nil {
		// Set bucket to the new node
		graph.buckets[index] = node
	} else {
		// Find the end of the list
		for {
			// Check the the edge we're adding doesn't already exist
			if geoAlmostEqual(&currentNode.from, fromVtx) &&
				geoAlmostEqual(&currentNode.to, toVtx) {
				// already exists, bail
				return currentNode
			}
			if currentNode.next != nil {
				currentNode = currentNode.next
			}

			if currentNode.next == nil {
				break
			}
		}
		// Add the new node to the end of the list
		currentNode.next = node
	}
	graph.size++
	return node
}

/**
 * Remove a node from the graph. The input node will be freed, and should
 * not be used after removal.
 * @param graph Graph to mutate
 * @param node  Node to remove
 * @return      0 on success, 1 on failure (node not found)
 */
func removeVertexNode(graph *VertexGraph, node *VertexNode) bool {
	// Determine location
	index := _hashVertex(&node.from, graph.res, graph.numBuckets)
	currentNode := graph.buckets[index]
	found := false
	if currentNode != nil {
		if currentNode == node {
			graph.buckets[index] = node.next
			found = true
		}
		// Look through the list
		for !found && currentNode.next != nil {
			if currentNode.next == node {
				// splice the node out
				currentNode.next = node.next
				found = true
			}
			currentNode = currentNode.next
		}
	}
	if found {
		node = nil
		graph.size--
		return false
	}
	// Failed to find the node
	return true
}

/**
 * Find the Vertex node for a given edge, if it exists
 * @param  graph   Graph to look in
 * @param  fromVtx Start vertex
 * @param  toVtx   End vertex, or NULL if we don't care
 * @return         Pointer to the vertex node, if found
 */
func findNodeForEdge(graph *VertexGraph, fromVtx *GeoCoord, toVtx *GeoCoord) *VertexNode {
	// Determine location
	index := _hashVertex(fromVtx, graph.res, graph.numBuckets)
	// Check whether there's an existing node in that spot
	node := graph.buckets[index]
	if node != nil {
		// Look through the list and see if we find the edge
		for {
			if geoAlmostEqual(&node.from, fromVtx) &&
				(toVtx == nil || geoAlmostEqual(&node.to, toVtx)) {
				return node
			}
			node = node.next

			if node == nil {
				break
			}
		}
	}
	// Iteration lookup fail
	return nil
}

/**
 * Find a Vertex node starting at the given vertex
 * @param  graph   Graph to look in
 * @param  fromVtx Start vertex
 * @return         Pointer to the vertex node, if found
 */
func findNodeForVertex(graph *VertexGraph, fromVtx *GeoCoord) *VertexNode {
	return findNodeForEdge(graph, fromVtx, nil)
}

/**
 * Get the next vertex node in the graph.
 * @param  graph Graph to iterate
 * @return       Vertex node, or NULL if at the end
 */
func firstVertexNode(graph *VertexGraph) *VertexNode {
	for _, node := range graph.buckets {
		if node != nil {
			return node
		}
	}

	return nil
}
