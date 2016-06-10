// Package joseki/rdf provides primitives to work with RDF
package rdf

import "errors"

// Interface which represent a generic node in a RDF Graph
type Node interface {
	Equals(n Node) (bool, error)
	Equivalent(n Node) (bool, error)
	String() string
}

// Type which represent a URI Node in a RDF Graph.
// RDF URI reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-Graph-URIref
type URI struct {
	Value string
}

// Type which represent a Literal Node in a RDF Graph.
// RDF Literal reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-Graph-Literal
type Literal struct {
	Value string
}

// Type which represent a Blank Node in a RDF Graph.
// RDF Blank Node reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-blank-nodes
type BlankNode struct {
	Variable string
}

// Type which represent a RDF triple.
// RDF Triple reference : https://www.w3.org/TR/2004/REC-rdf-concepts-20040210/#section-triples
type Triple struct {
	Subject   Node
	Predicate Node
	Object    Node
}

// Return True if Two URIs are equals, False if not
func (u URI) Equals(n Node) (bool, error) {
	other, ok := n.(URI)
	if ok {
		return u.Value == other.Value, nil
	} else {
		return false, errors.New("Error : mismatch type, can only compare two URIs")
	}
}

// Test if a URI is equivalent to a Node, assuming that a URI and a Blank Node are equals, like in the context of a SPARQL Query.s
// Return True if the two URIs are equivalent with this criteria, False if not
func (u URI) Equivalent(n Node) (bool, error) {
	equality, err := u.Equals(n)
	if err != nil {
		_, ok := n.(BlankNode)
		if ok {
			return true, nil
		} else {
			return false, errors.New("Error : can only compare a URI with another URI or a Blank Node")
		}
	} else {
		return equality, nil
	}
}

// Serialize a URI to string and return it
func (u URI) String() string {
	return "<" + u.Value + ">"
}

// Create a new URI
func NewURI(value string) URI {
	return URI{value}
}

// Return True if Two Literals are strictly equals, False if not
func (l Literal) Equals(n Node) (bool, error) {
	other, ok := n.(Literal)
	if ok {
		return l.Value == other.Value, nil
	} else {
		return false, errors.New("Error : mismatch type, can only compare two Literals")
	}
}

// Test if a Literal is equivalent to a Node, assuming that a Literal and a Blank Node are equals, like in the context of a SPARQL Query.
// Return True if the two nodes are equivalent with this criteria, False if not
func (l Literal) Equivalent(n Node) (bool, error) {
	equality, err := l.Equals(n)
	if err != nil {
		_, ok := n.(BlankNode)
		if ok {
			return true, nil
		} else {
			return false, errors.New("Error : can only compare a Literal with another Literal or a Blank Node")
		}
	} else {
		return equality, nil
	}
}

// Serialize a Literal to string and return it
func (l Literal) String() string {
	return "\"" + l.Value + "\""
}

// Create a new Literal
func NewLiteral(value string) Literal {
	return Literal{value}
}

// Return True if Two Blank Node are strictly equals, False if not
func (b BlankNode) Equals(n Node) (bool, error) {
	other, ok := n.(BlankNode)
	if ok {
		return b.Variable == other.Variable, nil
	} else {
		return false, errors.New("Error : mismatch type, can only compare two Blank Nodes")
	}
}

// Always return true, assuming that a Blank Node is equivalent to any other node in a SPARQL Query
func (b BlankNode) Equivalent(n Node) (bool, error) {
	return true, nil
}

// Serialize a Blank Node to string and return it
func (b BlankNode) String() string {
	return "_:" + b.Variable
}

// Create a new Literal
func NewBlankNode(variable string) BlankNode {
	return BlankNode{variable}
}

// Return True if two triples are strictly equals, False if not
func (t Triple) Equals(other Triple) (bool, error) {
	test_subj, err := t.Subject.Equals(other.Subject)
	if err != nil {
		return false, err
	}
	test_pred, err := t.Predicate.Equals(other.Predicate)
	if err != nil {
		return false, err
	}
	test_obj, err := t.Object.Equals(other.Object)
	if err != nil {
		return false, err
	}
	return test_subj && test_pred && test_obj, nil
}

// Test if a Triple is equivalent to another triple, assuming that blank node are equals to any other node types.
// Return True if the two triples are equivalent with this criteria, False if not
func (t Triple) Equivalent(other Triple) (bool, error) {
	test_subj, err := t.Subject.Equivalent(other.Subject)
	if err != nil {
		return false, err
	}
	test_pred, err := t.Predicate.Equivalent(other.Predicate)
	if err != nil {
		return false, err
	}
	test_obj, err := t.Object.Equivalent(other.Object)
	if err != nil {
		return false, err
	}
	return test_subj && test_pred && test_obj, nil
}

// Create a new Triple
func NewTriple(subject Node, predicate Node, object Node) Triple {
	return (Triple{subject, predicate, object})
}