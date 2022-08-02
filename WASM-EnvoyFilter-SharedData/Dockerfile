FROM  gcr.io/gke-release/asm/proxyv2:1.10.2-asm.2
RUN apt-get update && apt-get install -y \
    vim

COPY ./worker.wasm /tmp/worker.wasm
COPY ./singleton.wasm /tmp/singleton.wasm

ENTRYPOINT ["/usr/local/bin/envoy"]