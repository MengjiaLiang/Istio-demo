# replace http header

Lua:
- QPS: 770
- CPU: 35m -> 458m
- Memory: 69Mi -> 69Mi
- Avg response time: 13.44

WASM:
- QPS: 770
- CPU: 29m -> 474m
- Memory: 88Mi -> 93Mi
- Avg response time: 12.84

# recontruct path
Lua:
- QPS: 730
- CPU: 13m -> 466m
- Memory: 78Mi -> 80Mi
- Avg response time: 13.64

WASM:
- QPS: 730
- CPU: 13m -> 468m
- Memory: 88Mi -> 104Mi
- Avg response time: 13.6

# dispatch http call
Lua:
- QPS: 378
- CPU: 9m -> 470m
- Memory: 52Mi -> 60Mi
- Avg response time: -> 26.45 

WASM:
- QPS: 485
- CPU: 3m -> 461Mi
- Memory: 56Mi -> 108Mi
- Avg response time: 20.59


# regex
Lua:
- QPS: 687
- CPU: 3m -> 472m
- Memory: 57Mi -> 58Mi 
- Avg response time: 14.55

WASM:
- QPS: 403
- CPU: 4m -> 856m
- Memory: 60Mi -> 120Mi
- Avg response time: 24.7

# json
Lua:
- QPS: 687
- CPU: 13m -> 454m
- Memory: 55Mi -> 57Mi 
- Avg response time: 14.55

WASM:
- QPS: 703
- CPU: 4m -> 461m
- Memory: 60Mi -> 100Mi
- Avg response time: 14.22

# shared data
WASM:
- QPS: 701
- CPU: 4m -> 474m
- Memory: 62Mi -> 92Mi 
- Avg response time: 14.24