# The demo of using WASM in Envoy
This demo will show how to use wasm filter in Envoy on top of [proxy-wasm](https://github.com/proxy-wasm/spec/blob/master/docs/WebAssembly-in-Envoy.md).

Proxy-wasm has different language SDKs.
- AssemblyScript SDK
- C++ SDK
- Go (TinyGo) SDK
- Rust SDK
- Zig SDK

This demo will focus on the TinyGo and Rust.

The demo itself assumes you are running in an Ubuntu system.
The demo itself is verified in a ubuntu 20.04 in WSL2. 

## TinyGo

Here is an [reference](https://github.com/tetratelabs/proxy-wasm-go-sdk/tree/main/examples) for how to using TinyGo SDK

We need to first install a golang with version at lease v1.17
```
brew install golang
```

Install TinyGo then
```
wget https://github.com/tinygo-org/tinygo/releases/download/v0.23.0/tinygo_0.23.0_amd64.deb
sudo dpkg -i tinygo_0.23.0_amd64.deb
export PATH=$PATH:/usr/local/bin
```

Build the demo application
```
cd wasm-tinygo

# Download dependencies
go mod tidy

# Build
tinygo build -o plugin.wasm -scheduler=none -target=wasi main.go
```
We should see a binary file called `plugin.wasm` in your current directory.

Let's deploy this binary into envoy
```
# make sure you will be under WASM-EnvoyFilter directory after executing this
cd ..

# build a envoy image with specific version
docker build . -t envoy:wasm

# start an Envoy in container
docker run --rm -it \
           -v $(pwd)/local-envoy-config.yaml:/tmp/local-envoy-config.yaml \
           -v $(pwd)/wasm-tinygo/plugin.wasm:/tmp/plugin.wasm \
           -p 8000:8000 -t envoy:wasm -c /tmp/local-envoy-config.yaml
```

Now our wasm plugin is running with Envoy.
Open another terminal session
```
curl 127.0.0.1:8000
```

You will find new logs are dumping out inside the previous terminal session that Envoy runs.
In our demo TinyGo wasm plugin, we make it to print out the request header info inside wasm console.
On my side, the logs are like
```
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> :authority: 127.0.0.1:8000
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> :path: /
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> :method: GET
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> :scheme: http
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> user-agent: curl/7.68.0
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> accept: */*
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> x-forwarded-proto: http
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> x-request-id: da9ead29-aa2b-46a3-a609-bd0d112813d6
[2022-07-20 21:00:08.484][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: request header --> test: best
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- :status: 200
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- date: Wed, 20 Jul 2022 21:00:08 GMT
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- content-type: application/json; charset=utf-8
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- content-length: 365
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- connection: keep-alive
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- etag: W/"16d-KDQNNQhCImcUPgmcI2NDKSVPPdg"
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- vary: Accept-Encoding
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- set-cookie: sails.sid=s%3AdI0LH8cpgLaojhSX2yEVqCvgtyvYsiRT.5dcuIYQI0zc2CbARMRj3oMMRNlsQuJbS7NNz%2FfkfV2Q; Path=/; HttpOnly
[2022-07-20 21:00:08.830][42][info][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1218] wasm log: response header <-- x-envoy-upstream-service-time: 344
```

### TinyGo with connecting to Redis
This section is just demonstrate that TinyGo is a trimmed version of Golang, that compiler doesn't not support net package.

```
cd RedisTinyGo

tinygo build -o plugin.wasm -scheduler=none -target=wasi main.go
```
You should see the errors like
```
../../../../../../go/pkg/mod/github.com/go-redis/redis@v6.15.9+incompatible/options.go:114:22: netDialer.Dial undefined (type *net.Dialer has no field or method Dial)
../../../../../../go/pkg/mod/github.com/go-redis/redis@v6.15.9+incompatible/sentinel.go:235:13: DialTimeout not declared by package net
```
This is because go-redis uses net package in the bottom layer, which is not supported in TinyGo.

Writing a Redis client from scratch using [TinyNet](https://github.com/alphahorizonio/tinynet) is not feasible either due to [unisockets](https://github.com/alphahorizonio/unisockets) is not compatible to wasi based TinyGo compiler. 

## Rust
Install Rust
```
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# Install WASM dependencies
rustup toolchain install nightly
rustup target add wasm32-unknown-unknown

# For linux you probably need to do this as well
sudo apt-get update
sudo apt install build-essential
```

Build WASM plugin
```
cd ../wasm-rust-filter
cargo build --target=wasm32-unknown-unknown --release
```

There will be a wasm binary generated under `wasm-rust-filter/target/wasm32-unknown-unknown/release/plugin.wasm`

Let's deploy it in the Envoy
```
# make sure you will be under WASM-EnvoyFilter directory after executing this
cd ..

docker run --rm -it \
           -v $(pwd)/local-envoy-config.yaml:/tmp/local-envoy-config.yaml \
           -v $(pwd)/wasm-rust-filter/target/wasm32-unknown-unknown/release/plugin.wasm:/tmp/plugin.wasm \
           -p 8000:8000 -t envoy:wasm -c /tmp/local-envoy-config.yaml
```