SHELL := /bin/bash -eo pipefail

.DEFAULT_GOAL := insecure-node

clean:
	@rm -rfv cockroach-data certs private

secure-node: clean
	@mkdir -pv certs private
	@cockroach cert create-ca --certs-dir=certs --ca-key=private/ca.key
	@cockroach cert create-node localhost $(shell hostname) --certs-dir=certs --ca-key=private/ca.key
	@cockroach cert create-client root --certs-dir=certs --ca-key=private/ca.key
	@cockroach cert create-client maxroach --certs-dir=certs --ca-key=private/ca.key
	@cockroach start-single-node --certs-dir=certs --advertise-addr=$(shell hostname)

insecure-node: clean
	@cockroach start-single-node --insecure --listen-addr=localhost

insecure-sql-insert:
	cockroach sql --insecure -e 'CREATE USER IF NOT EXISTS maxroach; CREATE DATABASE IF NOT EXISTS bank; GRANT ALL ON DATABASE bank TO maxroach;'
	cockroach sql --insecure --database bank -e 'CREATE TABLE IF NOT EXISTS accounts (ID SERIAL PRIMARY KEY, balance INT);'

secure-sql-insert:
	cockroach sql --certs-dir=certs -e 'CREATE USER IF NOT EXISTS maxroach; CREATE DATABASE IF NOT EXISTS bank; GRANT ALL ON DATABASE bank TO maxroach;'
	cockroach sql --certs-dir=certs --database bank -e 'CREATE TABLE IF NOT EXISTS accounts (ID SERIAL PRIMARY KEY, balance INT);'
