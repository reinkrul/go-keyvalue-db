package main

import (
	"github.com/reinkrul/go-keyvalue-db/client"
	"github.com/reinkrul/go-keyvalue-db/server"
	"github.com/urfave/cli/v2"
	"log"
	"os"
)

func main() {
	app := &cli.App{
		Name:     "keyvaluedb",
		Usage:    "either start the server or use client mode to connect to a running instance",
		Commands: []*cli.Command{client.Command(), server.Command()},
	}
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}