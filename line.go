package main

import "strings"

// In joeson.coffee, those were Line, ILine and OLine.
// iorule is the recursive element of a tree
// expressing rules of a grammar in code.
//
// As grammar is parsed, iorule are transformed into cRules
// (computed rules).
//
// var rulez []iorule = Rules(
//	  o("EXPR", Rules(
//		o("CHOICE _"),
//		i("INT", "/[0-9]+/", func(it ctx) astnode { return Number{43} })
//	  ))
// )
//
// ^- each o(), i(), ost(), ist(), ocb() or icb() call
//    simply creates an iorule struct.
type iorule struct {
	isO   bool // if not, it's an I rule
	name  string
	rules []iorule   // When subrules. cb is nil in that case
	cb    ioCallback // it must be nil when contains subrules
}
type iorules []iorule

// o: any rule node, possibly part of a Rank node
// i: include line... Not included in the Rank order
// tokens: helper for declaring tokens
func rules(rules ...iorule) []iorule       { return rules }
func o(sym string, a ...[]iorule) iorule   { return iorule{true, sym, unpack(a), nil} }
func i(sym string, a ...[]iorule) iorule   { return iorule{false, sym, unpack(a), nil} }
func ocb(sym string, cb ioCallback) iorule { return iorule{true, sym, nil, cb} }
func icb(sym string, cb ioCallback) iorule { return iorule{false, sym, nil, cb} }
func ost(sym string, seq sequenceString, cb ioCallback) iorule {
	return iorule{true, sym, decompose(true, seq), cb}
}
func ist(sym string, seq sequenceString, cb ioCallback) iorule {
	return iorule{false, sym, decompose(false, seq), cb}
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

type Rule struct {
	name string
	rule Rule
}

// Now we are porting Line, a big technical piece
type Line struct {
}

// name:       The final and correct name for this rule
// rule:       A rule-like object
// parentRule: The actual parent Rule instance
// attrs:      {cb,...}, extends the result
func (line Line) getRule(
	name string,
	rule Rule,
	parentRule Rule,
	attrs Unknown,
) Rule {
	// in coffee, a rule could be a string, array or oline
}
