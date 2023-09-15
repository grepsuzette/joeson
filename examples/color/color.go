package main

import (
	"fmt"

	j "github.com/grepsuzette/joeson"
)

// named() creates a rule with a name
func named(name string, v any) j.NamedRule { return j.Named(name, v) }

// to make i and o rules
func i(a ...any) j.ILine { return j.I(a...) }
func o(a ...any) j.OLine { return j.O(a...) }

type Color struct {
	*j.Attr
	r, g, b int
}

func (c Color) String() string {
	return fmt.Sprintf("Color is rgb ( %d, %d, %d )", c.r, c.g, c.b)
}

func hslToRgb(h, s, l int) []int {
	// ask chatgpt for implementation
	return []int{12, 23, 34}
}

func toArrayInt(it j.Ast) []int {
	r := []int{}
	for _, v := range it.(*j.NativeArray).Array() {
		r = append(r, j.NativeIntFrom(v).(j.NativeInt).Int())
	}
	return r
}

func main() {
	gm := j.GrammarFromLines("color example", []j.Line{
		o(named("Color", []j.Line{
			o("'red'", func(it j.Ast) j.Ast { return Color{j.NewAttr(), 255, 0, 0} }),
			o("'green'", func(it j.Ast) j.Ast { return Color{j.NewAttr(), 0, 255, 0} }),
			o("'blue'", func(it j.Ast) j.Ast { return Color{j.NewAttr(), 0, 0, 255} }),
			o(named("Rgb", `'rgb' _ '(' _ Integer*_COMMA{3,3} _ ')'`), func(it j.Ast) j.Ast {
				a := toArrayInt(it)
				return Color{j.NewAttr(), a[0], a[1], a[2]}
			}),
			o(named("Hsl", `'hsl' _ '(' HslTrio ')'`), func(it j.Ast) j.Ast {
				a := toArrayInt(it)
				b := hslToRgb(a[0], a[1], a[2])
				return Color{j.NewAttr(), b[0], b[1], b[2]}
			}),
			i(named("HslTrio", `_ Integer _ ',' _ Integer '%' _ ',' _ Integer '%' _`)),
		})),
		i(named("Integer", `[1-9][0-9]*`), j.NativeIntFrom),
		i(named("_COMMA", `[ ,]*`)),
		i(named("_", `[ \t]*`)),
	})
	fmt.Println(gm.ParseString("blue"))
	fmt.Println(gm.ParseString("red"))
	fmt.Println(gm.ParseString("green"))
	fmt.Println(gm.ParseString("rgb(127, 49, 255)"))
	fmt.Println(gm.ParseString("hsl(127, 49%, 82%)"))
}
