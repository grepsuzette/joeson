package core

// helps storing grammar nodes,
// Grammar being part of ast/ & core/ must not depend on ast
type GrammarRuleCounter interface {
	IsReady() bool
	CountRules() int
}
