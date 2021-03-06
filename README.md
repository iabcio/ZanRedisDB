# ZanRedisDB

[![Codacy Badge](https://api.codacy.com/project/badge/Grade/5bb2847636f343e79edf048a0394de04)](https://www.codacy.com/app/cool8511/youzan_ZanRedisDB?utm_source=github.com&amp;utm_medium=referral&amp;utm_content=youzan/ZanRedisDB&amp;utm_campaign=Badge_Grade)
![Go](https://github.com/youzan/ZanRedisDB/workflows/Go/badge.svg)[![Build Status](https://travis-ci.com/youzan/ZanRedisDB.svg?branch=master)](https://travis-ci.com/youzan/ZanRedisDB) [![GitHub release](https://img.shields.io/github/release/youzan/ZanRedisDB.svg)](https://github.com/youzan/ZanRedisDB/releases/latest) [![codecov](https://codecov.io/gh/youzan/ZanRedisDB/branch/master/graph/badge.svg)](https://codecov.io/gh/youzan/ZanRedisDB) [![Go Report Card](https://goreportcard.com/badge/github.com/youzan/ZanRedisDB)](https://goreportcard.com/report/github.com/youzan/ZanRedisDB) [![Documentation Status](https://readthedocs.org/projects/youzan-zanredisdb/badge/?version=latest)](http://youzan-zanredisdb.readthedocs.io/en/latest/?badge=latest)
 

## What is ZanRedisDB
ZanRedisDB is a distributed redis cluster with strong consistency.
- Compatible with redis protocol: most redis command is supported in ZanRedisDB, so you can replace Redis with ZanRedisDB and get a more powerful cluster with unlimited data.
- Namespace and set supported: You can use the namespace and set as prefix to isolate different data to make operation more easily. 
- Strong consistence and persistence: Use raft and rocksdb to make sure all the data have the consistence replication and persistent on stable storage.
- Scalability: You can easily increase the cluster capacity just by adding more machines.
- Multi-DC support: multi data centers deployment can be supported with rackaware feature.

## Build

Install the compress library
<pre>
yum install snappy-devel (for CentOS)
apt-get install libsnappy1 libsnappy-dev (for Debian/Ubuntu)
brew install snappy (for Mac)
</pre>

Build the rocksdb with jemalloc
<pre>
git clone https://github.com/absolute8511/rocksdb.git
cd rocksdb
git checkout v6.4.6-patched
PORTABLE=1 WITH_JEMALLOC_FLAG=1 JEMALLOC=1 make static_lib
</pre>

Install the dependency (for old go version only, if using go1.13+, it will be done in go modules):
<pre>
CGO_CFLAGS="-I/path/to/rocksdb/include" CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lsnappy -lrt -ljemalloc" go get github.com/youzan/gorocksdb

CGO_CFLAGS="-I/path/to/rocksdb/include" CGO_LDFLAGS="-L/path/to/rocksdb -lrocksdb -lstdc++ -lm -lsnappy -ljemalloc" go get github.com/youzan/gorocksdb (for MacOS)
</pre>

use the `dep ensure` to install other dependencies or use go modules for go1.13+

Build zankv and placedriver from the source (only support go version 1.10+, gcc 4.9+ or xcode-command-line-tools on Mac):
<pre>
ROCKSDB=/path/to/rocksdb make
</pre>

If you want package the binary release run the scripts
<pre>
./pre-dist.sh
ROCKSDB=/path/to/rocksdb ./dist.sh
</pre>

## Deploy

 * Deploy the rsync daemon which is needed on all server node to transfer the snapshot data for raft
 * Deploy etcd cluster which is needed for the meta data for the namespaces
 * Deploy the placedriver which is used for data placement: `placedriver -config=/path/to/config`
 * Deploy the zankv for data storage server `zankv -config=/path/to/config`
 * Init a namespace using the create the namespace API in placedriver

## OS-Level Tuning

 * Setting `vm.swappiness=0`
 * Setting `vm.min_free_kbytes` to at least 1GB (8GB on larger memory system)
 * Disable NUMA zone reclaim with `vm.zone_reclaim_mode=0`
 * Disable THP(transparent huge pages)
 * Avoid the tcp delay ack by `echo 4 > /proc/sys/net/ipv4/tcp_delack_min` (for old OS only)


## API
placedriver has several HTTP APIs to manager the namespace
 * list the namespace: `GET /namespaces`
 * list the data nodes: `GET /datanodes`
 * list the placedriver nodes: `GET /listpd`
 * query the namespace meta info: `GET /query/namespace_name`
 * create the namespace (handle only by leader) : `POST /cluster/namespace/create?namespace=test_p16&partition_num=16&replicator=3`
 * delete the namespace (handle only by leader): `POST /cluster/namespace/delete?namespace=test_p16&partition=**`

storage server HTTP APIs for stats:
 * namespace stats : `GET /stats`
 * namespace raft internal stats : `GET /raft/stats`
 * optimize the data storage : `POST /kv/optimize`
 * get the raft leader of the namespace partition: `GET /cluster/leader/namespace-partition`

storage server also support the redis apis for read/write :
 * KV:
 * Hash Set:
 * List:
 * Sorted Set:
 * ZSet:

## Client
Golang client SDK : [client-sdk] , a redis proxy can be deployed 
based on this golang sdk if you want use the redis client in other language.

## Architechture

![arch](doc/resource/zankv-arch.png)
## Roadmap
* Redis data structures
  - [x] KV
  - [x] Hash
  - [x] List
  - [x] Set
  - [x] Sorted Set
  - [x] GeoHash
  - [x] Expires
  - [x] HyperLogLog
  - [x] JSON
* Distributed system
  - [x] Raft based replication
  - [x] Partitions
  - [x] Auto balance and migrate
  - [x] Support namespace
  - [x] High available
  - [x] Distributed scan on table
  - [x] Rackaware deployment used for across data center deployment
* Searchable and Indexing
  - [ ] Secondary index support on Hash fields
  - [ ] Secondary index support for json kv
  - [ ] Full text search support
* Operation
  - [x] Backup and restore for cluster
  - [x] More stats for read/write performance and errors.
* Client 
  - [x] High available for redis commands (Retry on fail)
  - [ ] Extand redis commands to support index and search
  - [x] Extand redis commands for advance scan
* Others (maybe)
  - [ ] Support configure for Column storage friendly for OLAP
  - [ ] Pebble as the storage engine
  - [ ] Lua scripts support
  - [ ] Support export data to other systems


[client-sdk]: https://github.com/youzan/go-zanredisdb

## Thanks

Many thanks for these great projects which make this project possible: etcd, RocksDB, ledisdb, pika.

