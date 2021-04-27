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
