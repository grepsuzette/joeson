package main

import "grepsuzette/joeson/helpers"

type sequenceRepr int

const (
	Single sequenceRepr = 0
	Array               = 1
	Object              = 2
)

type Sequence struct {
	GNode
	sequence      []astnode
	_internalType helpers.Varcache[sequenceRepr] // internal cache for internalType()
	_labels       helpers.Varcache[[]string]     // internal cache for Labels()
	_captures     helpers.Varcache[[]astnode]    // internal cache for Captures()
}

func NewSequence(a []GNode) Sequence {
	seq := Sequence{newGNode(), a}
	return seq
}

func (seq Sequence) GetGNode() GNode    { return seq.GNode }
func (seq Sequence) HandlesChildLabel() { return true }
func (seq Sequence) Prepare()           {}

func (seq Sequence) Labels() []string {
	return seq._labels.GetCacheOrSet(func() []string {
		a := []string{}
		for _, child := range seq.sequence {
			for _, label := range child.Labels() {
				a = append(a, label)
			}
		}
		return a
	})
}

func (seq Sequence) Captures() []astnode {
	return ref.GNode._captures.GetCacheOrSet(func() []astnode {
		a := []astnode{}
		for _, child := range seq.sequence {
			for _, capture := range child.Captures() {
				a = append(a, capture)
			}
		}
		return a
	})
}

func (seq Sequence) internalType() sequenceRepr {
	return seq._internalType.GetCacheOrSet(func() sequenceRepr {
		if len(seq.GNode.Labels()) == 0 {
			if len(seq.GNode.Captures()) > 1 {
				return Array
			} else {
				return Single
			}
		} else {
			return Object
		}
	})
}

func (seq Sequence) Parse(ctx *ParseContext) astnode {
	return _wrap(func(ctx, _) astnode {
		// TODO we can probably do better with the way Result works now
		switch seq.internalType() {
		case Array:
			a := make([]astnode, 0)
			for _, x := range seq.sequence {
				res := x.parse(ctx)
				if res == nil {
					return nil
				} else if x.GetGNode().capture {
					a = append(a, res)
					// ^ omg, result is not even what i though
					// is Result just some node of an ast tree?
				}
			}
			return NewNativeArray(a)
		case Single:
			var result *Result = nil
			for _, x := range seq.sequence {
				res := x.Parse(ctx)
				if res == nil {
					return nil
				} else if x.GetGNode().capture {
					result = res
				}
			}
			return result
		case Object:
			for _, child := range seq.sequence {
				var result *Result = &Result{nil, nil, make(map[string]astnode)}
				res := child.Parse(ctx)
				if res == nil {
					return nil
				}
				if child.GNode().label == "&" {
					if result == nil {
						result = res
					} else {
						for k, v := range res.labeled {
							result.labeled[k] = v
						}
					}
				} else if child.GNode().label == "@" {
					if result == nil {
						result = res
					} else {
						for k, v := range result.labeled {
							res.labeled[k] = v
						}
						result = res
					}
					// } else if child.GNode.label != nil {
				} else { // TODO check whether it must test if label != nil
					// ^ i have a doubt that maybe label should be *string.
					//   Because of joeson.go:832...
					result.labeled[child.label] = res
				}
			}
			return result
		default:
			panic("Unexpected type " + seq.Type())
		}
		panic("Error")
	})(ctx, seq)
}

func (seq Sequence) ContentString() string {
	var b strings.Builder
	for _, x := range seq.sequence {
		b.WriteString(x.contentString() + " ")
	}
	return Blue("(") + b.String() + Blue(")")
}
