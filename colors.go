package joeson

const esc string = "\x1b"
const reset string = esc + "[0m"

func black(s string) string       { return esc + "[30m" + s + reset }
func red(s string) string         { return esc + "[31m" + s + reset }
func green(s string) string       { return esc + "[32m" + s + reset }
func yellow(s string) string      { return esc + "[33m" + s + reset }
func blue(s string) string        { return esc + "[34m" + s + reset }
func magenta(s string) string     { return esc + "[35m" + s + reset }
func cyan(s string) string        { return esc + "[36m" + s + reset }
func white(s string) string       { return esc + "[37m" + s + reset }
func boldBlack(s string) string   { return esc + "[1;30m" + s + reset }
func boldRed(s string) string     { return esc + "[1;31m" + s + reset }
func boldGreen(s string) string   { return esc + "[1;32m" + s + reset }
func boldYellow(s string) string  { return esc + "[1;33m" + s + reset }
func boldBlue(s string) string    { return esc + "[1;34m" + s + reset }
func boldMagenta(s string) string { return esc + "[1;35m" + s + reset }
func boldCyan(s string) string    { return esc + "[1;36m" + s + reset }
func boldWhite(s string) string   { return esc + "[1;37m" + s + reset }

func bold(s string) string { return esc + "[1m" + s + reset }
