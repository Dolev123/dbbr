package main

import (
	"fmt"
	"log"
	"strings"

	"github.com/Dolev123/dbbr"
)

type parser struct{}

func (p parser) ParseLine(line string) (dbbr.RuleLogic, error) {
	return func(i interface{}) bool {
		return strings.Contains(i.(string), line)
	}, nil
}

func check(s string, r dbbr.Rule) {
	fmt.Println(fmt.Sprintf("checking '%v' : %v", s, r.Check(s)))
}

func main() {
	rules, err := dbbr.ParseFile("demo2.dbbr", parser{})
	if err != nil {
		log.Fatalln(err)
		return
	}
	rule := rules["check_string"]
	check("hello", rule)
	check("hello world", rule)
	check("hello john", rule)
	check("hello world goodbye", rule)
	check("hello john goodbye", rule)
	check("hello world john", rule)
	check("", rule)
}
