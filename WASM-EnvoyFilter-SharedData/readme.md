# Deploy in Envoy

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

# Deploy in Istio
```
minikube start

# Install Istioctl v1.10
curl -LO https://storage.googleapis.com/gke-release/asm/istio-1.10.2-asm.2-linux-amd64.tar.gz
tar xvf istio-1.10.2-asm.2-linux-amd64.tar.gz
mv istio-1.10.2-asm.2/bin/istioctl /usr/local/bin

# Label the default namespace
istioctl manifest apply --set profile=default

# Deploy a basic workload
kubectl apply -f ../ingress-gateway/echo-app/echo-deployment.yaml
kubectl apply -f ../ingress-gateway/echo-app/echo-service.yaml

# In another terminal
minikube tunnel --alsologtostderr

# this will create a gateway that adapt the traffic into the cluster for the requests whose host is 127.0.0.1
kubectl apply -f ../ingress-gateway/istio-resources/ingress-gateway/gateway.yaml

# this will create a virtualservice that routes the request to the services should serve them.
kubectl apply -f ../ingress-gateway/istio-resources/ingress-gateway/virtualservice.yaml

# this command should dump the request details
curl 127.0.0.1/echo/


# build the wasm binaries
tinygo build -o worker.wasm -scheduler=none -target=wasi worker/main.go
tinygo build -o singleton.wasm -scheduler=none -target=wasi singleton/main.go

kubectl delete cm -n istio-system worker-filter --ignore-not-found
kubectl create cm -n istio-system worker-filter --from-file=./worker.wasm

kubectl delete cm -n istio-system singleton-filter --ignore-not-found
kubectl create cm -n istio-system singleton-filter --from-file=./singleton.wasm

# Mount the wasm binaries in ingress controller
kubectl patch deployment -n istio-system istio-ingressgateway --patch-file patch.yaml
```