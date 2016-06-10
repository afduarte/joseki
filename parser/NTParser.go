package parser

import (
	"bufio"
	"errors"
	"github.com/Callidon/joseki/rdf"
	"os"
	"strings"
)

type NTParser struct {
}

func check(err error) {
	if err != nil {
		panic(err)
	}
}

func parseNode(elt string) (rdf.Node, error) {
	var node rdf.Node
	var err error
	if (string(elt[0]) == "<") && (string(elt[len(elt)-1]) == ">") {
		node = rdf.NewURI(elt[1 : len(elt)-2])
	} else if (string(elt[0]) == "\"") && (string(elt[len(elt)-1]) == "\"") {
		// TODO add a security when a xml type is given with the literal
		node = rdf.NewLiteral(elt[1 : len(elt)-2])
	} else if (string(elt[0]) == "_") && (string(elt[1]) == ":") {
		node = rdf.NewBlankNode(elt[2:])
	} else {
		err = errors.New("Error : cannot parse " + elt + " into a RDF Node")
	}
	return node, err
}

func (p *NTParser) Read(filename string) chan rdf.Triple {
	out := make(chan rdf.Triple)
	// walk through the file using a goroutine
	go func() {
		f, err := os.Open(filename)
		check(err)
		defer f.Close()
		scanner := bufio.NewScanner(bufio.NewReader(f))
		for scanner.Scan() {
			elts := strings.Split(scanner.Text(), " ")
			subject, err := parseNode(elts[0])
			check(err)
			predicate, err := parseNode(elts[1])
			check(err)
			object, err := parseNode(elts[2])
			check(err)
			out <- rdf.NewTriple(subject, predicate, object)
		}
		close(out)
	}()
	return out
}