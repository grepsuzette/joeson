package helpers

import (
	"fmt"
	"strings"
	"unicode/utf8"

	"golang.org/x/exp/constraints"
)

// escape characters that need to be and filter out non-ascii characters
func Escape(str string) string {
	return ToAscii(EscapeButKeepNonAscii(str))
}

var escaper = strings.NewReplacer(
	"\n", `\n`,
	"\r", `\r`,
	"\t", `\t`,
	`\`, `\\`,
	"\a", `\a`,
	"\f", `\f`,
	"\b", `\b`,
	"\v", `\v`,
)

func EscapeButKeepNonAscii(str string) string {
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
	// replacer := strings.NewReplacer( /*"\b", "\\b", "\f", "\\f", "\r", "\\r",*/ "\n", "n", "\u2028", "\\u2028", "\u2029", "\\u2029" /*`"`, `\\\"`,*/, `\0`, `\\0`)
	// return replacer.Replace(str)
	return escaper.Replace(str)
	return str
}

func Indent(c int) string {
	return strings.Join(make([]string, c+1), "  ")
}

// Characters \u0080-\uffff and \u0001b to be written in text form, e.g. "\\u" + code
func ToAscii(str string) string {
	return str // TODO
}

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

// [Deprecated] SliceString must act similar to javascript "string".slice(start, end)
// indexStart: The index of the first character to include in the returned substring.
// indexEnd: The index of the first character to exclude from the returned substring.
// slice() extracts up to but not including indexEnd. For example, str.slice(1, 4)
// extracts the second character through the fourth character (characters indexed 1, 2, and 3).
func SliceString(s string, indexStart int, indexEnd int) string {
	// If indexStart >= str.length, an empty string is returned.
	if indexStart >= len(s) {
		return ""
	}
	// If indexStart < 0, the index is counted from the end of the string. More
	// formally, in this case, the substring starts at max(indexStart + str.length, 0).
	if indexStart < 0 {
		indexStart = Max(indexStart+len(s), 0)
	}
	// If indexEnd < 0, the index is counted from the end of the string. More
	// formally, in this case, the substring ends at max(indexEnd + str.length, 0).
	if indexEnd < 0 {
		indexEnd = Max(indexEnd+len(s), 0)
	}
	// If indexEnd <= indexStart after normalizing negative values (i.e.
	// indexEnd represents a character that's before indexStart), an empty
	// string is returned.
	if indexEnd <= indexStart {
		return ""
	}
	if indexEnd > len(s) {
		indexEnd = len(s)
	}
	return s[indexStart:indexEnd]
}

func BoolToString(b bool) string {
	if b {
		return "y"
	} else {
		return "n"
	}
}

// take last `nTake` runes from `abc`, in natural order. ("abc", 2) -> "bc"
// When nTake is more than the number of available runes, return all of them.
// If RuneError at some point, panic.
func LastNRunes(abc string, nTake int) string {
	totalSize := 0 // size in bytes
	slice := abc
	for nTake > 0 {
		if len(slice) == 0 {
			break
		}
		c, csize := utf8.DecodeLastRuneInString(slice)
		if c == utf8.RuneError {
			panic(fmt.Sprintf("RuneError: %q\n", c))
		}
		totalSize += csize // csize is char len in bytes for rune c
		slice = slice[:len(slice)-csize]
		nTake--
	}
	return abc[len(abc)-totalSize:]
}

// take last `nTake` runes from `abc` in reverse order. ("abc", 2) -> "cb"
func LastNRunesReversed(abc string, nTake int) string {
	if nTake < 0 {
		return ""
	}
	_, sizeinit := utf8.DecodeLastRuneInString(abc)
	slice := abc
	var b strings.Builder
	for i := len(abc) - sizeinit; i >= len(abc)-sizeinit-nTake+1; i-- {
		r, size := utf8.DecodeLastRuneInString(slice)
		if r == utf8.RuneError {
			break
		}
		b.WriteRune(r)
		slice = slice[:len(slice)-size]
	}
	return b.String()
}
