package core

type ParseOptions struct {
	CbBuilder func(Astnode, *ParseContext, Astnode) Astnode
	SkipLog   bool
	SkipCache bool
	Debug     bool
}
