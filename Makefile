.PHONY: grammar
it: 
	go build . 

run: it
	./Joeson

test: 
	go test .

# test: it
# 	go test ./...
