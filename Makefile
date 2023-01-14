.PHONY: grammar
it: 
	go build . 

run: it
	./joeson

test: 
	go test ./...

