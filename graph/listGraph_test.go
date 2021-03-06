// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package graph

import (
	"github.com/Callidon/joseki/rdf"
	"math/rand"
	"testing"
)

func TestAddListGraph(t *testing.T) {
	graph := NewListGraph()
	subj := rdf.NewURI("http://dbpl.org#Thomas")
	pred := rdf.NewURI("http://foaf.com/age")
	obj := rdf.NewLiteral("22")
	triple := rdf.NewTriple(subj, pred, obj)
	graph.Add(triple)
	graphTriple, _ := graph.triples[0].Triple(graph.dictionnary)
	if test, err := graphTriple.Equals(triple); !test && (err != nil) {
		t.Error(triple, "hasn't been inserted into the graph")
	}
}

func TestFilterListGraph(t *testing.T) {
	skipTest("./watdiv.30k.nt", t)
	graph := NewListGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	subj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Offer10001")
	pred := rdf.NewURI("http://schema.org/eligibleRegion")
	obj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Country0")
	triple := rdf.NewTriple(subj, pred, obj)
	cpt := 0

	// select one triple specific triple pattern
	for result := range graph.Filter(subj, pred, obj) {
		if test, err := result.Equals(triple); !test || (err != nil) {
			t.Error("expected", triple, "but instead got", result)
		}
		cpt++
	}

	if cpt != 1 {
		t.Error("expected 1 result but instead got", cpt, "results")
	}

	// select all triples
	cpt = 0
	for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z")) {
		cpt++
	}
	if cpt != 30000 {
		t.Error("expected 30000 results but instead got", cpt, "results")
	}

	// select multiple triples with the same subject
	cpt = 0
	for _ = range graph.Filter(rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Offer1375"), rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}
	if cpt != 9 {
		t.Error("expected 9 results but instead got", cpt, "results")
	}

	// select multiple triples with the same predicate
	cpt = 0
	for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewURI("http://www.geonames.org/ontology#parentCountry"), rdf.NewVariable("w")) {
		cpt++
	}
	if cpt != 240 {
		t.Error("expected 240 results but instead got", cpt, "results")
	}

	// select multiple triples with the same object
	cpt = 0
	for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewLiteral("673")) {
		cpt++
	}
	if cpt != 6 {
		t.Error("expected 6 results but instead got", cpt, "results")
	}

	// select a triple that doesn't exist in the graph
	cpt = 0
	for _ = range graph.Filter(rdf.NewURI("http://example.org"), rdf.NewVariable("v1"), rdf.NewVariable("v2")) {
		cpt++
	}

	if cpt > 0 {
		t.Error("expected no result but instead found", cpt, "results")
	}
}

func TestFilterSubsetListGraph(t *testing.T) {
	skipTest("./watdiv.30k.nt", t)
	graph := NewListGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	nbDatas, limit, offset := 30000, 600, 800
	cpt := 0

	// test a FilterSubset with a simple Limit
	for _ = range graph.FilterSubset(rdf.NewVariable("x"), rdf.NewVariable("v"), rdf.NewVariable("w"), limit, -1) {
		cpt++
	}

	if cpt != limit {
		t.Error("expected ", limit, "results but instead found ", cpt, "results")
	}

	// test a FilterSubset with a simple offset
	cpt = 0
	for _ = range graph.FilterSubset(rdf.NewVariable("x"), rdf.NewVariable("v"), rdf.NewVariable("w"), -1, offset) {
		cpt++
	}

	if cpt != nbDatas-offset {
		t.Error("expected ", nbDatas-offset, "results but instead found ", cpt, "results")
	}

	// test with a offset than doesn't allow enough results to reach the limit
	cpt = 0
	offset = nbDatas - 10
	for _ = range graph.FilterSubset(rdf.NewVariable("x"), rdf.NewVariable("v"), rdf.NewVariable("w"), limit, offset) {
		cpt++
	}

	if cpt != nbDatas-offset {
		t.Error("expected ", nbDatas-offset, "results but instead found ", cpt, "results")
	}
}

func TestDeleteListGraph(t *testing.T) {
	var triple rdf.Triple
	graph := NewListGraph()
	nbDatas := 1000
	cpt := 0
	subj := rdf.NewURI("http://dblp.com#foo")

	// insert random triples in the graph
	for i := 0; i < nbDatas; i++ {
		triple = rdf.NewTriple(subj, rdf.NewURI(string(rand.Intn(nbDatas))), rdf.NewLiteral(string(rand.Intn(nbDatas))))
		graph.Add(triple)
	}

	// remove the last triple pattern inserted
	graph.Delete(triple.Subject, triple.Predicate, triple.Object)
	for _ = range graph.Filter(triple.Subject, triple.Predicate, triple.Object) {
		cpt++
	}

	if cpt > 0 {
		t.Error("the graph shouldn't contains the triple", triple)
	}

	// remove all triple with a given subject
	graph.Delete(subj, rdf.NewVariable("v"), rdf.NewVariable("w"))

	// select all triple of the graph
	cpt = 0
	for _ = range graph.Filter(subj, rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}

	if cpt > 0 {
		t.Error("the graph should be empty")
	}
}

func TestLoadFromFileListGraph(t *testing.T) {
	graph := NewListGraph()
	cpt := 0
	graph.LoadFromFile("../parser/datas/test.nt", "nt")

	// select all triple of the graph
	for _ = range graph.Filter(rdf.NewVariable("y"), rdf.NewVariable("v"), rdf.NewVariable("w")) {
		cpt++
	}

	if cpt != 5 {
		t.Error("the graph should contains 5 triples, but it contains", cpt, "triples")
	}
}

// Benchmarking with WatDiv 1K

func BenchmarkAddListGraph(b *testing.B) {
	b.Skip("skipped because it's currently not accurate")
	graph := NewListGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	triple := rdf.NewTriple(rdf.NewURI("http://example.org/subject"), rdf.NewURI("http://example.org/predicate"), rdf.NewURI("http://example.org/object"))
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		graph.Add(triple)
	}
}

func BenchmarkLoadFromFileListGraph(b *testing.B) {
	for i := 0; i < b.N; i++ {
		graph := NewListGraph()
		graph.LoadFromFile("./watdiv.30k.nt", "nt")
	}
}

func BenchmarkAllFilterListGraph(b *testing.B) {
	graph := NewListGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	cpt := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// select all triple of the graph
		for _ = range graph.Filter(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z")) {
			cpt++
		}
	}
}

func BenchmarkSpecificFilterListGraph(b *testing.B) {
	graph := NewListGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	subj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Offer10001")
	pred := rdf.NewURI("http://schema.org/eligibleRegion")
	obj := rdf.NewURI("http://db.uwaterloo.ca/~galuc/wsdbm/Country0")
	cpt := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// fetch the last inserted triple into the graph
		for _ = range graph.Filter(subj, pred, obj) {
			cpt++
		}
	}
}

func BenchmarkAllFilterSubsetListGraph(b *testing.B) {
	graph := NewListGraph()
	graph.LoadFromFile("./watdiv.30k.nt", "nt")
	limit, offset := 600, 200
	cpt := 0
	b.ResetTimer()

	for i := 0; i < b.N; i++ {
		// select all triple of the graph
		for _ = range graph.FilterSubset(rdf.NewVariable("v"), rdf.NewVariable("w"), rdf.NewVariable("z"), limit, offset) {
			cpt++
		}
	}
}
