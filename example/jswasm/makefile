EXAMPLE := align

serve: main.wasm
	cp `go env GOROOT`/misc/wasm/wasm_exec.js .
	python3 -m http.server 8080
	-$(RM) wasm_exec.js

main.wasm: $(wildcard ../../*.go)
	GOOS=js GOARCH=wasm go build -o main.wasm ../${EXAMPLE}

clean:
	-$(RM) main.wasm *.js