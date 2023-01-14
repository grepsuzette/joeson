package core

type ParseOptions struct {
	CbBuilder func(Ast, *ParseContext, Ast) Ast
	SkipLog   bool
	SkipCache bool
	Debug     bool
}
