// Package graph is a library for creating generic graph data structures and
// modifying, analyzing, and visualizing them.
//
// # Hashes
//
// A graph consists of vertices of type T, which are identified by a hash value
// of type K. The hash value for a given vertex is obtained using the hashing
// function passed to [graph.New]. A hashing function takes a T and returns a K.
//
// For primitive types like integers, you may use a predefined hashing function
// such as [graph.IntHash] – a function that takes an integer and uses that
// integer as the hash value at the same time:
//
//	g := graph.New(graph.IntHash)
//
// For storing custom data types, you need to provide your own hashing function.
// This example takes a City instance and returns its name as the hash value:
//
//	cityHash := func(c City) string {
//		return c.Name
//	}
//
// Creating a graph using this hashing function will yield a graph of vertices
// of type City identified by hash values of type string.
//
//	g := graph.New(cityHash)
//
// # Operations
//
// Adding vertices to a graph of integers is simple. [Graph.AddVertex] takes a
// vertex and adds it to the graph.
//
//	g := graph.New(graph.IntHash)
//
//	_ = g.AddVertex(1)
//	_ = g.AddVertex(2)
//
// Most functions accept and return only hash values instead of entire instances
// of the vertex type T. For example, [Graph.AddEdge] creates an edge between
// two vertices and accepts the hash values of those vertices. Because this
// graph uses the [graph.IntHash] hashing function, the vertex values and hash
// values are the same.
//
//	_ = g.AddEdge(1, 2)
//
// All operations that modify the graph itself are methods of [graph.Graph]. All
// other operations are top-level functions of by this library.
//
// For detailed usage examples, take a look at the README.
package graph

import "errors"

var (
	ErrVertexNotFound      = errors.New("vertex not found")
	ErrVertexAlreadyExists = errors.New("vertex already exists")
	ErrEdgeNotFound        = errors.New("edge not found")
	ErrEdgeAlreadyExists   = errors.New("edge already exists")
	ErrEdgeCreatesCycle    = errors.New("edge would create a cycle")
)

// Graph represents a generic graph data structure consisting of vertices of
// type T identified by a hash of type K.
type Graph[K comparable, T any] interface {
	// Traits returns the graph's traits. Those traits must be set when creating
	// a graph using New.
	Traits() *Traits

	// AddVertex creates a new vertex in the graph. If the vertex already exists
	// in the graph, ErrVertexAlreadyExists will be returned.
	//
	// AddVertex accepts a variety of functional options to set further edge
	// details such as the weight or an attribute:
	//
	//	_ = graph.AddVertex("A", "B", graph.VertexWeight(4), graph.VertexAttribute("label", "my-label"))
	//
	AddVertex(value T, options ...func(*VertexProperties)) error

	// Vertex returns the vertex with the given hash or ErrVertexNotFound if it
	// doesn't exist.
	Vertex(hash K) (T, error)

	// VertexWithProperties returns the vertex with the given hash along with
	// its properties or ErrVertexNotFound if it doesn't exist.
	VertexWithProperties(hash K) (T, VertexProperties, error)

	// AddEdge creates an edge between the source and the target vertex.
	//
	// If either vertex cannot be found, ErrVertexNotFound will be returned. If
	// the edge already exists, ErrEdgeAlreadyExists will be returned. If cycle
	// prevention has been activated using PreventCycles and if adding the edge
	// would create a cycle, ErrEdgeCreatesCycle will be returned.
	//
	// AddEdge accepts functional options to set further edge properties such as
	// the weight or an attribute:
	//
	//	_ = g.AddEdge("A", "B", graph.EdgeWeight(4), graph.EdgeAttribute("label", "my-label"))
	//
	AddEdge(sourceHash, targetHash K, options ...func(*EdgeProperties)) error

	// Edge returns the edge joining two given vertices or ErrEdgeNotFound if
	// the edge doesn't exist. In an undirected graph, an edge with swapped
	// source and target vertices does match.
	Edge(sourceHash, targetHash K) (Edge[T], error)

	// RemoveEdge removes the edge between the given source and target vertices.
	// If the edge cannot be found, ErrEdgeNotFound will be returned.
	RemoveEdge(source, target K) error

	// AdjacencyMap computes an adjacency map with all vertices in the graph.
	//
	// There is an entry for each vertex. Each of those entries is another map
	// whose keys are the hash values of the adjacent vertices. The value is an
	// Edge instance that stores the source and target hash values along with
	// the edge metadata.
	//
	// For a directed graph with two edges AB and AC, AdjacencyMap would return
	// the following map:
	//
	//	map[string]map[string]Edge[string]{
	//		"A": map[string]Edge[string]{
	//			"B": {Source: "A", Target: "B"},
	//			"C": {Source: "A", Target: "C"},
	//		},
	//		"B": map[string]Edge[string]{},
	//		"C": map[string]Edge[string]{},
	//	}
	//
	// This design makes AdjacencyMap suitable for a wide variety of algorithms.
	AdjacencyMap() (map[K]map[K]Edge[K], error)

	// PredecessorMap computes a predecessor map with all vertices in the graph.
	//
	// It has the same map layout and does the same thing as AdjacencyMap, but
	// for ingoing instead of outgoing edges of each vertex.
	//
	// For a directed graph with two edges AB and AC, PredecessorMap would
	// return the following map:
	//
	//	map[string]map[string]Edge[string]{
	//		"A": map[string]Edge[string]{},
	//		"B": map[string]Edge[string]{
	//			"A": {Source: "A", Target: "B"},
	//		},
	//		"C": map[string]Edge[string]{
	//			"A": {Source: "A", Target: "C"},
	//		},
	//	}
	//
	// For an undirected graph, PredecessorMap is the same as AdjacencyMap. This
	// is because there is no distinction between "outgoing" and "ingoing" edges
	// in an undirected graph.
	PredecessorMap() (map[K]map[K]Edge[K], error)

	// Clone creates a deep copy of the graph and returns that cloned graph.
	Clone() (Graph[K, T], error)

	// Order returns the number of vertices in the graph.
	Order() (int, error)

	// Size returns the number of edges in the graph.
	Size() (int, error)
}

// Edge represents an edge that joins two vertices. Even though these edges are
// always referred to as source and target, whether the graph is directed or not
// is determined by its traits.
type Edge[T any] struct {
	Source     T
	Target     T
	Properties EdgeProperties
}

// EdgeProperties represents a set of properties that each edge possesses. They
// can be set when adding a new edge using the corresponding functional options:
//
//	g.AddEdge("A", "B", graph.EdgeWeight(2), graph.EdgeAttribute("color", "red"))
//
// The example above will create an edge with a weight of 2 and an attribute
// "color" with value "red".
type EdgeProperties struct {
	Attributes map[string]string
	Weight     int
	Data       any
}

// Hash is a hashing function that takes a vertex of type T and returns a hash
// value of type K.
//
// Every graph has a hashing function and uses that function to retrieve the
// hash values of its vertices. You can either use one of the predefined hashing
// functions or provide your own one for custom data types:
//
//	cityHash := func(c City) string {
//		return c.Name
//	}
//
// The cityHash function returns the city name as a hash value. The types of T
// and K, in this case City and string, also define the types of the graph.
type Hash[K comparable, T any] func(T) K

// New creates a new graph with vertices of type T, identified by hash values of
// type K. These hash values will be obtained using the provided hash function.
//
// The graph will use the default in-memory store for persisting vertices and
// edges. To use a different [graph.Store], use [graph.NewWithStore].
func New[K comparable, T any](hash Hash[K, T], options ...func(*Traits)) Graph[K, T] {
	return NewWithStore(hash, newMemoryStore[K, T](), options...)
}

// NewWithStore creates a new graph same as New but uses the provided store
// instead of the default memory store.
func NewWithStore[K comparable, T any](hash Hash[K, T], store Store[K, T], options ...func(*Traits)) Graph[K, T] {
	var p Traits

	for _, option := range options {
		option(&p)
	}

	if p.IsDirected {
		return newDirected(hash, &p, store)
	}

	return newUndirected(hash, &p, store)
}

// StringHash is a hashing function that accepts a string and uses that exact
// string as a hash value. Using it as Hash will yield a Graph[string, string].
func StringHash(v string) string {
	return v
}

// IntHash is a hashing function that accepts an integer and uses that exact
// integer as a hash value. Using it as Hash will yield a Graph[int, int].
func IntHash(v int) int {
	return v
}

// EdgeWeight returns a function that sets the weight of an edge to the given
// weight. This is a functional option for the [Graph.Edge] and [Graph.AddEdge]
// methods.
func EdgeWeight(weight int) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Weight = weight
	}
}

// EdgeAttribute returns a function that adds the given key-value pair to the
// attributes of an edge. This is a functional option for the [Graph.Edge] and
// [Graph.AddEdge] methods.
func EdgeAttribute(key, value string) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Attributes[key] = value
	}
}

// EdgeData returns a function that sets the data of an edge to the given value.
// This is a functional option for the [Graph.Edge] and [Graph.AddEdge] methods.
func EdgeData(data any) func(*EdgeProperties) {
	return func(e *EdgeProperties) {
		e.Data = data
	}
}

// VertexProperties represents a set of properties that each vertex has. They
// can be set when adding a vertex using the corresponding functional options:
//
//	_ = g.AddVertex("A", "B", graph.VertexWeight(2), graph.VertexAttribute("color", "red"))
//
// The example above will create a vertex with a weight of 2 and an attribute
// "color" with value "red".
type VertexProperties struct {
	Attributes map[string]string
	Weight     int
}

// VertexWeight returns a function that sets the weight of a vertex to the given
// weight. This is a functional option for the [Graph.Vertex] and
// [Graph.AddVertex] methods.
func VertexWeight(weight int) func(*VertexProperties) {
	return func(e *VertexProperties) {
		e.Weight = weight
	}
}

// VertexAttribute returns a function that adds the given key-value pair to the
// vertex attributes. This is a functional option for the [Graph.Vertex] and
// [Graph.AddVertex] methods.
func VertexAttribute(key, value string) func(*VertexProperties) {
	return func(e *VertexProperties) {
		e.Attributes[key] = value
	}
}
