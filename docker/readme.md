# Cockroach DB for local Development

## Single Node DB

```sh
docker-compose up -d
docker-compose logs
docker-compose down --remove-orphans
```

## Init 3-Node DB Cluster

```sh
# 1. bootstrap db nodes
docker-compose -f docker-compose-cluster.yml up -d
# 2. init db cluster
docker-compose -f docker-compose-cluster.yml exec roach1 ./cockroach init --insecure
docker-compose -f docker-compose-cluster.yml logs
docker-compose -f docker-compose-cluster.yml down --remove-orphans
```

##  Startup Details

`docker-compose -f docker-compose-cluster.yml exec roach1 grep 'node starting' cockroach-data/logs/cockroach.log -A 11`

```txt
cockroach-data/logs/cockroach.log -A 11
CockroachDB node starting at 2021-04-27 20:39:01.201001562 +0000 UTC (took 4.6s)
build:               CCL v20.2.8 @ 2021/04/23 13:54:57 (go1.13.14)
webui:               ‹http://roach1:8080›
sql:                 ‹postgresql://root@roach1:26257?sslmode=disable›
RPC client flags:    ‹/cockroach/cockroach <client cmd> --host=roach1:26257 --insecure›
logs:                ‹/cockroach/cockroach-data/logs›
temp dir:            ‹/cockroach/cockroach-data/cockroach-temp024677352›
external I/O path:   ‹/cockroach/cockroach-data/extern›
store[0]:            ‹path=/cockroach/cockroach-data›
storage engine:      pebble
status:              restarted pre-existing node
clusterID:           ‹12582bc7-633d-4050-aff8-11a9c0955344›
```

## Built-in SQL client

`docker-compose -f docker-compose-cluster.yml exec roach1 ./cockroach sql --insecure`

```sh
# Welcome to the CockroachDB SQL shell.
# All statements must be terminated by a semicolon.
# To exit, type: \q.
# Enter \? for a brief introduction.
#
root@:26257/defaultdb> CREATE DATABASE bank;
CREATE DATABASE

Time: 39ms total (execution 39ms / network 0ms)
```

## Simulating client traffic

CockroachDB comes with a number of [built-in workloads](https://www.cockroachlabs.com/docs/v20.2/cockroach-workload) for simulating client traffic.

### Load the initial dataset

`docker-compose -f docker-compose-cluster.yml exec roach1 ./cockroach workload init movr 'postgresql://root@roach1:26257?sslmode=disable'`

```sh
workload/workloadsql/dataload.go:140  imported users (0s, 50 rows)
workload/workloadsql/dataload.go:140  imported vehicles (0s, 15 rows)
workload/workloadsql/dataload.go:140  imported rides (0s, 500 rows)
workload/workloadsql/dataload.go:140  imported vehicle_location_histories (0s, 1000 rows)
workload/workloadsql/dataload.go:140  imported promo_codes (0s, 1000 rows)
workload/workloadsql/workloadsql.go:113  starting 8 splits
workload/workloadsql/workloadsql.go:113  starting 8 splits
workload/workloadsql/workloadsql.go:113  starting 8 splits
```

### Run the workload for 5 minutes

`docker-compose -f docker-compose-cluster.yml exec roach1 ./cockroach workload run movr --duration=5m 'postgresql://root@roach1:26257?sslmode=disable'`
