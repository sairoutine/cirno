package cirno

import "time"

/*
import (
	"github.com/rs/xid"
)
*/

type App struct {
	idleTimeout time.Duration
}

func NewApp(timeout time.Duration) *App {

	return &App{
		idleTimeout: timeout,
	}
}
