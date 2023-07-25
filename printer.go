package joeson

import (
	"strconv"
	"strings"
)

// [experimental] might go away soon!
//
// In most cases, ast.String() should be used instead.
//
// Unlike Ast String(), Printer also prints the NAMES of the rules that
// produced the object being printed.
//
// Printer is intended to help quickly figuring out what is produced by
// some high level rules, e.g. to write mappers.
//
// For instance, `fmt.Println(ast.(Printer).Text())
type Printer interface {
	Text() string
}

// Only Native objects implement this
var (
	_ Printer = NativeArray{}
	_ Printer = NativeInt{}
	_ Printer = NativeMap{}
	_ Printer = NativeString{}
	_ Printer = NativeUndefined{}
)

func (na NativeArray) Text() string {
	var b strings.Builder
	b.WriteString("[")
	first := true
	for i, v := range na.Array {
		if !first {
			b.WriteString(", ")
		}
		b.WriteString("(" + strconv.Itoa(i) + ") ")
		if printer, ok := v.(Printer); ok {
			b.WriteString(printer.Text())
		} else {
			b.WriteString(v.String())
		}
		first = false
	}
	b.WriteString("]")
	return printRuleName(&na) + b.String()
}

func (ni NativeInt) Text() string {
	return printRuleName(ni) + strconv.Itoa(ni.int)
}

func (nm NativeMap) Text() string {
	var b strings.Builder
	b.WriteString("{")
	first := true
	for k, v := range nm.Map {
		if !first {
			b.WriteString(", ")
		}
		b.WriteString(k)
		b.WriteString(":")
		if printer, ok := v.(Printer); ok {
			b.WriteString(printer.Text())
		} else {
			b.WriteString(v.String())
		}
		first = false
	}
	b.WriteString("}")
	return printRuleName(nm) + b.String()
}

func (ns NativeString) Text() string {
	var b strings.Builder
	b.WriteString(`"`)
	b.WriteString(ns.Str)
	b.WriteString(`"`)
	return printRuleName(ns) + b.String()
}

func (nu NativeUndefined) Text() string {
	return printRuleName(nu) + "⊘"
}

func printRuleName(ast Ast) string {
	if ast.HasAttribute("RuleName") {
		return Green(ast.GetAttribute("RuleName").(string) + "=")
	} else {
		return ""
	}
}
