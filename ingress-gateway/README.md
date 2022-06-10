# Overview
This demo will show how to create an ingress gateway and set up the traffic routing via using Istio.
It focuses on these features provided by the Istio
- Ingress Gateway
- Traffic Routing and Traffic Control
- Request overwriting
- TLS/mTLS

For this demo, it will deploy two apps
- A simple flask app that can return the environment variable values
- A simple echo server app that can dump out the request details

This demo uses minikube to deploy a local k8s cluster in a Linux system.
Ideally it can be reproduced in any linux systems and x86_64 MacOS. Arm64 MacOS(M1 Chip) cannot deploy Istio locally for now.

# Demo Details

## Software installation and Cluster Creation
As for the prerequisites, these tools are assuming to be installed first
- kubectl
- Docker

We need to install minikube. Here I am using homebrew to install it.
```
# You can use this command to install homebrew on a linux system if you don't have it.
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

brew install minikube
```

Once the minikube is installed, use this command to create a local k8s cluster
```
minikube start
```

This command will help you to create a local k8s cluster with a single control plane node, which is enough for our demo.
The kubectl context will be set to this minikube cluster by default. so if you run
```
kubectl get nodes
```
you should see the node information like this
```
NAME       STATUS   ROLES                  AGE   VERSION
minikube   Ready    control-plane,master   87s   v1.23.3
```

Next, we will install the Istio inside this cluster.
```
# use homebrew to install istioctl
brew install istioctl

# install Istio inside the cluster, type 'y' for the prompt.
# The default profile will deploy ["Istio core" "Istiod" "Ingress gateways"], which are enough for us.
istioctl manifest apply --set profile=default
```

Now you can see the Istio components are installed inside the cluster
```
$ kubectl get pods -n istio-system
NAME                                    READY   STATUS    RESTARTS   AGE
istio-ingressgateway-75c6d79fcc-5n4gr   1/1     Running   0          12s
istiod-7ccff5bbc7-gfsqt                 1/1     Running   0          24s


$ kubectl get svc -n istio-system
NAME                   TYPE           CLUSTER-IP      EXTERNAL-IP   PORT(S)                                      AGE
istio-ingressgateway   LoadBalancer   10.96.158.150   <pending>     15021:30845/TCP,80:32732/TCP,443:31281/TCP   7m25s
istiod                 ClusterIP      10.107.96.99    <none>        15010/TCP,15012/TCP,443/TCP,15014/TCP        7m37s
```
Look at the `EXTERNAL-IP` of `istio-ingresgateway`, it shows `<pending>` for now. No worries, we will fix that later.
The reason behind that is this svc is defined as a `LoadBalancer` type, the external ip is usually provided by the cloud provider. As this is just a local cluster, cannot getting an external IP is expected here.

## Applications Deployment
Istio injection will help you to inject the istio-init and istio-proxy containers for your workloads automatically.
For the simplicities, we will deploy all our applications in the default namespace. Hence, we can label the namespace: `default` to enable the istio injection.
```
kubectl label namespace default istio-injection=enabled
```

```
# Deploy the echo app
kubectl apply -f echo-app/echo-deployment.yaml
kubectl apply -f echo-app/echo-service.yaml

# Deploy the flask app.
# We will deploy two versions of the flask app deployments, in order to demo the traffic control in the later.
kubectl apply -f flask-app/flask-deployment-v1.yaml
kubectl apply -f flask-app/flask-deployment-v2.yaml
kubectl apply -f flask-app/service.yaml
```

Make sure all the pods from these application are all running first, then we can verify if the services and deployments are working well.
```
# deploy a test pod to verify the cluster internal routing based on the service
kubectl apply -f echo-app/test-pod.yaml

kubectl exec -it test-pod -- /bin/bash

# inside the test-pod
# the 1st curl command should return the value of environment variable `version` we set in the deployment spec
# the 2nd curl command should dump out the request details.
curl http://flask-service.default.svc.cluster.local/env/version
curl http://echo-service.default.svc.cluster.local
```
You can exit this session now.

## Demo the Ingress Gateway
Remember our ingress loadbalancer is trying to get an external IP address?
We need to open another terminal session, and run this command
```
minikube tunnel --alsologtostderr
```
And leave this terminal session there, do not kill it!
If you check back that loadbalancer type of service, it should get an IP address now.

For me, it is using 127.0.0.1. So I am going to configure my resources with hosts containing this ip.

In order to achieve the full functionalities of routing an outside request to the mesh, we need two resources
- Gateway
- VirtualService

```
# this will create a gateway that adapt the traffic into the cluster for the requests whose hosts are localhost or 127.0.0.1
kubectl apply -f istio-resources/ingress-gateway/gateway.yaml

# this will create a virtualservice that routes the request to the services should serve them.
kubectl apply -f istio-resources/ingress-gateway/virtualservice.yaml
```

Now we can verify if the requests can be handled correctly
```
curl 127.0.0.1/flask/env/version
curl 127.0.0.1/echo/
```
The both curl commands should return what shows previously inside the test pod.

## Demo the traffic control
VirtualService + DestinationRule can help us to control the traffic path. For example, we can have both version of flask app, and VirtualService + DestinationRule can help us route X% of traffic to the v1 flask app and 1-X% of the traffic to the v2 flask app. This function is pretty useful for canary deployment.

As we already deployment both versions of flask app, the only thing we need to do is tweak the configuration of our virtual service and adding a destination rule.

```
kubectl apply -f istio-resources/traffic-control/virtualservice.yaml
kubectl apply -f istio-resources/traffic-control/destination-rule.yaml
```

If you run multiple times of
```
curl 127.0.0.1/flask/env/version
```
the responses could be either `v1` or `v2`, as we configured the traffics are evenly distributed between them.

You can also change one version weight as 100 and the other's as 0, then the curl command should return on the version whose weight is 100.

## Demo the Delegate Virtual Service
As we are trying to replace the nginx gateway by istio gateway, we pretty much like all our contribution experiences can be as much similar as before.

We used to take an Ingress resource to define the routing strategy, and a VirtualService resource to control the canary deployment. After using the pure Istio, we can still achieve that by using the delegate virtual service.

We can have a top layer of VS that behaves like the Ingress previously to define the general routing strategy, and an second layer of VS to control the canary deployment.

```
kubectl delete virtualservice demo-vs
kubectl apply -f istio-resources/delegate-vs/ingress-vs.yaml
kubectl apply -f istio-resources/delegate-vs/vs-canary.yaml
```
First, we delete the previous vs.
Then we create the top layer of the vs, inside it, we use delegate for flask app, and keep the original routing strategy for echo app.
Last, we create a dedicated vs for control flask's traffic.

## Demo the tls
This section shows how to expose a secure HTTPS service using simple tls.

1. Create a root certificate and private key to sign the certificates for the services:
```
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=example Inc./CN=example.com' -keyout example.com.key -out example.com.crt
```

2. Create a certificate and a private key for demo.example.com:
```
openssl req -out demo.example.com.csr -newkey rsa:2048 -nodes -keyout demo.example.com.key -subj "/CN=demo.example.com/O=demo organization"
openssl x509 -req -sha256 -days 365 -CA example.com.crt -CAkey example.com.key -set_serial 0 -in demo.example.com.csr -out demo.example.com.crt
```

3. Configure a TLS ingress gateway for a single host
