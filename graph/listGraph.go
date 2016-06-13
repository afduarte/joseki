package graph

import "github.com/Callidon/joseki/rdf"

// ListGraph is dummy implementation of a RDF Graph, using a simple slice to store RDF Triples.
//
// Very poorly optimized, should only be used for demonstration or benchmarking purposes.
type ListGraph struct {
	triples []rdf.Triple
}

// NewListGraph creates a new List Graph.
func NewListGraph() ListGraph {
	return ListGraph{make([]rdf.Triple, 0)}
}

// LoadFromFile load the content of a RDF graph stored in a file into the current graph.
func (g *ListGraph) LoadFromFile(filename, format string) {
	//TODO
}

// Add a new Triple pattern to the graph.
func (g *ListGraph) Add(triple rdf.Triple) {
	g.triples = append(g.triples, triple)
}

// Filter fetch triples form the graph that match a BGP given in parameters.
func (g *ListGraph) Filter(subject, predicate, object rdf.Node) chan rdf.Triple {
	results := make(chan rdf.Triple)
	refTriple := rdf.NewTriple(subject, predicate, object)
	// search for matching triple pattern in graph
	go func() {
		for _, triple := range g.triples {
			test, err := refTriple.Equivalent(triple)
			if (err == nil) && test {
				results <- triple
			}
		}
		close(results)
	}()
	return results
}

// Serialize the graph into a given format and return it as a string.
func (g *ListGraph) Serialize(format string) string {
	// TODO
	return ""
}
