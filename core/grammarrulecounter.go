package core

// Grammar is part of both ast/ & core/ package,
// this is to help preventing circular deps
type GrammarRuleCounter interface {
	IsReady() bool
	CountRules() int
	Options() TraceOptions
}
