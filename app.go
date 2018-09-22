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

const Version = "1.0.0"

type ListenFunc func(context.Context, string) error

type App struct {
	idleTimeout time.Duration
}

func NewApp(timeout time.Duration) *App {

	return &App{
		idleTimeout: timeout,
	}
}
