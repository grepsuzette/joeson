package main

import (
	j "github.com/grepsuzette/joeson"
)

type attr j.Origin

func (attr attr) GetLine() int { return attr.Line } func (attr attr) SetLine(n int) {
	attr.Line = n
	attr.Start = 0
	attr.End = 0
}
func (attr attr) GetOrigin() j.Origin { return j.Origin(attr) }
func (attr attr) SetOrigin(o j.Origin) {
	attr.Code = o.Code
	attr.Line = o.Line
	attr.Start = o.Start
	attr.End = o.End
	attr.RuleName = o.RuleName
}
func (attr attr) HasAttribute(key interface{}) bool        { return false }
func (attr attr) GetAttribute(key interface{}) interface{} { panic("not implemented") }
func (attr attr) SetAttribute(key interface{}, value interface{}) {
	panic("not implemented, SetAttribute " + key.(string))
}
