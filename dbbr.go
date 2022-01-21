// dbbr Implements the DBBR format for binary tree creation.
// It is used to create simple (and even complex) true/false rules
// from a true/false logic tree.
//
// For more info, look at the README.md file.
package dbbr

import (
	"errors"
)

// ParseFile creates rules from a given file, using the given parser.
func ParseFile(path string, parser LogicParser) (map[string]Rule, error) {
	rules := map[string]Rule{}
	if parser == nil {
		return rules, errors.New("cannot create a builder with a nil parser")
	}
	lines, err := ReadLinesFromFile(path)
	if err != nil {
		return rules, err
	}
	sections := SectionLines(lines)
	builder := NewBuilder(parser)
	rules, err = builder.Build(sections)
	if err != nil {
		return rules, err
	}
	return rules, nil
}
