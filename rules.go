package main

import "strconv"
import "strings"

// a grammar is either formed using i(), o() and rules() (as in joeson_test.coffee)
// or directly through the use of ast nodes (as in joeson.coffee:783)
type iorule struct {
	isO           bool   // if not, it's an I rule
	nameOrContent string //
	                     // name Or content (e.g. "CHOICE" Or "_ LABELED")
	                     // o() and i() accept only one string,
	                     // and this is either a simple name (/^[A-z._]+$)
	                     // or an unnamed sequence (e.g. "_ FOO", or "FOO BAR*")
	content string       // for I rule
	node  astnode        // interface
	parent *iorule
	index int
	rules []iorule   // When subrules. cb is nil in that case
	cb    ioCallback // it must be nil when contains subrules
}
type iorules []iorule

// o: any rule node, possibly part of a Rank node
// i: include line... Not included in the Rank order
// tokens: helper for declaring tokens
func rules(rules ...iorule) []iorule       { return rules }
func o(sym string, a ...[]iorule) iorule   { return iorule{true, sym, "", unpack(a), nil} }
func i(sym string, content string, cb ...ioCallback) iorule   { 
	return iorule{false, sym, content, unpack(a), nil} 
}
func ocb(sym string, cb ioCallback) iorule { return iorule{true, sym, "", nil, cb} }
func icb(sym string, cb ioCallback) iorule { return iorule{false, sym, "", nil, cb} }
func ost(sym string, seq sequenceString, cb ioCallback) iorule {
	return iorule{true, sym, "", decompose(true, seq), cb}
}
// Not sure if really used as i() already has extra `content` field
// func ist(sym string, seq sequenceString, cb ioCallback) iorule {
// 	return iorule{false, sym, "", decompose(false, seq), cb}
// }

// toRules() is only for ILine
func (rule iorule) toRules() map[string]GNode {
	if rule.isO {
		panic("can not call toRules() on OLine")
	}
	if rule.rules == nil {
		panic("ILine can not have nil rules")
	}
	return rule.rules
}

// OLine specific (toRule() is only for OLine)
// As in the original joeson.coffee,
// it will call getRule(), which will call this func again.
// This 2 func recursion is the mechanism to compile the grammar.
//
// For parameters,
//  - name: "" for none
//  - idx: -1 for non
func (line iorule) toRule(parentRule *iorule, name string, idx int) iorule {
	if !line.isO {
		panic("can not call toRule() on ILine")
	}
	rule := getArgs() // TODO
	ruleName := ""
	if name == "" {
		if regexp.MustCompile(`^[A-z_.]+$`).MatchString(line.nameOrContent) {
			ruleName = line.nameOrContent
		} else if idx > -1 && parentRule != nil {
			// unnamed (i.e. line.nameOrContent like "BAZ _ FOO+ BAR*")
			ruleName = parentRule.name + "[" + strconv.Itoa(idx) + "]"
		} else {
			panic("Name undefined for 'o' line")
		}
	}
}

type Attrs struct {
	cb        func()
	skipCache bool
	skipLog   bool
}

func (rule iorule) init(args anything...) {
	for k, arg := range args {
		if (k == "skipCache") 
		else if (k == "skipCache") 
		else if (cb == "skipCache") 
		else {

		}
	}
}
func (rule iorule) getArgs() (rule iorule, attrs Attrs) {
	(rule, rest...) = @args
}

// conveniency: just so o() and i() can be one-liners with their optional ...[]iorule
func unpack(a [][]iorule) []iorule {
	if len(a) > 0 {
		return a[0]
	} else {
		return []iorule{}
	}
}

// conveniency: just so ost() and ist() are shorter
func decompose(isO bool, s sequenceString) []iorule {
	a := []iorule{}
	var f func(sym string, a ...[]iorule) iorule
	if isO {
		f = o
	} else {
		f = i
	}
	for _, elt := range strings.Split(string(s), "|") {
		a = append(a, f(elt))
	}
	return a
}
