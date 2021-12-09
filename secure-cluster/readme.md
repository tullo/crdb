# Local Secure Cluster

## Step 1. [Generate certificates](https://www.cockroachlabs.com/docs/v21.1/secure-a-cluster#step-1-generate-certificates)

```sh
make dirs      # Create dirs: certs & private.
make ca        # Create the CA certificate and key pair.
make node-cert # Create certificate and key pair for your nodes.
client-cert    # Create a client certificate and key pair for the root user.
```

## Step 2. [Start the cluster](https://www.cockroachlabs.com/docs/v21.1/secure-a-cluster#step-2-start-the-cluster)

Start the first node:

```sh
# make crdb-node1
* INFO: initial startup completed.
* Node will now attempt to join a running cluster, or wait for `cockroach init`.
* Client connections will be accepted after this completes successfully.
* Check the log file(s) for progress. 
```

Start two more nodes:

```sh
# make crdb-node2
* INFO: initial startup completed.
* Node will now attempt to join a running cluster, or wait for `cockroach init`.
* Client connections will be accepted after this completes successfully.
* Check the log file(s) for progress. 

# make crdb-node3
* INFO: initial startup completed.
* Node will now attempt to join a running cluster, or wait for `cockroach init`.
* Client connections will be accepted after this completes successfully.
* Check the log file(s) for progress. 
```

Perform a one-time initialization of the cluster:

```sh
# make cluster-init
Cluster successfully initialized
```

## Step 3. [Use the built-in SQL client](https://www.cockroachlabs.com/docs/v21.1/secure-a-cluster#step-3-use-the-built-in-sql-client)

You can use any node as a SQL gateway.

```sh
# make sql-client-node1
root@localhost:26257/defaultdb> CREATE DATABASE bank;
root@localhost:26257/defaultdb> CREATE TABLE bank.accounts (id INT PRIMARY KEY, balance DECIMAL);
root@localhost:26257/defaultdb> INSERT INTO bank.accounts VALUES (1, 1000.50);
root@localhost:26257/defaultdb> SELECT * FROM bank.accounts;
  id | balance
-----+----------
   1 | 1000.50
(1 row)
root@localhost:26257/defaultdb> \q                                    


# make sql-client-node2
root@localhost:26257/defaultdb> SELECT * FROM bank.accounts;
  id | balance
-----+----------
   1 | 1000.50
(1 row)

root@localhost:26257/defaultdb> CREATE USER max WITH PASSWORD 'roach';
root@localhost:26257/defaultdb> \q                                    
```


## Step 4. [Run a sample workload](https://www.cockroachlabs.com/docs/v21.1/secure-a-cluster#step-4-run-a-sample-workload)

Load the initial dataset for the movr workload.

```sh
# make movr-initial-dataset 
I211209 15:19:01.012148 1 workload/workloadsql/dataload.go:146  [-] 1  imported users (0s, 50 rows)
I211209 15:19:01.023156 1 workload/workloadsql/dataload.go:146  [-] 2  imported vehicles (0s, 15 rows)
I211209 15:19:01.060973 1 workload/workloadsql/dataload.go:146  [-] 3  imported rides (0s, 500 rows)
I211209 15:19:01.088180 1 workload/workloadsql/dataload.go:146  [-] 4  imported vehicle_location_histories (0s, 1000 rows)
I211209 15:19:01.117669 1 workload/workloadsql/dataload.go:146  [-] 5  imported promo_codes (0s, 1000 rows)
I211209 15:19:01.121802 1 workload/workloadsql/workloadsql.go:113  [-] 6  starting 8 splits
I211209 15:19:01.221583 1 workload/workloadsql/workloadsql.go:113  [-] 7  starting 8 splits
I211209 15:19:01.316522 1 workload/workloadsql/workloadsql.go:113  [-] 8  starting 8 splits
```

Run the workload for 5 minutes:

```sh
# make movr-run
...
  296.0s        0            2.0            3.0     10.5     12.6     12.6     12.6 addUser
  296.0s        0            0.0            1.0      0.0      0.0      0.0      0.0 addVehicle
  296.0s        0            0.0            1.0      0.0      0.0      0.0      0.0 applyPromoCode
  296.0s        0            0.0            0.3      0.0      0.0      0.0      0.0 createPromoCode
  296.0s        0            1.0            0.7     17.8     17.8     17.8     17.8 endRide
  296.0s        0           66.0          196.7      0.5      0.7      0.7      0.7 readVehicles
  296.0s        0            5.0            4.1     12.6     16.3     16.3     16.3 startRide
  296.0s        0            8.0           10.1    109.1    113.2    113.2    113.2 updateActiveRides
  297.0s        0            2.0            3.0     13.1     16.8     16.8     16.8 addUser
  297.0s        0            0.0            1.0      0.0      0.0      0.0      0.0 addVehicle
  297.0s        0            1.0            1.0     14.2     14.2     14.2     14.2 applyPromoCode
  297.0s        0            0.0            0.3      0.0      0.0      0.0      0.0 createPromoCode
  297.0s        0            2.0            0.7     10.5     10.5     10.5     10.5 endRide
  297.0s        0          169.0          196.6      0.5      0.6      0.7      0.8 readVehicles
  297.0s        0            3.0            4.1     13.1     14.2     14.2     14.2 startRide
  297.0s        0            8.0           10.1    104.9    117.4    117.4    117.4 updateActiveRides
  298.0s        0            4.0            3.0     14.7     18.9     18.9     18.9 addUser
  298.0s        0            1.0            1.0      8.1      8.1      8.1      8.1 addVehicle
  298.0s        0            0.0            1.0      0.0      0.0      0.0      0.0 applyPromoCode
  298.0s        0            0.0            0.3      0.0      0.0      0.0      0.0 createPromoCode
_elapsed___errors__ops/sec(inst)___ops/sec(cum)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)
  298.0s        0            1.0            0.7     10.0     10.0     10.0     10.0 endRide
  298.0s        0           70.0          196.2      0.5      0.7      0.7     22.0 readVehicles
  298.0s        0            2.0            4.1     14.7     18.9     18.9     18.9 startRide
  298.0s        0            8.0           10.1    109.1    121.6    121.6    121.6 updateActiveRides
  299.0s        0            1.0            3.0      9.4      9.4      9.4      9.4 addUser
  299.0s        0            0.0            1.0      0.0      0.0      0.0      0.0 addVehicle
  299.0s        0            2.0            1.0     10.0     14.2     14.2     14.2 applyPromoCode
  299.0s        0            0.0            0.3      0.0      0.0      0.0      0.0 createPromoCode
  299.0s        0            1.0            0.7     16.8     16.8     16.8     16.8 endRide
  299.0s        0          111.1          195.9      0.5      0.6      0.7      0.8 readVehicles
  299.0s        0            4.0            4.1     12.1     14.2     14.2     14.2 startRide
  299.0s        0            8.0           10.1    104.9    109.1    109.1    109.1 updateActiveRides
  300.0s        0            0.0            3.0      0.0      0.0      0.0      0.0 addUser
  300.0s        0            2.0            1.0     14.2     16.3     16.3     16.3 addVehicle
  300.0s        0            0.0            1.0      0.0      0.0      0.0      0.0 applyPromoCode
  300.0s        0            2.0            0.3      9.4     10.0     10.0     10.0 createPromoCode
  300.0s        0            1.0            0.7     14.2     14.2     14.2     14.2 endRide
  300.0s        0          256.9          196.1      0.5      0.6      0.8     22.0 readVehicles
  300.0s        0            2.0            4.1     10.0     14.7     14.7     14.7 startRide
  300.0s        0            7.0           10.1    100.7    104.9    104.9    104.9 updateActiveRides

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__total
  300.0s        0            894            3.0      8.5      9.4     17.8     22.0     24.1  addUser

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__total
  300.0s        0            290            1.0     11.0     11.0     19.9     22.0     24.1  addVehicle

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__total
  300.0s        0            306            1.0      9.9      9.4     18.9     23.1     26.2  applyPromoCode

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__total
  300.0s        0             97            0.3      8.1      8.9     17.8     22.0     26.2  createPromoCode

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__total
  300.0s        0            220            0.7      8.1      7.3     17.8     22.0     27.3  endRide

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__total
  300.0s        0          58833          196.1      0.5      0.5      0.6      0.9     71.3  readVehicles

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__total
  300.0s        0           1219            4.1     11.2     11.0     19.9     24.1     32.5  startRide

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__total
  300.0s        0           3026           10.1     80.1    100.7    151.0    176.2    260.0  updateActiveRides

_elapsed___errors_____ops(total)___ops/sec(cum)__avg(ms)__p50(ms)__p95(ms)__p99(ms)_pMax(ms)__result
  300.0s        0          64885          216.3      4.6      0.5     18.9    117.4    260.0  
```

## Step 5. [Access the DB Console](https://www.cockroachlabs.com/docs/v21.1/secure-a-cluster#step-5-access-the-db-console)

The CockroachDB DB Console gives you insight into the overall health of your cluster as well as the performance of the client workload.

On secure clusters, certain pages of the DB Console can only be accessed by admin users.

```sh
# make sql-client-node1
root@localhost:26257/defaultdb> GRANT admin TO max;
root@localhost:26257/defaultdb> \q
```

Go to https://localhost:8080 and login (max:roach)

On the [Cluster Overview](https://localhost:8080/#/overview/list), notice that three nodes are live, with an identical replica count on each node.

This demonstrates CockroachDB's [automated replication](https://www.cockroachlabs.com/docs/v21.1/demo-replication-and-rebalancing) of data via the Raft consensus protocol.

## Step 6. [Simulate node failure](https://www.cockroachlabs.com/docs/v21.1/secure-a-cluster#step-6-simulate-node-failure)

Run the cockroach quit command against a node to simulate a node failure:

```sh
# make crdb-node3-quit 
node is draining... remaining: 8
node is draining... remaining: 0 (complete)
ok
killing node3 process
Killed 319954
```

This demonstrates CockroachDB's use of the Raft consensus protocol to maintain availability and consistency in the face of failure; as long as a majority of replicas remain online, the cluster and client traffic continue uninterrupted.

## Step 7. [Scale the cluster](https://www.cockroachlabs.com/docs/v21.1/secure-a-cluster#step-7-scale-the-cluster)

Adding capacity is as simple as starting more nodes with `cockroach start`.

```sh
# make crdb-node4
* INFO: initial startup completed.
* Node will now attempt to join a running cluster, or wait for `cockroach init`.
* Client connections will be accepted after this completes successfully.
* Check the log file(s) for progress. 

# make crdb-node5
* INFO: initial startup completed.
* Node will now attempt to join a running cluster, or wait for `cockroach init`.
* Client connections will be accepted after this completes successfully.
* Check the log file(s) for progress. 
```

At first, the replica count will be lower for nodes 4 and 5. Very soon, however, you'll see those numbers even out across all nodes, indicating that data is being automatically rebalanced to utilize the additional capacity of the new nodes.

## Step 8. [Stop the cluster](https://www.cockroachlabs.com/docs/v21.1/secure-a-cluster#step-8-stop-the-cluster)

When you're done with your test cluster, use the `cockroach node drain` command to gracefully shut down each node.

```sh
# make cluster-stop
Terminating node1 process by sending SIGTERM
Terminated 325771
Terminating node2 process by sending SIGTERM
initiating graceful shutdown of server
Terminated 325811
initiating graceful shutdown of server
Terminating node3 process by sending SIGTERM
Terminated 325858
Terminating node4 process by sending SIGTERM
initiating graceful shutdown of server
Terminated 325897
Terminating node5 process by sending SIGTERM
initiating graceful shutdown of server
Terminated 325938
waiting 15 seconds ...
initiating graceful shutdown of server
server drained and shutdown completed
server drained and shutdown completed
Killing node1 process  by sending SIGKILL
killprocess.sh: line 4: kill: (325997) - No such process
Killing node2 process  by sending SIGKILL
killprocess.sh: line 4: kill: (326003) - No such process
Killing node3 process  by sending SIGKILL
Killed 325858
Killing node4 process  by sending SIGKILL
Killed 325897
Killing node5 process  by sending SIGKILL
Killed 325938
```