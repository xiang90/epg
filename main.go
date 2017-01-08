package main

import (
	"context"
	"database/sql"
	"log"

	"github.com/coreos/etcd/clientv3"
	"github.com/coreos/etcd/clientv3/mirror"
	_ "github.com/lib/pq"
	"github.com/xiang90/edb"
)

func main() {
	db, err := sql.Open("postgres", "")

	if err != nil {
		log.Fatal(err)
	}

	etcd, err := clientv3.New(clientv3.Config{
		Endpoints: []string{"localhost:2379"},
	})
	if err != nil {
		log.Fatal(err)
	}
	syncer := mirror.NewSyncer(etcd, "", 0)

	dbs := &edb.SQLDBSyncer{
		Table: "kvs",

		Syncer: syncer,
		DB:     db,
	}

	log.Fatal(dbs.Sync(context.Background()))
}
