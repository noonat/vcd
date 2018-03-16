package main

import (
	"context"
	"log"

	"github.com/namsral/flag"
	"github.com/noonat/vcd"
)

func main() {
	var (
		listenAddr string
		mysqlDSN   string
	)
	flag.StringVar(&listenAddr, "listen-addr", ":8080", "HTTP listen address")
	flag.StringVar(&mysqlDSN, "mysql-dsn", "", "MySQL connection string")
	flag.Parse()

	ctx := context.Background()

	if err := vcd.Run(ctx, listenAddr, mysqlDSN); err != nil {
		log.Fatalf("%+v", err)
	}
}
