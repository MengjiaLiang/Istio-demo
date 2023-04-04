# CPU-Profiling for Envoy
This page will talk about how to do CPU-Profiling on the envoy proxy in istio ingressgateway.
Here is the [reference](https://github.com/istio/istio/wiki/Analyzing-Istio-Performance) from Istio Wiki.

Istioctl seems not support enable CPU profiling yet, that means we need to use the [envoy admin interface](https://www.envoyproxy.io/docs/envoy/latest/operations/admin) to enable CPU profiling via the cURL command.

## Pre-requisites
1. Enable coreDump when install the istio, because CPU profiling needs write to the file system, this permission is disabled by default.
2. Use a distro image of envoy proxy instead of a distroless image, otherwise you cannot use a cURL command to enable the profiling.

## Instructions
1. Create a minikube cluster
```
minikube start --kubernetes-version=v1.22.2
```

2. Install Istioctl and install Istio components into the cluster
```
export ISTIO_VERSION=1.14.1
curl -sL https://istio.io/downloadIstioctl | sh -

istioctl manifest apply --set values.global.proxy.enableCoreDump=true
```

3. Replace the proxy image in the deployment - `istio-ingressgateway`

```
kubectl edit deploy istio-ingressgateway -n istio-system
```
change `docker.io/istio/proxyv2:1.14.1` to `sfbrdevhelmweacr.azurecr.io/wasm-proxy:1.14.1`


4. Install testing required resources
We will use the routing E2E test framework as an example. It will install
- Gateway
- Dummy LS resources (deployment , configmap)
- Echo server resources (deployment, virtual service)
- EnvoyFilters

```
kubectl create namespace uipath --dry-run=client -o yaml | kubectl apply -f -
kubectl label namespace uipath istio-injection=enabled
kubectl label namespace default istio-injection=enabled
kubectl create configmap service-standard-url-status -n uipath --from-file=./resources/service-standard-url-status.json
kubectl apply -f ./resources/gateway.yml
kubectl apply -f ./resources/echo-server.yml
kubectl apply -f ./resources/dummy-ls.yml
kubectl wait deployment -n uipath  platform-location-service --for condition=Available=True --timeout=2m
kubectl apply -f ./resources/wasm-plugin.yml
```

5. Enable minikube tunnel
Open another terminal session
```
minikube tunnel --alsologtostderr
```
you may be asked for sudo password as the loadbalancer of istio-ingressgateway uses 80 port.

6. Enable wasm log level to verify if everything works fine
```
# pod name is the pod of istio-ingressgateway.  E.g.  istio-ingressgateway-6854d49464-bcprt
istioctl proxy-config log <pod-name>.istio-system --level wasm:info
```

Check the log of this pod, see if anything complains errors.
```
curl -i 127.0.0.1/mjorg/portal_

kubectl logs -n istio-system <pod-name>
```
if you don't find any erros, and the wasm info log shows some LS calls and their response, you should be good.

You may need to delete the HPA of istio-ingresssgateway to let the single pod be overwhelmed to collect more data.

7. Start load testing
For simplicity of the load testing, we will use this tool. https://github.com/link1st/go-stress-testing
Clone it to you local first.
```
go mod tidy
go mod download

go run main.go -c 10 -n 20000 -u http://127.0.0.1/mjorg/portal_
```

8. CPU Profiling
When you detect your isio-ingressgateway pod starts to consume more resource

```
 kubectl exec -it -n istio-system <pod> -- /bin/bash

# Inside the pod
curl -X POST -s "http://localhost:15000/cpuprofiler?enable=y"

# wait for the time you want to profile.
curl -X POST -s "http://localhost:15000/cpuprofiler?enable=n"

exit

# in your local session
POD="istio-ingressgateway-787968d58d-hgmrg"
kubectl cp -n istio-system "$POD":/var/lib/istio/data /tmp/envoy -c istio-proxy
kubectl cp -n istio-system "$POD":/lib/x86_64-linux-gnu /tmp/envoy/lib -c istio-proxy
kubectl cp -n istio-system "$POD":/usr/local/bin/envoy /tmp/envoy/lib/envoy -c istio-proxy
```

Install pprof
```
go install github.com/google/pprof@latest
```
The executable file is in $HOME/go/bin, do `export PATH=$PATH:$(go env GOPATH)/bin`

Install other dependencies that pprof requires
```
brew install graphviz
```

Run pprof
```
PPROF_BINARY_PATH=/tmp/envoy/lib/ pprof -pdf /tmp/envoy/lib/envoy /tmp/envoy/envoy.prof
```