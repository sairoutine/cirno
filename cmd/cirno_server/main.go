package main

import (
	"context"
	"fmt"
	"log"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"

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

		listenFunc, addr, err := setupListenFunc(cirnoConfig{
			port:     c.Uint("port"),
			sockpath: c.String("sock"),
			timeout:  c.Uint("timeout"),
		})
		if err != nil {
			return err
		}

		wg.Add(1)
		go mainListener(ctx, &wg, listenFunc, addr)

		wg.Wait()
		log.Println("Shutdown completed.")

		return nil
	}

	err := app.Run(os.Args)
	if err != nil {
		log.Fatal(err)
	}
}

func signalHandler(ctx context.Context, cancel context.CancelFunc, wg *sync.WaitGroup) {
	defer wg.Done()
	trapSignals := []os.Signal{
		syscall.SIGHUP,
		syscall.SIGINT,
		syscall.SIGTERM,
		syscall.SIGQUIT,
	}
	sigCh := make(chan os.Signal, 1)
	signal.Notify(sigCh, trapSignals...)
	select {
	case sig := <-sigCh:
		log.Printf("received signal %s", sig)
		cancel()
	case <-ctx.Done():
	}
}

func setupListenFunc(conf cirnoConfig) (cirno.ListenFunc, string, error) {
	var timeout time.Duration
	if conf.timeout == 0 {
		timeout = cirno.InfiniteTimeout
	} else {
		timeout = time.Duration(conf.timeout) * time.Second
	}

	app, err := cirno.NewApp(&timeout)
	if err != nil {
		return nil, "", err
	}

	if conf.sockpath != "" {
		return app.ListenSock, conf.sockpath, nil
	} else {
		return app.ListenTCP, fmt.Sprintf(":%d", conf.port), nil
	}
}

func mainListener(ctx context.Context, wg *sync.WaitGroup, listenFunc cirno.ListenFunc, addr string) {
	defer wg.Done()
	if err := listenFunc(ctx, addr); err != nil {
		log.Println("Listen failed", err)
		os.Exit(1)
	}
}
