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
kubectl apply -f flask-app/flask-service.yaml
```

Make sure all the pods from these application are all running first, then we can verify if the services and deployments are working well.
```
# deploy a test pod to verify the cluster internal routing based on the service
kubectl apply -f echo-app/test-pod.yaml

kubectl exec -it test-pod -- /bin/bash

# inside the test-pod
# the 1st curl command should return the value of environment variable `version` we set in the deployment spec
# the 2nd curl command should dump out the request details.
$ curl http://flask-service.default.svc.cluster.local/env/version
v1

$ curl http://echo-service.default.svc.cluster.local
{"host":{"hostname":"echo-service.default.svc.cluster.local","ip":"::ffff:127.0.0.6","ips":[]},"http":{"method":"GET","baseUrl":"","originalUrl":"/","protocol":"http"},"request":{"params":{"0":"/"},"query":{},"cookies":{},"body":{},"headers":{"host":"echo-service.default.svc.cluster.local","user-agent":"curl/7.82.0","accept":"*/*","x-forwarded-proto":"http","x-request-id":"84fc4172-5cf5-4c32-9a44-25bc6af5fa02","x-envoy-attempt-count":"1","x-forwarded-client-cert":"By=spiffe://cluster.local/ns/default/sa/default;Hash=3f21fd3abb9bcccbf0c34bd1767941de33a8b19cfa31e85611724c8b75932e06;Subject=\"\";URI=spiffe://cluster.local/ns/default/sa/default","x-b3-traceid":"dfb248c38b4fb2e2b0dfa041783aaafc","x-b3-spanid":"2e571c2c1b66b0bc","x-b3-parentspanid":"b0dfa041783aaafc","x-b3-sampled":"0"}},"environment":{"PATH":"/usr/local/sbin:/usr/local/bin:/usr/sbin:/usr/bin:/sbin:/bin","HOSTNAME":"echo-deployment-6b7ddbcb47-qffcn","version":"v1","ECHO_SERVICE_SERVICE_PORT_HTTP":"80","ECHO_SERVICE_PORT_80_TCP_PORT":"80","KUBERNETES_PORT_443_TCP":"tcp://10.96.0.1:443","KUBERNETES_PORT_443_TCP_PROTO":"tcp","KUBERNETES_PORT_443_TCP_PORT":"443","ECHO_SERVICE_SERVICE_HOST":"10.96.118.49","ECHO_SERVICE_PORT":"tcp://10.96.118.49:80","ECHO_SERVICE_PORT_80_TCP_ADDR":"10.96.118.49","KUBERNETES_SERVICE_PORT":"443","KUBERNETES_PORT":"tcp://10.96.0.1:443","KUBERNETES_PORT_443_TCP_ADDR":"10.96.0.1","ECHO_SERVICE_SERVICE_PORT":"80","ECHO_SERVICE_PORT_80_TCP":"tcp://10.96.118.49:80","ECHO_SERVICE_PORT_80_TCP_PROTO":"tcp","KUBERNETES_SERVICE_HOST":"10.96.0.1","KUBERNETES_SERVICE_PORT_HTTPS":"443","NODE_VERSION":"16.15.0","YARN_VERSION":"1.22.18","HOME":"/root"}}
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
# this will create a gateway that adapt the traffic into the cluster for the requests whose host is 127.0.0.1
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
mkdir certs
openssl req -x509 -sha256 -nodes -days 365 -newkey rsa:2048 -subj '/O=example Inc./CN=example.com' -keyout certs/example.com.key -out certs/example.com.crt
```

2. Create a certificate and a private key for demo.example.com:
```
openssl req -out certs/demo.example.com.csr -newkey rsa:2048 -nodes -keyout certs/demo.example.com.key -subj "/CN=demo.example.com/O=demo organization"
openssl x509 -req -sha256 -days 365 -CA certs/example.com.crt -CAkey certs/example.com.key -set_serial 0 -in certs/demo.example.com.csr -out certs/demo.example.com.crt
```

3. Configure a TLS ingress gateway for a single host
```
# Create a secret for the ingress gateway
# Note: the namespace must be istio-system
kubectl create -n istio-system secret tls demo-credential --key=certs/demo.example.com.key --cert=certs/demo.example.com.crt

kubectl apply -f istio-resources/tls/gateway.yaml
```

Testing the new configuration with
```
curl -v -HHost:demo.example.com --resolve "demo.example.com:443:127.0.0.1" \
--cacert certs/example.com.crt "https://demo.example.com:443/flask/env/version"
```
The version value should return at last.

And if you not attach the ca cert
```
curl -v -HHost:demo.example.com --resolve "demo.example.com:443:127.0.0.1" "https://demo.example.com:443/flask/env/version"
```
you should see
```
* Added demo.example.com:443:127.0.0.1 to DNS cache
* Hostname demo.example.com was found in DNS cache
*   Trying 127.0.0.1:443...
* TCP_NODELAY set
* Connected to demo.example.com (127.0.0.1) port 443 (#0)
* ALPN, offering h2
* ALPN, offering http/1.1
* successfully set certificate verify locations:
*   CAfile: /etc/ssl/certs/ca-certificates.crt
  CApath: /etc/ssl/certs
* TLSv1.3 (OUT), TLS handshake, Client hello (1):
* TLSv1.3 (IN), TLS handshake, Server hello (2):
* TLSv1.3 (IN), TLS handshake, Encrypted Extensions (8):
* TLSv1.3 (IN), TLS handshake, Certificate (11):
* TLSv1.3 (OUT), TLS alert, unknown CA (560):
* SSL certificate problem: unable to get local issuer certificate
* Closing connection 0
curl: (60) SSL certificate problem: unable to get local issuer certificate
More details here: https://curl.haxx.se/docs/sslcerts.html

curl failed to verify the legitimacy of the server and therefore could not
establish a secure connection to it. To learn more about this situation and
how to fix it, please visit the web page mentioned above.
```

## Demo the mTLS

1. Create a self-signed certificate
```
openssl req -x509 -newkey rsa:4096 -keyout certs/key.pem -out certs/cert.pem -sha256 -days 365 -nodes -subj '/CN=example.com'
```

2. Create the secret to store this cert.
```
# the name must be `<secret>-cacert`. Since we create the previous secret as `demo-credential`, we will create this secrete as `demo-credential-cacert`
kubectl create -n istio-system secret generic demo-credential-cacert --from-file=ca.crt=certs/cert.pem
```

3. Update the gateway resource
```
kubectl apply -f istio-resources/tls/gateway.yaml
```

4. Verify the client request will fail due to it doesn't present the certificate
```
curl -v -HHost:demo.example.com --resolve "demo.example.com:443:127.0.0.1" \
--cacert certs/example.com.crt "https://demo.example.com:443/flask/env/version"
```
The command should contains
```
curl: (56) OpenSSL SSL_read: error:1409445C:SSL routines:ssl3_read_bytes:tlsv13 alert certificate required, errno 0
```
at the end.

5. Verify the client request will succeed after presenting the certificate
```
curl -v -HHost:demo.example.com --resolve "demo.example.com:443:127.0.0.1" \
--cacert certs/example.com.crt --cert certs/cert.pem --key certs/key.pem "https://demo.example.com:443/flask/env/version"
```
It should return the version number now

# Teardown the resources
```
# Local files
rm -rf certs

# delete minikube cluster
minikube delete --all
```