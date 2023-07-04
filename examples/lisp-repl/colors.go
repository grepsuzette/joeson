package main

const (
	esc   string = "\x1b"
	reset string = esc + "[0m"
)

func black(s string) string        { return esc + "[30m" + s + reset }
func red(s string) string          { return esc + "[31m" + s + reset }
func green(s string) string        { return esc + "[32m" + s + reset }
func yellow(s string) string       { return esc + "[33m" + s + reset }
func blue(s string) string         { return esc + "[34m" + s + reset }
func magenta(s string) string      { return esc + "[35m" + s + reset }
func cyan(s string) string         { return esc + "[36m" + s + reset }
func white(s string) string        { return esc + "[36m" + s + reset }
func bold_black(s string) string   { return esc + "[1;30m" + s + reset }
func bold_red(s string) string     { return esc + "[1;31m" + s + reset }
func bold_green(s string) string   { return esc + "[1;32m" + s + reset }
func bold_yellow(s string) string  { return esc + "[1;33m" + s + reset }
func bold_blue(s string) string    { return esc + "[1;34m" + s + reset }
func bold_magenta(s string) string { return esc + "[1;35m" + s + reset }
func bold_cyan(s string) string    { return esc + "[1;36m" + s + reset }
func bold_white(s string) string   { return esc + "[1;36m" + s + reset }

func colorKeyword(s string) string { return (s) }
func colorParen(s string) string   { return bold_cyan(s) }
func colorComma(s string) string   { return bold_green(s) }
func quoted(s string) string       { return red(`"`) + blue(s) + red(`"`) }
