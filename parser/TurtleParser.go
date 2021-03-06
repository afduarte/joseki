// Copyright (c) 2016 Thomas Minier. All rights reserved.
// Use of this source code is governed by a MIT License
// license that can be found in the LICENSE file.

package parser

import (
	"bufio"
	"github.com/Callidon/joseki/rdf"
	"io"
	"os"
	"strings"
)

// TurtleParser is a parser for reading & loading triples in Turtle format.
//
// Turtle reference : https://www.w3.org/TR/turtle/
type TurtleParser struct {
	prefixes map[string]string
	cutter   *lineCutter
}

// scanTurtle read a file in Turtle format, identify and extract token with their values.
//
// The results are sent through a channel, which is closed when the scan of the file has been completed.
func scanTurtle(reader io.Reader, out chan<- rdfToken, l *lineCutter) {
	// walk through the file using a goroutine
	go func() {
		defer close(out)
		var prefixName, prefixValue string
		var scanPrefixesDone bool

		scanner := bufio.NewScanner(reader)
		lineNumber := 1
		for scanner.Scan() {
			line := l.extractSegments(scanner.Text())
			rowNumber := 1
			// skip blank lines & comments
			if (len(line) == 0) || (line[0] == "#") {
				lineNumber++
				continue
			}
			scanPrefixesDone = (line[0] != "@prefix")
			// scan elements of the line
			for _, elt := range line {
				// skip comments
				if string(elt[0]) == "#" {
					break
				}
				if !scanPrefixesDone {
					switch {
					case elt == "@prefix" || elt == ":":
						continue
					case elt == ".":
						out <- newTokenPrefix(prefixName, prefixValue)
						prefixName, prefixValue = "", ""
					case prefixName == "":
						if string(elt[len(elt)-1]) != ":" {
							out <- newTokenIllegal("Unexpected token : "+elt, lineNumber, rowNumber)
							return
						}
						prefixName = elt[0 : len(elt)-1]
					case prefixValue == "":
						if string(elt[0]) != "<" && string(elt[len(elt)-1]) != ">" {
							out <- newTokenIllegal("Unexpected token : "+elt, lineNumber, rowNumber)
							return
						}
						prefixValue = elt[1 : len(elt)-1]
					default:
						out <- newTokenIllegal("Unexpected token when scanning '"+elt+"', expected a prefix definition", lineNumber, rowNumber)
					}
				} else {
					switch {
					case elt == ".", elt == "]":
						out <- newTokenEnd(lineNumber, rowNumber)
					case elt == ";", elt == ",", elt == "[":
						out <- newTokenSep(elt, lineNumber, rowNumber)
					case string(elt[0]) == "<" && string(elt[len(elt)-1]) == ">":
						out <- newTokenURI(elt[1 : len(elt)-1])
					case string(elt[0]) == "\"" && string(elt[len(elt)-1]) == "\"", string(elt[0]) == "'" && string(elt[len(elt)-1]) == "'":
						out <- newTokenLiteral(elt[1 : len(elt)-1])
					case len(elt) >= 2 && elt[0:2] == "^^":
						out <- newTokenType(elt[2:], lineNumber, rowNumber)
					case string(elt[0]) == "@":
						out <- newTokenLang(elt[1:], lineNumber, rowNumber)
					case string(elt[0]) == "_" && string(elt[1]) == ":":
						out <- newTokenBlankNode(elt[2:])
					case string(elt[0]) == "?":
						out <- newTokenBlankNode(elt[1:])
					case strings.Index(elt, ":") > -1:
						out <- newTokenPrefixedURI(elt, lineNumber, rowNumber)
					default:
						out <- newTokenIllegal("Unexpected token when scanning '"+elt+"'", lineNumber, rowNumber)
					}
				}
				rowNumber += len(elt) + 1
			}
			lineNumber++
		}
	}()
}

// NewTurtleParser creates a new TurtleParser
func NewTurtleParser() *TurtleParser {
	return &TurtleParser{make(map[string]string), newLineCutter(wordRegexp)}
}

// Prefixes returns the prefixes read by the parser during the last parsing.
func (p TurtleParser) Prefixes() map[string]string {
	return p.prefixes
}

// Read a file containg RDF triples in Turtle format & convert them in triples.
//
// Triples generated are send throught a channel, which is closed when the parsing of the file has been completed.
func (p *TurtleParser) Read(filename string) chan rdf.Triple {
	tokenPipe := make(chan rdfToken, bufferSize)
	out := make(chan rdf.Triple, bufferSize)
	stack := newStack()

	// scan the file & analyse the tokens using a goroutine
	go func() {
		defer close(out)
		f, err := os.Open(filename)
		check(err)
		defer f.Close()
		// launch the scan, then interpret each token produced
		go scanTurtle(bufio.NewReader(f), tokenPipe, p.cutter)
		for token := range tokenPipe {
			err = token.Interpret(stack, &p.prefixes, out)
			check(err)
		}
	}()
	return out
}
