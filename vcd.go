// Package vcd provides an HTTP server which hosts the Vehicle Configuration
// Database for Captain Forever. Captain Forever is a space arcade game, where
// you blow up enemy ships and use their component parts to improve the design
// of your own ship. The game allows you to export your ship to various
// formats, and VCD allows you to easily share the HTML export of your ship
// with others.
package vcd

import (
	"context"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	"github.com/pkg/errors"
)

// Run starts the VCD app.
func Run(ctx context.Context, listenAddr, mysqlDSN string) error {
	db, err := openDB(ctx, mysqlDSN)
	if err != nil {
		return err
	}

	router := httprouter.New()
	router.GET("/", handleList(db))
	router.GET("/new", handleNew(db))
	router.GET("/vessels/:id", handleVessel(db))
	router.GET("/vessels/:id/pilot", handleVesselPilot(db))
	router.NotFound = http.FileServer(http.Dir("./public"))

	log.Printf("Starting server on %s", listenAddr)
	if err := http.ListenAndServe(listenAddr, router); err != nil {
		return errors.Wrap(err, "error listening")
	}

	return nil
}
