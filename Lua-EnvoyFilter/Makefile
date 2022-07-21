HUB ?= mengjial
PORT ?= 8000

build:
	docker build . -t $(HUB)/lua-envoy-test:latest

redis:
	docker run --rm --name redis -p 6379:6379 -d redis

run-luasocket-module:
	docker run --rm -it -p $(PORT):8000 -t $(HUB)/lua-envoy-test:latest -c /etc/envoy/envoy-luasocket-module.yaml

run-redis-connection:
	docker run --rm -it -p $(PORT):8000 -t $(HUB)/lua-envoy-test:latest -c /etc/envoy/envoy-redis-connection.yaml