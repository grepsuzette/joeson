package main

var rules_tokens = rules(
	o(named("token", "quoted"), x("token")),
	i(named("quoted", "'Q' [a-zA-Z]* 'Q'"), x("quoted")),
)
