package main

import (
	"errors"
	"fmt"
	"log"
	"strconv"
	"strings"

	"github.com/Dolev123/dbbr"
)

type parser struct{}

func (p parser) ParseLine(line string) (dbbr.RuleLogic, error) {
	parts := strings.Split(line, " ")
	if len(parts) == 2 {
		if parts[0] != "divBy" {
			return nil, errors.New(fmt.Sprintf("invalid word:", parts[0]))
		}
		num, err := strconv.Atoi(parts[1])
		if nil != err {
			return nil, err
		}
		return func(i interface{}) bool {
			return i.(int)%num == 0
		}, nil
	}
	if len(parts) == 3 {
		if parts[0] != "not" {
			return nil, errors.New(fmt.Sprintf("invalid word:", parts[0]))
		}
		if parts[1] != "divBy" {
			return nil, errors.New(fmt.Sprintf("invalid word:", parts[1]))
		}
		num, err := strconv.Atoi(parts[2])
		if nil != err {
			return nil, err
		}
		return func(i interface{}) bool {
			return i.(int)%num != 0
		}, nil
	}
	return nil, errors.New("could not parse line")
}

func main() {
	builder := dbbr.NewBuilder(parser{})
	lines, err := dbbr.ReadLinesFromFile("demo3.dbbr")
	if err != nil {
		log.Fatalln(err)
		return
	}
	sections := dbbr.SectionLines(lines)
	rules, err := builder.Build(sections)
	if err != nil {
		log.Fatalln(err)
		return
	}
	rule1 := rules["not_div_by_5_to_13"]
	rule2 := rules["power_of_2"]
	rule3 := rules["power_of_3"]
	rule4 := rules["check_number"]

	for i := 1; i <= 16; i += 1 {
		fmt.Println("number:", i, "|", rule1.Check(i), "|", rule2.Check(i), "|", rule3.Check(i), "|", rule4.Check(i))
	}
}
