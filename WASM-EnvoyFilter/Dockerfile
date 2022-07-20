FROM sfbrdevhelmweacr.azurecr.io/rancher/mirrored-istio-proxyv2:1.10.4
RUN apt-get update && apt-get install -y \
    vim

ENTRYPOINT ["/usr/local/bin/envoy"]