
# 1. Install Nginx
```
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx

NAMESPACE=nginx

helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update

helm install ingress-nginx ingress-nginx/ingress-nginx \
  --create-namespace \
  --namespace $NAMESPACE \
  --set controller.service.annotations."service\.beta\.kubernetes\.io/azure-load-balancer-health-probe-request-path"=/healthz \
  --debug
```

# 2. Create echo service
`kubectl apply -f echo.yaml`

# 3. Create Ingress rule
`kubectl apply -f ingress-v1.yaml`
Note: Go to the Azure resource group to give a DNS name of the LB of Nginx ingress control service. In my case, it is `mj-nginx-istio.eastus.cloudapp.azure.com` 

Now the routing based on Nginx should work
`curl -i mj-nginx-istio.eastus.cloudapp.azure.com/echo_`

# 4. Install istio

Use istioctl to install istio components. This command will install istio-core, istiod and ingress gateway
`istioctl manifest apply --set profile=default`

Label the namespace where the echo server is
`kubectl label namespace default istio-injection=enabled`

Restart the echo server
`kubectl rollout restart Deployment -n istio-system istio-ingressgateway`

# 5. Set Istio gateway, VS and envoy filter
The Istio ingress gateway svc should be a LB by default, so it has a public IP address. In my case, it is 104.45.186.250, put it into the gateway and virtual service
```
kubectl apply -f istio/v1/gateway.yaml
kubectl apply -f istio/v1/vs.yaml
kubectl apply -f istio/envoyfilter.yaml
```
This should set up a naive routing rule and enable a naive envoy filter to log the current host and path of a inbound request.

Let's enable the lua log of the ingress gateway controller
`istioctl proxy-config log istio-ingressgateway-6cf59d9885-l6tbq.istio-system --level lua:info`

Send a request
`curl -i 104.45.186.250/echo_/foo`
It should returns 200

Check the istio ingress gateway log
`kubectl logs -n istio-system istio-ingressgateway-6cf59d9885-l6tbq`
it should contain something like
```
2023-05-03T21:53:49.886354Z	info	envoy lua	script log: Current path is: /echo_/foo
2023-05-03T21:53:49.886407Z	info	envoy lua	script log: Current Host is: 104.45.186.250
```

This demonstrate a basic istio ingress routing configuration.
In next section, we will let the nginx control the ingress prior to the whole istio components.

# 6. Let nginx delegate the ingress for istio
delete old resources
```
kubectl delete -f ingress-v1.yaml
kubectl delete -f istio/v1/gateway.yaml
kubectl delete -f istio/v1/vs.yaml
```

apply new resources
```
kubectl apply -f ingress-v2.yaml
kubectl apply -f istio/v2/gateway.yaml
kubectl apply -f istio/v2/vs.yaml
```

now send a request
`curl -i mj-nginx-istio.eastus.cloudapp.azure.com/echo_/bar`
it still returns 200

and if you check the log of istio ingress gateway log, the envoy filter works as expected
```
2023-05-03T22:00:50.328724Z	info	envoy lua	script log: Current path is: /echo_/bar
2023-05-03T22:00:50.328761Z	info	envoy lua	script log: Current Host is: mj-nginx-istio.eastus.cloudapp.azure.com
```