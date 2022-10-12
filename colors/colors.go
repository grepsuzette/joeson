package colors

const ESC string = ""
const RESET string = ESC + "[0m"

func Black(s string) string       { return ESC + "[30m" + s + RESET }
func Red(s string) string         { return ESC + "[31m" + s + RESET }
func Green(s string) string       { return ESC + "[32m" + s + RESET }
func Yellow(s string) string      { return ESC + "[33m" + s + RESET }
func Blue(s string) string        { return ESC + "[34m" + s + RESET }
func Magenta(s string) string     { return ESC + "[35m" + s + RESET }
func Cyan(s string) string        { return ESC + "[36m" + s + RESET }
func White(s string) string       { return ESC + "[37m" + s + RESET }
func BoldBlack(s string) string   { return ESC + "[1;30m" + s + RESET }
func BoldRed(s string) string     { return ESC + "[1;31m" + s + RESET }
func BoldGreen(s string) string   { return ESC + "[1;32m" + s + RESET }
func BoldYellow(s string) string  { return ESC + "[1;33m" + s + RESET }
func BoldBlue(s string) string    { return ESC + "[1;34m" + s + RESET }
func BoldMagenta(s string) string { return ESC + "[1;35m" + s + RESET }
func BoldCyan(s string) string    { return ESC + "[1;36m" + s + RESET }
func BoldWhite(s string) string   { return ESC + "[1;37m" + s + RESET }
