.PHONY: grammar
it: 
	go build . 

run: it
	./Joeson

test: 
	go test . --run TestDebugLabel -v

# test: it
# 	go test ./...
