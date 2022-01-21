package dbbr

import (
	"errors"
	"fmt"
	"strings"
)

type BuildTree struct {
	Name  string
	True  *BuildTree
	False *BuildTree
	Logic string
	Line  int
	// if set to true, use Logic as rule name
	IsPointer bool
}

func NewBuildTree(name, logic string, line int) *BuildTree {
	return &BuildTree{name, nil, nil, logic, line, false}
	// return &BuildTree{name, True, False, logic, line, false}
}

func (bt *BuildTree) UpdateNext(b bool, next *BuildTree) {
	if b {
		bt.True = next
	} else {
		bt.False = next
	}
}

// Builder builds
type Builder struct {
	//	ruleNames []string
	rules  map[string]Rule
	parser LogicParser
}

type indentTuple struct {
	bt    *BuildTree
	count int
}

// NewBuilder creates a Builder with builtin True and False.
func NewBuilder(parser LogicParser) *Builder {
	rules := make(map[string]Rule)
	rules["True"] = True
	rules["False"] = False
	return &Builder{
		rules:  rules,
		parser: parser,
	}
}

// Build builds Rules from a sections of lines, each slice corresponds to 1 rule.
func (b *Builder) Build(sections [][]Line) (map[string]Rule, error) {
	ruleTrees := make([]*BuildTree, 0)
	for _, lines := range sections {
		bt, err := b.BuildRuleTree(lines)
		if err != nil {
			return nil, err
		}
		ruleTrees = append(ruleTrees, bt)
	}

	for _, bt := range ruleTrees {
		if bt == nil {
			continue
		}
		newRule, err := b.walkBuildTree(bt)
		if err != nil {
			return b.rules, err
		}
		b.rules[newRule.Name()] = newRule
	}
	return b.rules, nil
}

// create rules over binary tree recursively
func (b *Builder) walkBuildTree(bt *BuildTree) (Rule, error) {
	var (
		logic RuleLogic
		err   error
	)
	if !bt.IsPointer {
		logic, err = b.parser.ParseLine(bt.Logic)
		if err != nil {
			return nil, err
		}
	} else {
		logic = func(i interface{}) bool {
			return b.rules[bt.Logic].Check(i)
		}
	}
	nTrue := True
	nFalse := False
	if bt.False != nil {
		nFalse, err = b.walkBuildTree(bt.False)
		if err != nil {
			return nil, err
		}
	}
	if bt.True != nil {
		nTrue, err = b.walkBuildTree(bt.True)
		if err != nil {
			return nil, err
		}
	}
	return NewRule(bt.Name, logic, nTrue, nFalse), nil
}

func CalcIndent(s string) int {
	return len(s) - len(strings.TrimLeft(s, "\t"))
}

// BulidRuleTree builds a single BuildTree for a rule, from a section of lines.
func (b *Builder) BuildRuleTree(lines []Line) (base *BuildTree, err error) {
	if len(lines) == 0 {
		return nil, nil
	}

	currIndent := 0
	indents := []indentTuple{}
	var name string
	for _, line := range lines {
		indent := CalcIndent(line.Value)
		var currNode *BuildTree
		if indent == 0 {
			// comment
			if strings.TrimSpace(line.Value)[0] == '#' {
				continue
			}
			// base node of rule
			if len(indents) != 0 {
				return nil, errors.New(fmt.Sprintf("base already exists for rule, error at line %v", line.Number))
			}
			parts := strings.SplitN(line.Value, ":", 2)
			if len(parts) != 2 {
				return nil, errors.New(fmt.Sprintf("syntax error at line: %v", line.Number))
			}
			name = strings.TrimSpace(parts[0])
			if b.HasRule(name) {
				return nil, errors.New(fmt.Sprintf("duplicate rule name at line: %v", line.Number))
			}
			isPointer, logic, err := b.determineLogic(strings.TrimSpace(parts[1]), line.Number)
			if err != nil {
				return nil, err
			}
			currNode = NewBuildTree(name, logic, line.Number)
			currNode.IsPointer = isPointer
		} else if currIndent+1 < indent {
			return nil, errors.New(fmt.Sprintf("indent to large at line %v", line.Number))
		} else if indents[indent-1].count >= 2 {
			return nil, errors.New(fmt.Sprintf("more than 2 values for rule node (indent level %v), found at line %v", indent-1, line.Number))
		} else {
			// an inner rule node
			nodeBool, currNodeInner, err := b.buildLineTreeNode(line)
			if err != nil {
				return nil, err
			}
			// a comment line
			if currNodeInner == nil {
				continue
			}
			currNode = currNodeInner
			prevIndent := indents[indent-1]
			prevIndent.bt.UpdateNext(nodeBool, currNode)
			prevIndent.count += 1
		}
		// save last node for indent level
		if indent >= len(indents) {
			indents = append(indents, indentTuple{currNode, 0})
		} else {
			indents[indent] = indentTuple{currNode, 0}
		}
		currIndent = indent
	}
	if len(indents) != 0 {
		base = indents[0].bt
		// use nil as a placeholder, add rule's name to builder
		b.rules[base.Name] = nil
	}
	return
}

// buildLineTreeNode builds a rule node from a single line.
func (b *Builder) buildLineTreeNode(line Line) (nodeBool bool, bt *BuildTree, err error) {
	strLine := strings.TrimSpace(line.Value)
	// return nil for a comment
	if strLine[0] == '#' {
		return false, nil, nil
	}
	parts := strings.SplitN(strLine, " ", 2)
	if len(parts) != 2 {
		return false, nil, errors.New(fmt.Sprintf("no boolean indicator found at line: %v", line.Number))
	}
	switch parts[0] {
	case "+":
		nodeBool = true
	case "-":
		nodeBool = false
	default:
		return false, nil, errors.New(fmt.Sprintf("unknown boolean indicator at line: %v", line.Number))
	}
	logic := parts[1]
	isPointer, logic, err := b.determineLogic(logic, line.Number)
	if err != nil {
		return false, nil, err
	}
	bt = NewBuildTree(fmt.Sprintf("builtin-%v", line.Number), logic, line.Number)
	bt.IsPointer = isPointer
	return nodeBool, bt, nil
}

// determineLogic checks if logic is a rule (true) or a sentance (false)
func (b *Builder) determineLogic(logic string, lineNumber int) (bool, string, error) {
	isPointer := false
	if len(logic) > 2 {
		if logic[0] == ':' && logic[len(logic)-1] == ':' {
			logic = logic[1 : len(logic)-1]
			if _, ok := b.rules[logic]; !ok {
				return false, logic, errors.New(fmt.Sprintf("unknown rule name '%v' at line %v", logic, lineNumber))
			}
			isPointer = true
		}
	}
	return isPointer, logic, nil
}

// HasRule checks if Builder has a given rule.
func (b *Builder) HasRule(name string) bool {
	if _, ok := b.rules[name]; ok {
		return true
	}
	return false
}
