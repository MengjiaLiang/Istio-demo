FROM mengjial/lua-envoy:v1.18.4-exported-modules

RUN apt-get -y update && \
    apt-get -y install luarocks wget && \
    luarocks install luasocket && \
    wget https://raw.githubusercontent.com/rxi/json.lua/master/json.lua -O /usr/local/share/lua/5.1/json.lua && \
    wget https://raw.githubusercontent.com/nrk/redis-lua/version-2.0/src/redis.lua  -O /usr/local/share/lua/5.1/redis.lua

COPY ./envoy-luasocket-module/envoy-luasocket-module.yaml /etc/envoy/envoy-luasocket-module.yaml
COPY ./envoy-redis-connection/envoy-redis-connection.yaml /etc/envoy/envoy-redis-connection.yaml
ENTRYPOINT ["/usr/local/bin/envoy"]
