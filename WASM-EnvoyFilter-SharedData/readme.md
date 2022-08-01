```
tinygo build -o worker.wasm -scheduler=none -target=wasi worker/main.go

tinygo build -o singleton.wasm -scheduler=none -target=wasi singleton/main.go
```