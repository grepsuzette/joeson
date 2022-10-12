package main

/*
 Result is a parse result.
 Literally it's the byproduct of `GNode.parse()`.

 As an example, the long string parameter of o below:
   o("(label:LABEL ':')? &:(DECORATED|PRIMARY)")
 Could be parsed ("compiled grammar") giving something like:
    S(E(S(L("label",R("LABEL")), St(':'))), L('&',C(R("DECORATED"),R("PRIMARY"))))

 Where S, E, L, St are short function names for Sequence, Existential, gnode.label and Str
 These would be expressed using this Result type (in that case, using the
 `labeled` field over the `scalar` one, since there are labels and captures)
*/
// TODO this is WRONG, see str.parse, which wants to return string,
//                     see sequence.parse which returns complex format...
//                     the result is astnode-specific..
//      And we will have to rework this for golang.

type ResultType int

const (
	ResultIsString ResultType = iota
	ResultIsWhatever
	ResultIsArrayOfString
	ResultIsArray
)

type Result struct {
	kind   ResultType
	origin *Origin
	// m      map[string]*Result
	scalar  []astnode          // [for Sequence.parse] when gnode.Type() is "single" or "array"
	labeled map[string]astnode // [for Sequence.parse] when gnode.Type() is "object"
	str     string             // [when kind == ResultIsString]
	astr    []string

	// hum
	// from joeson.go:547 i get the feeling
	//    Result can also be a string.
	//    so Result could be anything, right?
	//    That would be the `it` passed by callback
	//    to the init() of astnodes
}

func (r Result) toString() string {
	return "TODO result.toString"
}

func NewResultIsArrayString(a []string) *Result {
	return &Result{ResultIsArrayOfString, nil, nil, nil, "", a}
}
func NewResultIsArray(a []astnode) *Result {
	return &Result{ResultIsArray, nil, a, nil, "", nil}
}
func NewResultIsString(s string) *Result { return &Result{ResultIsString, nil, nil, nil, s, nil} }

func (r Result) ExtractString() string {
	if r.kind == ResultIsString {
		return r.str
	} else {
		panic("expecting ResultIsString")
	}
}

func (r Result) ExtractArrayOfString() []string {
	if r.kind == ResultIsArrayOfString {
		// TODO should we check for nil?
		return r.astr
	} else {
		panic("expecting ResultIsArrayOfString")
	}
}
