# Cockroach Go examples

- [Bank](https://github.com/tullo/examples-go/tree/master/bank)
  - Continuously performs balance transfers between accounts using concurrent transactions.
- [Bank)](https://github.com/tullo/examples-go/tree/master/bank2) (high and low contention)
  - Transfers money between accounts, creating new ledger transactions in the form of a transaction record and two transaction "legs" per database transaction.
  - Each transfer additionally queries and updates account balances.
- [Block writer](https://github.com/tullo/examples-go/tree/master/block_writer) (split and rebalance)
  - The program is a write-only workload intended to insert a large amount of data into cockroach quickly.
  - This example is intended to trigger range splits and rebalances.
- [Real Time](https://github.com/tullo/examples-go/tree/master/fakerealtime)
  - Uses a log-style table in an approximation of the "fake real time" system used at Friendfeed a Facebook acquisition.
- [Filesystem](https://github.com/tullo/examples-go/tree/master/filesystem) (Filesystem in Userspace) 
  - A fuse filesystem using cockroach as a backing store.
- [Hotspot](https://github.com/tullo/examples-go/tree/master/hotspot)
  - A read/write workload intended to always hit the exact same value.
  - It performs reads and writes to simulate a super contentious load.
- [Ledger](https://github.com/tullo/examples-go/tree/master/ledger) (causing complete deadlock in the more contended modes)
  - Simulate a ledger and a certain type of workload against it.
  - A general ledger is a complete record of financial transactions over the life of a bank (or other company).
  - Aims to model a bank in a more realistic setting and tickles contention issues.
- [Photos](https://github.com/tullo/examples-go/tree/master/photos) (Benchmark)
  - Create artificial load using a simple database schema containing
users, photos and comments.
  - Users have photos, photos have comments.
  - Users can author comments on any photos.
  - User actions are simulated
using an exponential distribution on user IDs, so lower IDs see
more activity than high ones.

## wikifeedia

Sample application highlighting the benefits of using CockroachDB, specifically `Follower Reads`.

https://github.com/tullo/wikifeedia (Go & React)