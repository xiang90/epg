## Why

A lot of etcd users want to have advanced query over the data they stored into etcd. But
etcd does not support secondary index currently (We do have plan to support it eventually though).

With GIN indexes support in Postgres, it is easier to build secondary indexes for JSONB or unstructured data. Why not query etcd data by utilizing the powerful SQL and GIN indexes?

etcd is the ONLY consistent data store that supports streaming watch (as of 1.7.2017), which allows us to sync data into other data sink consistently and efficiently.

epg syncs etcd data into Postgres consistently. Users can use Postgres to query the data synced in.

The Postgres server can run on tmpfs or any other violate storages. It can be stateless. All data is
stored in etcd, and will be synced back by `epg` when you lose Postgres data.

## DEMO

Setup etcd

```bash
./etcd
```

Setup pg

Start a local postgres server

```sql
CREATE DATABASE etcd;
\c etcd
CREATE TABLE kvs (key bytea, value bytea, primary key (key) );
```

Start epg

```
./PGUSER=$USERNAME PGPASSWORD=$PASSWORD PGSSLMODE=disable PGDATABASE=etcd ./epg
```

Write some keys into etcd

```bash
$ etcdctl put foo1 bar1
OK
$ etcdctl put foo2 bar2
OK
```

Select the table on pg

```sql

SELECT * FROM kvs;

    key     |   value
------------+------------
 \x666f6f31 | \x62617231
 \x666f6f32 | \x62617232
(2 rows)

```

Enjoy!
