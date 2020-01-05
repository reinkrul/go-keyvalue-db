package server

import (
    "errors"
    "fmt"
    "github.com/reinkrul/go-keyvalue-db/server/bolt"
    "github.com/reinkrul/go-keyvalue-db/server/grpc"
    "github.com/reinkrul/go-keyvalue-db/server/spi"
    "github.com/urfave/cli/v2"
    "log"
    "os"
    "os/signal"
    "syscall"
)

func Command() *cli.Command {
    return &cli.Command{
        Name:  "server",
        Usage: "server mode",
        Action: func(context *cli.Context) error {
            return start(":4446")
        },
    }

}

func start(addr string) error {
    log.Printf("Starting server on %s", addr)
    database, err := startDatabase()
    if err != nil {
        return err
    }
    api, err := startAPI(addr, database)
    if err != nil {
        return err
    }
    log.Print("Server started.")
    closeFn := closeServices(api, database)

    signalChannel := make(chan os.Signal, 1)
    signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
    for {
        <- signalChannel
        closeFn()
        break
    }
    log.Print("Bye bye.")
    return nil
}

func startDatabase() (spi.DataStore, error) {
    database, err := bolt.Connect("data.db")
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Unable to connect to database: %v", err))
    }
    return database, nil
}

func startAPI(apiAddr string, store spi.DataStore) (spi.API, error) {
    api, err := grpc.Start(apiAddr, store)
    if err != nil {
        return nil, errors.New(fmt.Sprintf("Unable to start API: %v", err))
    }
    return api, nil
}

func closeServices(services ...spi.Service) func() {
    return func() {
        for _, service := range services {
            err := service.Close()
            if err != nil {
                log.Printf("Unable to close %T: %v", service, err)
            }
        }
    }
}