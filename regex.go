package main

import "regexp"

type Regex struct {
	GNode
	reStr string
	re    regexp.Regexp
}

func NewRegex(s string) Regex {
	var re Regex
	if compiledRegexp, e := regexp.Compile("(" + regexString + ")"); e != nil {
		panic("Invalid regex: " + regexString)
	} else {
		re = Regex{newGNode(), s, *compiledRegexp}
	}
	return re
}

func (re Regex) GetGNode() GNode       { return re.GNode }
func (re Regex) ContentString() string { return Magenta(re.re.String()) }
func (re Regex) HandlesChildLabel()    { return false }
func (re Regex) Labels() []string      { return re.GNode.Labels() }
func (re Regex) Captures() []astnode   { return re.GNode.Captures() }
func (re Regex) Parse(ctx *ParseContext) astnode {
	return _wrap(func(ctx, _) astnode {
		if didMatch, sMatch := ctx.code.MatchRegexp(re.re); !didMatch {
			return nil
		} else {
			return NewNativeString(sMatch)
		}
	})(re, ctx)
}
