package main

import (
	"fmt"
	"log"
	"os"
	"strconv"

	"github.com/jackc/pgx"
)

func connect() (*pgx.ConnPool, error) {
	connConfig, err := pgx.ParseEnvLibpq()

	connConfig.PreferSimpleProtocol = true // enabled for a postgres-pooler-odyssey to work

	if err != nil {
		return nil, err
	}

	poolCfg := pgx.ConnPoolConfig{
		ConnConfig:     connConfig,
		AfterConnect:   nil,
		AcquireTimeout: 15 * 1000 * 1000 * 1000,
	}

	if max := os.Getenv("PGMAXCONNECTIONS"); max != "" {
		if maxInt, err := strconv.Atoi(max); err == nil {
			poolCfg.MaxConnections = maxInt
		} else {
			return nil, err
		}
	}

	conn, err := pgx.NewConnPool(poolCfg)

	if err != nil {
		return nil, err
	}

	return conn, nil
}

func onStartUpErr(err error) {
	panic(fmt.Sprintf("Failure on start up: %s", err))
}

func main() {
	log.Printf("Hello!")

	pgPool, err := connect()

	if err != nil {
		onStartUpErr(err)
	}

	defer pgPool.Close()

	log.Printf("--- Migration ---")

	// pgConn, err := pgPool.Acquire()
	// if err != nil {
	// 	onStartUpErr(err)
	// }
	// defer pgPool.Release(pgConn)

	// if err := startMigration(pgConn); err != nil {
	if err := startMigration(); err != nil {
		onStartUpErr(err)
	}
	log.Printf("--- End of Migration ---")

	// time.Sleep(1 * time.Minute)

	rows, err := pgPool.Query("SELECT COUNT(1) FROM information_schema.tables WHERE table_name = $1 AND table_schema = (SELECT current_schema()) LIMIT 1", "t")
	if err != nil {
		panic(err)
	}

	// rows.Close is called by rows.Next when all rows are read
	// or an error occurs in Next or Scan. So it may optionally be
	// omitted if nothing in the rows.Next loop can panic. It is
	// safe to close rows multiple times.
	defer rows.Close()

	var sum int32

	// Iterate through the result set
	for rows.Next() {
		var n int32
		err = rows.Scan(&n)
		if err != nil {
			panic(err)
		}
		sum += n
		log.Print(sum)
	}

	// Any errors encountered by rows.Next or rows.Scan will be returned here
	if rows.Err() != nil {
		panic(err)
	}

	log.Printf("Bye!")
}
