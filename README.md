go-sentinel
===========

Redis Sentinel support for [redigo](https://github.com/gomodule/redigo) library.

Documentation
-------------

- [API Reference](http://godoc.org/github.com/FZambia/sentinel)

test
-------------

make cp
make up
go test -v ./...

Alternative solution
--------------------

You can alternatively configure Haproxy between your application and Redis to proxy requests to Redis master instance if you only need HA:

```
listen redis
    server redis-01 127.0.0.1:6380 check port 6380 check inter 2s weight 1 inter 2s downinter 5s rise 10 fall 2 on-marked-down shutdown-sessions on-marked-up shutdown-backup-sessions
    server redis-02 127.0.0.1:6381 check port 6381 check inter 2s weight 1 inter 2s downinter 5s rise 10 fall 2 backup
    bind *:6379
    mode tcp
    option tcpka
    option tcplog
    option tcp-check
    tcp-check send PING\r\n
    tcp-check expect string +PONG
    tcp-check send info\ replication\r\n
    tcp-check expect string role:master
    tcp-check send QUIT\r\n
    tcp-check expect string +OK
    balance roundrobin
```

This way you don't need to use this library.

License
-------

Library is available under the [Apache License, Version 2.0](http://www.apache.org/licenses/LICENSE-2.0.html).

```json
2020/11/24 16:15:44 map[
    down-after-milliseconds:30000
    flags:slave
    info-refresh:996
    ip:172.23.0.1
    last-ok-ping-reply:315
    last-ping-reply:315
    last-ping-sent:0
    link-pending-commands:0
    link-refcount:1
    master-host:172.25.184.229
    master-link-down-time:0
    master-link-status:ok
    master-port:6379
    name:172.23.0.1:6381
    port:6381
    role-reported:slave
    role-reported-time:1577043
    runid:4e2f5677d0465b001c30049240d706f2385c55d9
    slave-priority:100
    slave-repl-offset:364985]
```
