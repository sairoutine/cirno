package cirno

import (
	"context"
	"time"
)

/*
import (
	"github.com/rs/xid"
)
*/

const (
	Version         = "1.0.0"
	InfiniteTimeout = time.Duration(0)
)

type ListenFunc func(context.Context, string) error

type App struct {
	timeout *time.Duration
}

func NewApp(timeout *time.Duration) (*App, error) {

	return &App{
		timeout: timeout,
	}, nil
}

func (c *App) ListenSock(context.Context, string) error {

	return nil
}

func (c *App) ListenTCP(context.Context, string) error {

	return nil
}
