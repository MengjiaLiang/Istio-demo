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

Open another terminal session
```
curl 127.0.0.1:8000/hello
```

You should see the output as
```
Hello, World!
```
This is what we configured inside our Rust code.

### Rust with connecting to Redis
This section is just demonstrate that proxy-wasm does not support the TCP connection package even the compiler can compile the plugin well.

```
# Kill the running envoy first
cd ../wasm-rust-filter-redis

cargo build --target=wasm32-unknown-unknown --release
```
Ignore the warning message as the code itself is just for demoing purpose

In this version of Rust code, we add the redis client to let the plugin try to connect with the redis.

Let's start a redis in container
```
docker run --rm --name redis -p 6379:6379 -d redis
```

Then deploy the new Rust plugin in Envoy
```
# make sure you will be under WASM-EnvoyFilter directory after executing this
cd ..

docker run --rm -it \
           -v $(pwd)/local-envoy-config.yaml:/tmp/local-envoy-config.yaml \
           -v $(pwd)/wasm-rust-filter-redis/target/wasm32-unknown-unknown/release/plugin.wasm:/tmp/plugin.wasm \
           -p 8000:8000 -t envoy:wasm -c /tmp/local-envoy-config.yaml
```

Open another terminal session and execute
```
curl 127.0.0.1:8000/hello
```
In your envoy session, you will see the error logs dumping out like
```
[2022-07-20 22:00:37.928][36][critical][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1227] wasm log: panicked at 'failed to connect to Redis: operation not supported on this platform', src/lib.rs:74:10
[2022-07-20 22:00:37.928][36][error][wasm] [external/envoy/source/extensions/common/wasm/wasm_vm.cc:39] Function: proxy_on_request_headers failed: Uncaught RuntimeError: unreachable
Proxy-Wasm plugin in-VM backtrace:
  0:  0x3015a - __rust_start_panic
  1:  0x300d6 - rust_panic
  2:  0x300a6 - _ZN3std9panicking20rust_panic_with_hook17h1c368a27f9b0afe1E
  3:  0x2f69a - _ZN3std9panicking19begin_panic_handler28_$u7b$$u7b$closure$u7d$$u7d$17h8e1f8b682ca33009E
  4:  0x2f5d9 - _ZN3std10sys_common9backtrace26__rust_end_short_backtrace17h7f7da41799766719E
  5:  0x2fd18 - rust_begin_unwind
  6:  0x317ea - _ZN4core9panicking9panic_fmt17hcdb13a4b2416cf82E
  7:  0x338f6 - _ZN4core6result13unwrap_failed17he825aa6f43b16604E
  8:  0x451f - _ZN4core6result19Result$LT$T$C$E$GT$6expect17hd84817813fccb93fE
  9:  0x47da - _ZN6plugin7connect17h80f53ab29bcaeccbE
```
The reason behind that is proxy-wasm doesn't support the underlying packages in Rust redis.

Similarly, writing a redis client from scratch by using Rust TCP connection support is not feasible either
```
# Kill the envoy
cd ../wasm-rust-filter-redis-tcp

cargo build --target=wasm32-unknown-unknown --release

cd ..

docker run --rm -it \
           -v $(pwd)/local-envoy-config.yaml:/tmp/local-envoy-config.yaml \
           -v $(pwd)/wasm-rust-filter-redis-tcp/target/wasm32-unknown-unknown/release/plugin.wasm:/tmp/plugin.wasm \
           -p 8000:8000 -t envoy:wasm -c /tmp/local-envoy-config.yaml
```

You will see a similar error when executing `curl 127.0.0.1:8000/hello`
```
[2022-07-20 22:15:38.156][51][critical][wasm] [external/envoy/source/extensions/common/wasm/context.cc:1227] wasm log: panicked at 'called `Result::unwrap()` on an `Err` value: Error { kind: Unsupported, message: "operation not supported on this platform" }', src/lib.rs:65:67
```