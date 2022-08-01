```
cd WASM-EnvoyFilter-SharedData

go mod tidy

tinygo build -o worker.wasm -scheduler=none -target=wasi worker/main.go

tinygo build -o singleton.wasm -scheduler=none -target=wasi singleton/main.go

docker run --rm -it \
           -v $(pwd)/local-envoy-config.yaml:/tmp/local-envoy-config.yaml \
           -v $(pwd)/worker.wasm:/tmp/worker.wasm \
           -p 8000:8000 -t envoy:wasm -c /tmp/local-envoy-config.yaml

curl 127.0.0.1:8000
```