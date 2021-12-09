# Cockroach DB for local Development

Docs:

- [Install Docker-Compose](install.md)
- [Start a Cluster in Docker](https://www.cockroachlabs.com/docs/v21.1/start-a-local-cluster-in-docker-linux)
- [Cockroach Commands](https://www.cockroachlabs.com/docs/v21.2/cockroach-commands.html)

## Single Node DB

```sh
cd crdb/docker
docker-compose up -d
docker-compose logs -f
docker-compose down --remove-orphans --volumes
```

## Init 3-Node DB Cluster

```sh
docker-compose -f docker-compose-cluster.yml up -d
# docker-compose -f docker-compose-cluster.yml exec roach1 ./cockroach init --insecure --host=roach1:26257 --cluster-name=tulloroach
docker-compose -f docker-compose-cluster.yml logs
docker-compose -f docker-compose-cluster.yml down --remove-orphans --volumes
```

##  Startup Details

```sh
# docker-compose -f docker-compose-cluster.yml exec roach1 grep 'node starting' cockroach-data/logs/cockroach.log -A 13
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +CockroachDB node starting at 2021-12-09 12:36:24.146899839 +0000 UTC (took 5.4s)
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +build:               CCL v21.2.2 @ 2021/12/01 14:35:45 (go1.16.6)
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +webui:               ‹http://roach1:8080›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +sql:                 ‹postgresql://root@roach1:26257/defaultdb?sslmode=disable›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +sql (JDBC):          ‹jdbc:postgresql://roach1:26257/defaultdb?sslmode=disable&user=root›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +RPC client flags:    ‹/cockroach/cockroach <client cmd> --host=roach1:26257 --insecure›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +logs:                ‹/cockroach/cockroach-data/logs›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +temp dir:            ‹/cockroach/cockroach-data/cockroach-temp785232904›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +external I/O path:   ‹/cockroach/cockroach-data/extern›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +store[0]:            ‹path=/cockroach/cockroach-data›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +storage engine:      pebble
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +status:              restarted pre-existing node
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +cluster name:        ‹tulloroach›
I211209 12:36:24.147803 65 1@cli/start.go:759 ⋮ [-] 80 +clusterID:           ‹9bb62702-fd6f-43cb-9aa5-4be023757c4f›
```

## Built-in SQL client

```sh
# docker-compose -f docker-compose-cluster.yml exec roach1 ./cockroach sql --insecure
root@:26257/defaultdb> CREATE USER IF NOT EXISTS johndoe; CREATE DATABASE IF NOT EXISTS bank; GRANT ALL ON DATABASE bank TO johndoe;

root@:26257/defaultdb> show databases;
  database_name | owner | primary_region | regions | survival_goal
----------------+-------+----------------+---------+----------------
  bank          | root  | NULL           | {}      | NULL
  defaultdb     | root  | NULL           | {}      | NULL
  postgres      | root  | NULL           | {}      | NULL
  system        | node  | NULL           | {}      | NULL
(4 rows)
```

---

## Simulating client traffic

CockroachDB comes with a number of [built-in workloads](https://www.cockroachlabs.com/docs/v21.2/cockroach-workload) for simulating client traffic.

### Load the initial dataset

```sh
# docker-compose -f docker-compose-cluster.yml exec roach1 ./cockroach workload init movr 'postgresql://root@roach1:26257?sslmode=disable'
I211209 12:48:07.010159 1 workload/workloadsql/dataload.go:146  [-] 1  imported users (0s, 50 rows)
I211209 12:48:07.037228 1 workload/workloadsql/dataload.go:146  [-] 2  imported vehicles (0s, 15 rows)
I211209 12:48:07.119618 1 workload/workloadsql/dataload.go:146  [-] 3  imported rides (0s, 500 rows)
I211209 12:48:07.191365 1 workload/workloadsql/dataload.go:146  [-] 4  imported vehicle_location_histories (0s, 1000 rows)
I211209 12:48:07.249157 1 workload/workloadsql/dataload.go:146  [-] 5  imported promo_codes (0s, 1000 rows)
I211209 12:48:07.254055 1 workload/workloadsql/workloadsql.go:113  [-] 6  starting 8 splits
I211209 12:48:07.587873 1 workload/workloadsql/workloadsql.go:113  [-] 7  starting 8 splits
I211209 12:48:07.916997 1 workload/workloadsql/workloadsql.go:113  [-] 8  starting 8 splits
```

### Run the workload for 5 minutes

```sh
# docker-compose -f docker-compose-cluster.yml exec roach1 ./cockroach workload run movr --duration=5m 'postgresql://root@roach1:26257?sslmode=disable'
...
  136.0s        0            0.0            0.8      0.0      0.0      0.0      0.0 addUser
  136.0s        0            0.0            0.2      0.0      0.0      0.0      0.0 addVehicle
  136.0s        0            0.0            0.3      0.0      0.0      0.0      0.0 applyPromoCode
  136.0s        0            0.0            0.1      0.0      0.0      0.0      0.0 createPromoCode
  136.0s        0            0.0            0.3      0.0      0.0      0.0      0.0 endRide
  136.0s        0            0.0           51.2      0.0      0.0      0.0      0.0 readVehicles
  136.0s        0            0.0            1.0      0.0      0.0      0.0      0.0 startRide
  136.0s        0            0.0            2.6      0.0      0.0      0.0      0.0 updateActiveRides
  137.1s        0            0.0            0.8      0.0      0.0      0.0      0.0 addUser
  137.1s        0            0.0            0.2      0.0      0.0      0.0      0.0 addVehicle
  137.1s        0            0.0            0.3      0.0      0.0      0.0      0.0 applyPromoCode
  137.1s        0            0.0            0.1      0.0      0.0      0.0      0.0 createPromoCode
  137.1s        0            0.0            0.3      0.0      0.0      0.0      0.0 endRide
  137.1s        0            0.0           50.8      0.0      0.0      0.0      0.0 readVehicles
  137.1s        0            0.0            0.9      0.0      0.0      0.0      0.0 startRide
  137.1s        0            0.0            2.6      0.0      0.0      0.0      0.0 updateActiveRides
  138.0s        0            0.0            0.8      0.0      0.0      0.0      0.0 addUser
  138.0s        0            0.0            0.2      0.0      0.0      0.0      0.0 addVehicle
  138.0s        0            0.0            0.3      0.0      0.0      0.0      0.0 applyPromoCode
  138.0s        0            0.0            0.1      0.0      0.0      0.0      0.0 createPromoCode
_elapsed___errors__ops/sec(inst)___ops/sec(cum)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)
...
```
