.PHONY: grammar
it: 
	go build . 

run: it
	./Joeson

# test: it
# 	go test ./...
