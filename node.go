package main

import . "grepsuzette/joeson/colors"
import "grepsuzette/joeson/helpers"

type Node struct {
	id int
	// children map[string]*string // key: desc?
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
