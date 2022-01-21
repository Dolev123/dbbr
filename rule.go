package dbbr

var (
	// True is alwyas true.
	True Rule = &rule{
		name:   "True",
		nTrue:  nil,
		nFalse: nil,
		logic:  func(_ interface{}) bool { return true },
	}
	//False is always false.
	False Rule = &rule{
		name:   "False",
		nTrue:  nil,
		nFalse: nil,
		logic:  func(_ interface{}) bool { return false },
	}
)

// RuleLogic repressents the function which detemines the logic of a rule.
type RuleLogic func(interface{}) bool

// LogicParser is an interface for parsing the logic of a rule when building.
type LogicParser interface {
	ParseLine(string) (RuleLogic, error)
}

// Rule is a binary decision tree with a logic functions.
type Rule interface {
	Name() string
	// Check runs the input against the full rule tree logic
	Check(interface{}) bool
	// CheckLogic runs the input against only the current rule
	CheckLogic(interface{}) bool
	// NextTrue returns the next rule node if true
	NextTrue() Rule
	// NextFalse returns the next rule node if false
	NextFalse() Rule
	// Next returns the next rule by boolean value
	Next(bool) Rule
}

type rule struct {
	name   string
	nTrue  Rule
	nFalse Rule
	logic  RuleLogic
}

// NewRule creates a new Rule node.
func NewRule(name string, logic RuleLogic, nTrue Rule, nFalse Rule) Rule {
	return &rule{
		name:   name,
		logic:  logic,
		nTrue:  nTrue,
		nFalse: nFalse,
	}
}

func (r *rule) Name() string {
	return r.name
}

func (r *rule) Check(i interface{}) bool {
	res := r.logic(i)
	if r.Next(res) != nil {
		return r.Next(res).Check(i)
	}
	return res
}

func (r *rule) CheckLogic(i interface{}) bool {
	return r.logic(i)
}

func (r *rule) NextTrue() Rule {
	return r.nTrue
}

func (r *rule) NextFalse() Rule {
	return r.nFalse
}

func (r *rule) Next(b bool) Rule {
	if b {
		return r.NextTrue()
	}
	return r.NextFalse()
}
