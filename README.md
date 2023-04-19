# couchbase-replication-ping

Tool that allows to perform round-trip communication (write-read-overwrite-read) based on [XDCR](https://docs.couchbase.com/server/current/learn/clusters-and-availability/xdcr-overview.html) in order to measure replication latency in both directions simultaneously.

## build

make build

## run

```
cd bin
./cp-repl ping --config=../configuration/local.yaml
```

notes: example config can be found in [configuration/local.yaml](./configuration/local.yaml)

## output

```
Using config file: ../configuration/local.yaml
Start measuring latency ...
[2023-04-19 22:16:10] ping 1)	duration: 5.991002ms, err: %!s(<nil>)
[2023-04-19 22:16:11] ping 2)	duration: 7.77831ms, err: %!s(<nil>)
[2023-04-19 22:16:12] ping 3)	duration: 5.46673ms, err: %!s(<nil>)
```
