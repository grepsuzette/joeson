package videogamedb

// this videogame database example demonstrates how to create an Ast type
// yourself.
//
// - embed j.Attr
// - define String() string

import (
	"fmt"
	"testing"

	j "github.com/grepsuzette/joeson"
)

// named() creates a rule with a name
func named(name string, v any) j.NamedRule { return j.Named(name, v) }

// to make i and o rules
func i(a ...any) j.ILine { return j.I(a...) }
func o(a ...any) j.OLine { return j.O(a...) }

type VideoGame struct {
	*j.Attr   // help implementing j.Ast easily
	id        int
	title     string
	year      int
	developer string
	score     int // tentative score by chatgpt
}

var db = []VideoGame{
	{j.NewAttr(), 1, "Civilization", 1991, "Sid Meier", 95},
	{j.NewAttr(), 3, "Rogue", 1980, "Michael Toy, Glenn Wichman, Ken Arnold", 80},
	{j.NewAttr(), 4, "Doom", 1993, "John Carmack", 90},
	{j.NewAttr(), 5, "Tetris", 1984, "Alexey Pajitnov", 90},
}
var notfound = VideoGame{j.NewAttr(), 0, "Not found", -1, "", 0}

func (v VideoGame) String() string {
	return fmt.Sprintf(
		`#%d: "%s" by %s in %d, %d/100\n`,
		v.id,
		v.title,
		v.developer,
		v.year,
		v.score,
	)
}

func findVideoGameById(it j.Ast) j.Ast {
	id := j.NewNativeIntFrom(it).Int()
	for _, v := range db {
		if v.id == id {
			return v
		}
	}
	return notfound
}

func findVideoGameByTitle(it j.Ast) j.Ast {
	title := j.NativeStringFrom(it).String()
	for _, v := range db {
		if v.title == title {
			return v
		}
	}
	return notfound
}

func findBestVideoGameOfYear(it j.Ast) j.Ast {
	// whichever was listed first in db is probably good enough
	year := j.NewNativeIntFrom(it).Int()
	for _, v := range db {
		if v.year == year {
			return v
		}
	}
	return notfound
}

func makeVideoGame(it j.Ast) j.Ast {
	// this callback maps id or title to a videogame entry
	m := it.(*j.NativeMap)
	if id, exists := m.GetExists("id"); exists {
		return findVideoGameById(id)
	} else if title, exists := m.GetExists("title"); exists {
		return findVideoGameByTitle(title)
	}
	panic("assert")
}

func TestVideogameDb(t *testing.T) {
	gm := j.GrammarFromLines("videogamedb", []j.Line{
		o(named("VideoGame", []j.Line{
			o(named("VideoGameId", `[1-9][0-9]*`), findVideoGameById),
			o(named("VideoGameTitle", `'"' [^"]* '"'`), findVideoGameByTitle),
			o(named("VideoGameBestOfYear", `'bestIn:' _ Year`), findBestVideoGameOfYear),
			i(named("Year", `a:('19'|'20'|'21') b:[0-9] c:[0-9]`), func(it j.Ast) j.Ast {
				return it.(*j.NativeMap).Concat()
			}),
		})),
		i(named("_", `[ \t]*`)),
	})
	if gm.ParseString("1").(VideoGame).title != "Civilization" {
		t.Fail()
	}
	if gm.ParseString("73264983242").(VideoGame).title != notfound.title {
		t.Fail()
	}
	if gm.ParseString(`"Doom"`).(VideoGame).developer != "John Carmack" {
		t.Fail()
	}
	if gm.ParseString("bestIn: 1984").(VideoGame).title != "Tetris" {
		t.Fail()
	}
}
