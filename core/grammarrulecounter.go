package core

// Grammar is part of both ast/ & core/ package,
// this helps preventing circular deps
type GrammarRuleCounter interface {
	IsReady() bool
	CountRules() int
}
