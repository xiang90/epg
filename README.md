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
