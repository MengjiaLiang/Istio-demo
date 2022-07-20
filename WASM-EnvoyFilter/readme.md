# The demo of using WASM in Istio

Envoy supports to use proxy-wasm inside its filter.
Before Istio v1.12, we can use WASM inside EnvoyFilter

This demo will show how to achieve that locally.

## Development Environment Setup
```
# 1. Install golang
brew install golang

# 2. Install TinyGo
wget https://github.com/tinygo-org/tinygo/releases/download/v0.23.0/tinygo_0.23.0_amd64.deb
sudo dpkg -i tinygo_0.23.0_amd64.deb
export PATH=$PATH:/usr/local/bin

# 3. Install Istioctl v1.10
curl -LO https://storage.googleapis.com/gke-release/asm/istio-1.10.2-asm.2-linux-amd64.tar.gz
tar xvf istio-1.10.2-asm.2-linux-amd64.tar.gz
mv istio-1.10.2-asm.2/bin/istioctl /usr/local/bin

```


## Testing Environment Setup
```
# 1. Setpup the minikube cluster
minikube start

# 2. Install the Istio
istioctl manifest apply --set profile=default

# 3. Lable the Namespace with Injection
kubectl label namespace default istio-injection=enabled
```

## Build the WASM code
```
# Set environment variables get from
tinygo env

# Download pkg
cd wasm.tinygo
go mod edit -require=github.com/tetratelabs/proxy-wasm-go-sdk@main
go mod download github.com/tetratelabs/proxy-wasm-go-sdk

# Build go binary
tinygo build -o plugin.wasm -scheduler=none -target=wasi main.go
```

## Test with local envoy
```
# Go to root of `WASM-ENVOYFILTER`
cd ..

docker run --rm -it \
           -v $(pwd)/local-envoy-config.yaml:/tmp/local-envoy-config.yaml \
           -v $(pwd)/wasm-tinygo/plugin.wasm:/tmp/plugin.wasm \
           -p 8000:8000 -t envoy:wasm -c /tmp/local-envoy-config.yaml
```

## WASM-Rust
```
# install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh

# install wasm-pack
curl https://rustwasm.github.io/wasm-pack/installer/init.sh -sSf | sh

# install packages
rustup toolchain install nightly
rustup target add wasm32-unknown-unknown

# for linux you probably need to do this as well
sudo apt-get update
sudo apt install build-essential

# Build the Rust wasm filter
cargo build --target=wasm32-unknown-unknown --release

docker build . -t envoy:wasm

docker run --rm -it \
           -v $(pwd)/local-envoy-config.yaml:/tmp/local-envoy-config.yaml \
           -v /Users/mengjia.liang/Documents/github.com/Sophichia/Istio-demo/wasm-rust-filter/target/wasm32-unknown-unknown/release/plugin.wasm:/tmp/plugin.wasm \
           -p 8000:8000 -t envoy:wasm -c /tmp/local-envoy-config.yaml
```