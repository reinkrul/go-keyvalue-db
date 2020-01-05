package client

import (
	"context"
	"fmt"
	"github.com/reinkrul/go-keyvalue-db/api"
	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"log"
	"time"
)

func Command() *cli.Command {
	addrFlag := cli.StringFlag{Name: "addr", Usage: "database address to connect to", Required: true}
	bucketFlag := cli.StringFlag{Name: "bucket", Usage: "bucket to operate on", Required: true}
	keyFlag := cli.StringFlag{Name: "key", Usage: "key to operate on", Required: true}
	valueFlag := cli.StringFlag{Name: "value", Usage: "value to write", Required: true}
	return &cli.Command{
		Name:  "client",
		Usage: "CLI client mode",
		Subcommands: []*cli.Command{
			{
				Name:  "set",
				Usage: "write a value",
				Flags: []cli.Flag{&addrFlag, &bucketFlag, &keyFlag, &valueFlag},
				Action: func(c *cli.Context) error {
					return execute(c.String("addr"), func(ctx context.Context, client api.KeyValueDatabaseClient) error {
						bucket := c.String("bucket")
						key := c.String("key")
						fmt.Println(fmt.Sprintf("Writing value to %s/%s", bucket, key))
						_, err := client.Set(ctx, &api.SetValueRequest{
							Bucket: bucket,
							Key:    key,
							Value:  c.String("value"),
						})
						return err
					})
				},
			},
			{
				Name:  "get",
				Usage: "read a value",
				Flags: []cli.Flag{&addrFlag, &bucketFlag, &keyFlag},
				Action: func(c *cli.Context) error {
					return execute(c.String("addr"), func(ctx context.Context, client api.KeyValueDatabaseClient) error {
						response, err := client.Get(ctx, &api.GetValueRequest{
							Bucket: c.String("bucket"),
							Key:    c.String("key"),
						})
						if response != nil {
							fmt.Println(response.Value)
						}
						return err
					})
				},
			},
		},
	}
}

func execute(addr string, consumer func(context.Context, api.KeyValueDatabaseClient) error) error {
	conn, err := grpc.Dial(addr, grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := api.NewKeyValueDatabaseClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	return consumer(ctx, c)
}
