package helpers

import (
	"fmt"
	"golang.org/x/exp/constraints"
	"strings"
)

// escape characters that need to be and filter out non-ascii characters
func Escape(str string) string {
	return ToAscii(EscapeButKeepNonAscii(str))
}

func EscapeButKeepNonAscii(str string) string {
	replacer := strings.NewReplacer("\b", "\\b", "\\", "\\\\", "\f", "\\f", "\r", "\\r", "\n", "\\n", "\u2028", "\\u2028", "\u2029", "\\u2029", `"`, `\\\"`, `\0`, `\\0`)
	return replacer.Replace(str)
}

func Indent(c int) string {
	return strings.Join(make([]string, c+1), "  ")
}

// Characters \u0080-\uffff and \u0001b to be written in textual
// form, e.g. "\\u" + code
func ToAscii(str string) string {
	// %v	the value in a default format
	// 	when printing structs, the plus flag (%+v) adds field names
	// %#v	a Go-syntax representation of the value
	// %T	a Go-syntax representation of the type of the value
	// %%	a literal percent sign; consumes no value
	//
	// String and slice of bytes (treated equivalently with these verbs):
	//
	// %s	the uninterpreted bytes of the string or slice
	// %q	a double-quoted string safely escaped with Go syntax
	// %x	base 16, lower-case, two characters per byte
	// %X	base 16, upper-case, two characters per byte
	return fmt.Sprintf("%s", str) // TODO likely not right
}

/*
@toAscii = toAscii = (str) ->
  return str.replace /[\u001b\u0080-\uffff]/g, (ch) ->
    code = ch.charCodeAt(0).toString(16)
    code = "0" + code while code.length < 4
    "\\u"+code
*/

/*
@escape = (str, asciiOnly=yes) ->
  str = str.replace /[\\\b\f\n\r\t\x22\u2028\u2029\0]/g, (s) ->
    switch s
      when "\\" then "\\\\"
      when "\b" then "\\b"
      when "\f" then "\\f"
      when "\n" then "\\n"
      when "\r" then "\\r"
      when "\u2028" then "\\u2028"
      when "\u2029" then "\\u2029"
      when '"'  then "\\\""
      when "\0" then "\\0"
      else s
  str = toAscii str if asciiOnly
  return str
*/

func TypeOfToString(t any) string {
	return fmt.Sprintf("%T", t)
}

func PadLeft(s string, n int) string {
	l := len(s)
	if n > l {
		return s + strings.Repeat(" ", n-l)
	} else {
		return s
	}
}

func PadRight(s string, n int) string {
	l := len(s)
	if n > l {
		return strings.Repeat(" ", n-l) + s
	} else {
		return s
	}
}

func Min[T constraints.Ordered](a, b T) T {
	if a < b {
		return a
	}
	return b
}

func Max[T constraints.Ordered](a, b T) T {
	if a > b {
		return a
	}
	return b
}

// SliceString must act similar to javascript "string".slice(start, end)
func SliceString(s string, start int, end int) string {
	if start < 0 || start > len(s) || end < start {
		return ""
	}
	if end > len(s) {
		end = len(s)
	}
	return s[start:end]
}
func SliceStringFrom(s string, start int) string {
	return SliceString(s, start, len(s))
}
