wasm:
	GOARCH=wasm GOOS=js go build -o web/app.wasm

build:
	make wasm
	go build