package main

import (
	"context"
	"database/sql"
	"fmt"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/mirror"
)

type dbSyncer struct {
	table string

	src  mirror.Syncer
	dest *sql.DB
}

func (ds *dbSyncer) sync(ctx context.Context) error {
	stmt, err := ds.dest.Prepare(fmt.Sprintf("INSERT INTO %s(key,value) VALUES ($1,$2) ON CONFLICT (key) DO UPDATE SET value=$2", ds.table))
	if err != nil {
		return err
	}

	dstmt, err := ds.dest.Prepare(fmt.Sprintf("DELETE FROM %s WHERE value=$1", ds.table))
	if err != nil {
		return err
	}

	rc, errc := ds.src.SyncBase(ctx)
	select {
	case r := <-rc:
		for _, kv := range r.Kvs {
			if _, err := stmt.Exec(kv.Key, kv.Value); err != nil {
				return err
			}
		}
	case e := <-errc:
		return e
	}

	wch := ds.src.SyncUpdates(ctx)
	for wr := range wch {
		if wr.Err() != nil {
			return wr.Err()
		}

		for _, e := range wr.Events {
			switch e.Type {
			case clientv3.EventTypePut:
				if _, err := stmt.Exec(e.Kv.Key, e.Kv.Value); err != nil {
					return err
				}
			case clientv3.EventTypeDelete:
				if _, err := dstmt.Exec(e.Kv.Key); err != nil {
					return err
				}
			}
		}
	}

	return ctx.Err()
}
