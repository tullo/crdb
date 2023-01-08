SHELL := /bin/bash -eo pipefail

.DEFAULT_GOAL := cluster

clean: cluster-stop
	@rm -rf certs node1 node2 node3 node4 node5 private

dirs:
	@mkdir -p certs private

ca: dirs
	@cockroach cert create-ca --certs-dir=certs --ca-key=private/ca.key

node-cert:
	@cockroach cert create-node ::1 localhost $(shell hostname) --certs-dir=certs --ca-key=private/ca.key

client-cert:
	@cockroach cert create-client root --certs-dir=certs --ca-key=private/ca.key

certs: ca node-cert client-cert

CLUSTER = node1 node2 node3 node4 node5

cluster: certs crdb-node1 crdb-node2 crdb-node3 cluster-init
	@sleep 3
	@make cluster-logs

cluster-scale: crdb-node4 crdb-node5
	@sleep 3
	@cockroach node status --certs-dir=certs

cluster-status:
	cockroach node status --certs-dir=certs --host=$(shell hostname):26257

cluster-stop:
	@for n in $(CLUSTER); do bash terminateprocess.sh $$n; done
	@echo "waiting 15 seconds ..."
	@sleep 15
	@for n in $(CLUSTER); do bash killprocess.sh $$n; done

crdb-node1:
	@cockroach start --certs-dir=certs --store=node1 \
	--advertise-addr=$(shell hostname):26258 \
	--sql-addr=:26257 --listen-addr=:26258 \
	--join=$(shell hostname):26258,$(shell hostname):26260,$(shell hostname):26262 \
	--http-addr=:8080 \
	--cache=.25 \
	--max-sql-memory=.25 \
	--background

crdb-node2:
	@cockroach start --certs-dir=certs --store=node2 \
	--advertise-addr=$(shell hostname):26260 \
	--sql-addr=:26259 --listen-addr=:26260 \
	--join=$(shell hostname):26258,$(shell hostname):26260,$(shell hostname):26262 \
	--http-addr=:8081
	--cache=.25 \
	--max-sql-memory=.25 \
	--background

crdb-node3:
	@cockroach start --certs-dir=certs --store=node3 \
	--advertise-addr=$(shell hostname):26262 \
	--sql-addr=:26261 --listen-addr=:26262 \
	--join=$(shell hostname):26258,$(shell hostname):26260,$(shell hostname):26262 \
	--http-addr=:8082
	--cache=.25 \
	--max-sql-memory=.25 \
	--background

crdb-node3-quit:
	@cockroach node drain 3 --certs-dir=certs --host=:26262
	@bash killprocess.sh node3

crdb-node4:
	@cockroach start --certs-dir=certs --store=node4 \
	--advertise-addr=$(shell hostname):26264 \
	--sql-addr=:26263 --listen-addr=:26264 \
	--http-addr=localhost:8083 \
	--join=$(shell hostname):26258,$(shell hostname):26260,$(shell hostname):26262 \
	--http-addr=:8083
	--cache=.25 \
	--max-sql-memory=.25 \
	--background

crdb-node5:
	@cockroach start --certs-dir=certs --store=node5 \
	--advertise-addr=$(shell hostname):26266 \
	--sql-addr=:26265 --listen-addr=:26266 \
	--http-addr=localhost:8084 \
	--join=$(shell hostname):26258,$(shell hostname):26260,$(shell hostname):26262 \
	--http-addr=:8084
	--cache=.25 \
	--max-sql-memory=.25 \
	--background

cluster-init:
	cockroach init --certs-dir=certs --host=$(shell hostname):26258

cluster-logs:
	@grep 'node starting' node1/logs/cockroach.log -A 11

sql-client-node1:
	@cockroach sql --certs-dir=certs --host=$(shell hostname):26257

sql-client-node2:
	@cockroach sql --certs-dir=certs --host=$(shell hostname):26259

movr-initial-dataset:
	@cockroach workload init movr \
	'postgresql://root@$(shell hostname):26257?sslcert=certs%2Fclient.root.crt&sslkey=certs%2Fclient.root.key&sslmode=verify-full&sslrootcert=certs%2Fca.crt'

movr-run: movr-initial-dataset
	@cockroach workload run movr \
	--duration=5m 'postgresql://root@$(shell hostname):26257?sslcert=certs%2Fclient.root.crt&sslkey=certs%2Fclient.root.key&sslmode=verify-full&sslrootcert=certs%2Fca.crt'

.PHONY : crdb-node1 crdb-node2 crdb-node3 crdb-node4 crdb-node5 cluster-init sql-client-node2 movr-initial-dataset movr-run crdb-node3-quit cluster-stop
