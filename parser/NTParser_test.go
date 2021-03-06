// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"github.com/Callidon/joseki/rdf"
	"strings"
	"testing"
)

func TestReadNTParser(t *testing.T) {
	parser := NewNTParser()
	cpt := 0
	datas := []rdf.Triple{
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/Document")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://www.w3.org/1999/02/22-rdf-syntax-ns#type"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/Document")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://purl.org/dc/terms/title"),
			rdf.NewLangLiteral("N-Triples", "en")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://purl.org/dc/terms/title"),
			rdf.NewTypedLiteral("My Typed Literal", "<http://www.w3.org/2001/XMLSchema#string>")),
		rdf.NewTriple(rdf.NewURI("http://www.w3.org/2001/sw/RDFCore/ntriples"),
			rdf.NewURI("http://xmlns.com/foaf/0.1/maker"),
			rdf.NewBlankNode("art")),
	}

	for elt := range parser.Read("datas/test.nt") {
		if test, err := elt.Equals(datas[cpt]); !test || (err != nil) {
			t.Error(datas[cpt], "should be equal to", elt)
		}
		cpt++
	}

	if cpt != len(datas) {
		t.Error("read", cpt, "nodes of the file instead of", len(datas))
	}
}

func TestIllegalTokenNTParser(t *testing.T) {
	input := "illegal_token"
	expectedMsg := "Unexpected token when scanning 'illegal_token' at line : 1 row : 1"
	out := make(chan rdfToken, bufferSize)
	scanNtriples(strings.NewReader(input), out, newLineCutter(wordRegexp))

	token := <-out
	tokenErr := token.Interpret(nil, nil, nil).Error()
	if tokenErr != expectedMsg {
		t.Error("expected illegal token", expectedMsg, "but instead got", tokenErr)
	}
}
