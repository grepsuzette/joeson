package main

import . "grepsuzette/joeson/colors"
import "grepsuzette/joeson/helpers"
import "strconv"
import "strings"

// GNode class functions and fields of joeson.coffee have been put
// in gnode_exclassfields.go
/*
   In addition to the attributes defined by subclasses,
     the following attributes exist for all nodes.
   node.rule = The topmost node of a rule.
   node.rule = rule # sometimes true.
   node.name = name of the rule, if this is @rule.
*/
type GNode struct {
	Node
	id                int                // joeson.coffee:604: `node.id = @numRules++`, in Grammar.
	Name              string             // rule name if this.rule == this
	label             string             // "" by default <- i have a doubt that maybe should be *string.  Because of joeson.go:832...
	rules             map[string]astnode // see ref.go, `grammar.rules[xx].Parse()` implies it is astnode
	rule              astnode
	parent            astnode                                     // note the grammar tree should be a DAG
	grammar           Grammar                                     // see joeson.coffee:592, and joeson.coffee:530
	handlesChildLabel func() bool                                 // or nil. See Existential and Sequence
	skipCache         bool                                        // option
	skipLog           bool                                        // option
	capture           bool                                        // true by default, it's false for instance for Str
	cb                func(GNode, astnode, *ParseContext) astnode // option
	parse             func(*ParseContext) astnode                 // note: The ParseContext is modifiable by Parse(), hence the pointer
	_labels           helpers.Varcache[[]string]                  // internal cache for labels()
	_captures         helpers.Varcache[[]astnode]                 // internal cache for captures()
	_type             helpers.Varcache[string]                    // see sequence.go .Type() and parse
	_origin           Origin                                      // automatically set by prepareResult when a node is being parsed (prepareResult is called by wrap)
}

func NewGNode() GNode {
	return GNode{rules: map[string]astnode{}, capture: true}
}

func NewGNodeNamed(name string, rules map[string]astnode) GNode {
	return GNode{name: name, rules: rules, capture: true}
}

func (gn GNode) Labels() []string {
	return gn._labels.GetCacheOrSet(func() []string {
		if gn.label != "" {
			return []string{gn.label}
		} else {
			return []string{}
		}
	})
}

func (gn GNode) Captures() []astnode {
	return gn._captures.GetCacheOrSet(func() []astnode {
		if gn.capture {
			return []astnode{gn}
		} else {
			return []astnode{}
		}
	})
}

func (gn GNode) Prepare() {} // please put nothing in here

func (gn GNode) ContentString() string { return "<naked GNode, please redefine>" }
func (gn GNode) HandlesChildLabel()    { return false }
func (gn GNode) ToString() string {
	s := ""
	if gn == gn.rule {
		s = Red(gn.Name + ": ")
	} else if gn.label != "" {
		s = Cyan(gn.label + ":")
	}
	return s + " " + gn.ContentString()
}

func (gn GNode) Include(name string, rule astnode) {
	if rule.GNode.Name == "" {
		rule.GNode.Name = name
	}
	gn.rules[name] = rule
}

// find a parent in the ancestry chain that satisfies condition
func (gn GNode) FindParent(fcond func(candidateParent astnode) bool) astnode {
	var parent = gn.parent
	for {
		if fcond(parent) {
			return parent
		} else {
			parent = parent.parent
		}
	}
}
