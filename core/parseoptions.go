package core

type ParseOptions struct {
	// `O("WORD", func(it Astnode) Astnode { return ast.NewRef(it) })`
	CbBuilder func(Astnode, *ParseContext) Astnode
	SkipLog   bool
	SkipCache bool
	Debug     bool
}
