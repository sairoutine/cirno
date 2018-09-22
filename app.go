package cirno

import (
	"context"
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

	return c.Listen(ctx, l)
}

// ListenTCP starts to listen on addr "host:port".
func (c *App) ListenTCP(ctx context.Context, addr string) error {
	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}

	return c.Listen(ctx, l)
}

func (c *App) Listen(ctx context.Context, l net.Listener) error {

	return nil
}
