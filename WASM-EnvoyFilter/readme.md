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
tinygo build -o plugin.wasm -scheduler=none -target=wasi main.go

docker build . -t envoy:wasm
```