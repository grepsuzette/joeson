package main

// TODO this will probably go to trash if 1 option
type Opts struct {
	debug bool
	// returnContext bool  // see joeson.coffee:657  this is never used and won't be implemented
}

type Option func(*Opts)

func SetDebug(b bool) Option {
	return func(opts *Opts) {
		opts.debug = b
	}
}
