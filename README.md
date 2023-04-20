# couchbase-replication-ping a.k.a cb-tracker

Tool that allows to perform round-trip communication (write-read-overwrite-read) based on [XDCR](https://docs.couchbase.com/server/current/learn/clusters-and-availability/xdcr-overview.html) in order to measure replication latency in both directions simultaneously.

# Do you need a fast track to examine replication lag?

1. Prepare config that contains configuration that allows to connect to buckets: [example](./configuration/local.yaml)
2. There is ready to use docker:

```
CONFIG=/path/to/your/config.yaml
docker run -v ${CONFIG}:/config.yml ziollek/cb-tracker:latest oneway
```
Please check below information to know what else can you achieve and how to interpret command output.

## build

make build

## run

```
./bin/cp-tracker ping --config=configuration/local.yaml
```

notes: example config can be found in [configuration/local.yaml](./configuration/local.yaml)

## configuration file

The example configuration is presented below:

```yaml
source:
  uri: couchbase://localhost
  bucket: default
  user: default
  password: default
  name: src

destination:
  uri: couchbase://localhost
  bucket: default
  user: default
  password: default
  name: dst

generator:
  ttl: 150s
  size: 200

key: cb-repl-ping-document-key
```

source & destination keys has exactly the same structure and represents connection string to couchbase bucket, details are described below:

| key      | description                                                                |
|----------|----------------------------------------------------------------------------|
| uri      | couchbase url in format couchbase://[couchnase-host-ip-or-hostname]        |
| bucket   | bucket name                                                                |
| user     | credentials: username                                                      |
| password | credentials: password                                                      |
| name| name that is used internally to distinguished source & destinantion bucket |

additionally you can find here parameters that are used to generate documents while testing replication:

| key  | description                                                                         |
|------|-------------------------------------------------------------------------------------|
| ttl  | ttl of documents that are generated during tests |
| size | size in bytes of data field of documents that are generated during tests, it worth mentioning that document consists of several additional fields and are a little bigger than defined size|
| key  | depending of command it is key that is used for storing testing documents or prefix of documents key |


## what kind of tests/modes it supports?

### ping mode - checks full round trip replication time

1. PING

   a. store -> bucketA -> replication -> bucketB

   b. fetch <- bucketB (with retires)
2. PONG

   a. store -> bucketB -> replication -> bucketA

   b. fetch <- bucketA (with retires)

#### example output

```
./bin/cb-tracker ping --config=configuration/local.yaml

INFO[2023-04-20T17:45:51+02:00] Using config file: configuration/local.yaml
INFO[2023-04-20T17:46:02+02:00] Start measuring latency ...
INFO[2023-04-20T17:46:02+02:00] ping                                          no=1 ping=5.250046ms pong=3.247143ms total=8.497189ms
INFO[2023-04-20T17:46:03+02:00] ping                                          no=2 ping=6.287856ms pong=3.777971ms total=10.065827ms
INFO[2023-04-20T17:46:04+02:00] ping                                          no=3 ping=4.662636ms pong=3.46971ms total=8.132346ms
```

#### interpretation

| field   | description                                  |
|---------|----------------------------------------------|
| total   | time consumed by whole operation             |
| ping    | time consumed by ping phase                  |
| pong    | time consumed by pong phase                  |
| retries | total number of retrying reads on both sides |
| err | error message |

### oneway mode - check detailed times consumed by unidirectional replication

TBD

### half ping mode - run two processes that connected only with one side (source or destination)

This approach allows to mitigate variety of RTT while connecting to both sides of replication from single host

TBD


## remarks

Measured times are tightly connected with RTT time between the machine where the test is fired and both clusters/buckets.
If you want to achieve accurate results you should consider half-ping mode fired on hosts with little RTT to each side of testing buckets.