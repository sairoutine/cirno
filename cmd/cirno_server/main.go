package main

import (
	"context"
	"log"
	"os"
	"sync"

	"github.com/sairoutine/cirno"
	"github.com/urfave/cli"
)

type cirnoConfig struct {
	port     uint
	timeout  uint
	sockpath string
}

func main() {
	app := cli.NewApp()
	app.Name = "cirno server"
	app.Usage = "ID generator server by xid algorithms."
	app.Version = "1.0.0"

	// command's option
	app.Flags = []cli.Flag{
		cli.UintFlag{
			Name:  "port, p",
			Value: 11212,
			Usage: "port to listen",
		},
		cli.StringFlag{
			Name:  "sock, s",
			Value: "",
			Usage: "unix domain socket to listen. ignore port option when set this",
		},
		cli.UintFlag{
			Name:  "timeout, t",
			Value: 5,
			Usage: "connection will be closed if there are no packets over the seconds. 0 means infinite",
		},
	}

	app.Action = func(c *cli.Context) error {
		var wg sync.WaitGroup
		ctx, cancel := context.WithCancel(context.Background())

		wg.Add(1)
		go signalHandler(ctx, cancel, &wg)

		// main listener
		fn, addr, err := listenFunc(cirnoConfig{
			port:     c.Uint("port"),
			sockpath: c.String("sock"),
			timeout:  c.Uint("timeout"),
		})
		if err != nil {
			return err
		}

		wg.Add(1)
		go mainListener(ctx, &wg, fn, addr)

		wg.Wait()
		log.Println("Shutdown completed.")

		return nil
	}

	// execute
	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func signalHandler(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {

}

func listenFunc(conf cirnoConfig) (cirno.ListenFunc, string, error) {

	return nil, "", nil
}

func mainListener(ctx context.Context, wg *sync.WaitGroup, fn cirno.ListenFunc, addr string) {

}
