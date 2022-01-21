# dbbr

A golang package to use the DBBR format.

The DBBR format allows to repressent a binary trees, in a simplistic manner.  
Each tree describes a rule, which will be exported by name.  
Each rule is made out of nodes. Each node has 2 sub-nodes: `+` for true and `-` for false, 
based on the output of the node's logic.
The indentation difference between a node and it's sub-nodes is 1 tab.

The format does not specify the logic used to chek the rules, which should be passed to 
the builder, thus allowing to use the format for different purposes.  

---

## DBBR Format Structure:

general rule structure:

```
	<rule_name>: SENTENCE
	[<\t>	<+|-> SENTENCE ]

	SENTENCE := <boolean_logic>|<:rule_name:>
```

builtin rules:
- TRUE: always returns true.
- FALSE: always returns false.

> Note: True and False are defaulted to `+` and `-` respectively, and are not needed to 
be explicitly written.

---

## Usage

In order to create a builder, which will build the rules, a logic parser must be 
implemented. The parser should implement the `dbbr.LogicParser` interface. 
basicly, it needs to implement the function `ParseLine(string) (dbbr.RuleLogic, error)`. 

Examples and demos can be found at the demo folder.
To run each demo:

```bash
	go run <demofile.go>
```

