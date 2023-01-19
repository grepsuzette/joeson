package colors

// in their own package so these colors can be imported with a dot
//  not importing anything else

const esc string = ""
const reset string = esc + "[0m"

func Black(s string) string       { return esc + "[30m" + s + reset }
func Red(s string) string         { return esc + "[31m" + s + reset }
func Green(s string) string       { return esc + "[32m" + s + reset }
func Yellow(s string) string      { return esc + "[33m" + s + reset }
func Blue(s string) string        { return esc + "[34m" + s + reset }
func Magenta(s string) string     { return esc + "[35m" + s + reset }
func Cyan(s string) string        { return esc + "[36m" + s + reset }
func White(s string) string       { return esc + "[37m" + s + reset }
func BoldBlack(s string) string   { return esc + "[1;30m" + s + reset }
func BoldRed(s string) string     { return esc + "[1;31m" + s + reset }
func BoldGreen(s string) string   { return esc + "[1;32m" + s + reset }
func BoldYellow(s string) string  { return esc + "[1;33m" + s + reset }
func BoldBlue(s string) string    { return esc + "[1;34m" + s + reset }
func BoldMagenta(s string) string { return esc + "[1;35m" + s + reset }
func BoldCyan(s string) string    { return esc + "[1;36m" + s + reset }
func BoldWhite(s string) string   { return esc + "[1;37m" + s + reset }

func Bold(s string) string { return esc + "[1m" + s + reset }
