package joeson

const (
	Esc   string = "\x1b"
	Reset string = Esc + "[0m"
)

func Black(s string) string       { return Esc + "[30m" + s + Reset }
func Red(s string) string         { return Esc + "[31m" + s + Reset }
func Green(s string) string       { return Esc + "[32m" + s + Reset }
func Yellow(s string) string      { return Esc + "[33m" + s + Reset }
func Blue(s string) string        { return Esc + "[34m" + s + Reset }
func Magenta(s string) string     { return Esc + "[35m" + s + Reset }
func Cyan(s string) string        { return Esc + "[36m" + s + Reset }
func White(s string) string       { return Esc + "[37m" + s + Reset }
func BoldBlack(s string) string   { return Esc + "[1;30m" + s + Reset }
func BoldRed(s string) string     { return Esc + "[1;31m" + s + Reset }
func BoldGreen(s string) string   { return Esc + "[1;32m" + s + Reset }
func BoldYellow(s string) string  { return Esc + "[1;33m" + s + Reset }
func BoldBlue(s string) string    { return Esc + "[1;34m" + s + Reset }
func BoldMagenta(s string) string { return Esc + "[1;35m" + s + Reset }
func BoldCyan(s string) string    { return Esc + "[1;36m" + s + Reset }
func BoldWhite(s string) string   { return Esc + "[1;37m" + s + Reset }

func Bold(s string) string { return Esc + "[1m" + s + Reset }
