package core

// poor man js's .bind()/.call() in go for the parsefunc

type BindableParseFunc struct {
	Bounded Astnode // nil by default
	F       func(ctx *ParseContext, x Astnode) Astnode
}

func NewBindableParseFunc(f func(*ParseContext, Astnode) Astnode) *BindableParseFunc {
	return &BindableParseFunc{nil, f}
}

func (bpf *BindableParseFunc) Bind(x Astnode) *BindableParseFunc {
	bpf.Bounded = x
	return bpf
}

func (bpf *BindableParseFunc) Call(ctx *ParseContext) Astnode {
	return bfp.F(ctx, bpf.Bounded)
}

func (bpf *BindableParseFunc) CallWith(ctx *ParseContext, x Astnode) Astnode {
	return bpf.F(ctx, x)
}
