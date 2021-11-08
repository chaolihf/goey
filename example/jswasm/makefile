EXAMPLE ?= align

serve: main.wasm wasm_exec.js index.html
	python3 -m http.server 8080
	-$(RM) wasm_exec.js

main.wasm: $(wildcard ../../*.go) $(wildcard ../../windows/*.go)
	GOOS=js GOARCH=wasm go build -o main.wasm ../${EXAMPLE}

wasm_exec.js:
	cp `go env GOROOT`/misc/wasm/wasm_exec.js .

clean:
	-$(RM) *.wasm *.js
