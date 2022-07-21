## Envoy luasocket test

Note that this is not recommended, since `luasocket` is blocking.

## Build and run

### Build the container image
```bash
$ make
```

### Run with the LuaSockets module imported
```bash
$ make run-luasocket-module
```
If you are seeing the logs like
```
qemu: uncaught target signal 6 (Aborted) - core dumped
make: *** [run-luasocket-module] Error 137
```
at the end, you need to check the log details to see what happened and fix it. If there's no obvious error logs inside like, you can kill the process by `docker stop` and redo the same command.

If you are seeing the logs like
```
[2022-02-24 18:17:40.664][1][info][config] [source/server/listener_manager_impl.cc:888] all dependencies initialized. starting workers
[2022-02-24 18:17:40.672][1][warning][main] [source/server/server.cc:642] there is no configured limit to the number of allowed active connections. Set a limit via the runtime key overload.global_downstream_max_connections
```
at the end, you can move to next step. This means your envoy is up and running inside the container.

// From another terminal session
```bash
$ curl localhost:8000
{"args":{},"headers":{"x-forwarded-proto":"https","x-forwarded-port":"443","host":"localhost","x-amzn-trace-id":"Root=1-5f6d39ed-75c8f0de990ca942e91ab258","content-length":"0","user-agent":"curl/7.58.0","accept":"*/*","x-request-id":"a3f82815-15cb-42cd-b4bd-57026c10443a","foo":"LuaSocket 3.0-rc1","x-envoy-expected-rq-timeout-ms":"15000","x-envoy-original-path":"/"},"url":"https://localhost/get"}
```

From the response, you can see that `"foo":"LuaSocket 3.0-rc1"` is set by Lua script.


### Run with the redis connection in Lua
Start a redis in background then start the envoy
```bash
$ make redis
$ make run-redis-connection
```

Same as what is mentioned above, you can send the request only when the log shows your envoy is up and running well.

// From another terminal session
```bash
$ curl localhost:8000
{"args":{},"headers":{"x-forwarded-proto":"https","x-forwarded-port":"443","host":"localhost","x-amzn-trace-id":"Root=1-6217cbc8-2df1640c4df99fef770db9d2","user-agent":"curl/7.64.1","accept":"*/*","x-request-id":"1d163471-f62d-403e-86f4-2bc9fc6a0e48","foo":"bar from redis","x-envoy-expected-rq-timeout-ms":"15000","x-envoy-original-path":"/"},"url":"https://localhost/get"}
```

From the response, you can see that `"foo":"bar from redis"` is set by Lua script. The `bar` is the value of the key `foo` in redis.

