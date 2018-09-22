package cirno

import (
	"context"
	"log"
	"net"
	"time"

	_ "github.com/rs/xid"
)

const (
	Version         = "1.0.0"
	InfiniteTimeout = time.Duration(0)
)

type ListenFunc func(context.Context, string) error

type App struct {
	startedAt time.Time
	timeout   *time.Duration
}

func NewApp(timeout *time.Duration) (*App, error) {

	return &App{
		startedAt: time.Now(),
		timeout:   timeout,
	}, nil
}

// ListenSock starts to listen Unix Domain Socket on sockpath.
func (c *App) ListenSock(ctx context.Context, sockpath string) error {
	l, err := net.Listen("unix", sockpath)
	if err != nil {
		return err
	}

	return c.listen(ctx, l)
}

// ListenTCP starts to listen on addr "host:port".
func (c *App) ListenTCP(ctx context.Context, addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return c.listen(ctx, l)
}

func (c *App) listen(ctx context.Context, l net.Listener) error {
	go func() {
		<-ctx.Done()
		if err := l.Close(); err != nil {
			log.Println(err)
		}
	}()

	for {
		conn, err := l.Accept()
		if err != nil {
			select {
			case <-ctx.Done():
				log.Println("Shutting down listener")
				return nil
			default:
				return err
			}
		}

		go c.handleConn(ctx, conn)
	}
	return nil
}

func (c *App) handleConn(ctx context.Context, conn net.Conn) error {

	return nil
}
