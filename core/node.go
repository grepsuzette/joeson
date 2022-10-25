package core

// An element of a walk-able tree is a Node.
type Node struct {
	id int
}

// ignore @defineChildren
// ignore? withChildren:  HUM walk uses it
// ignore? walk
//   only in Grammar.init
// ignore validateType
// ignore serialize
// ignore validate

// So, if we want to not implment withChildren,
//  we will need to seriously consider how to
//  reimplement Grammar.init
