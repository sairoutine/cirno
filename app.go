package cirno

import (
	"bufio"
	"context"
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strings"
	"time"
)

var (
	Version         = "1.0.0"
	InfiniteTimeout = time.Duration(0)
	respError       = []byte("ERROR\r\n")
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

func (c *App) handleConn(ctx context.Context, conn net.Conn) {
	defer func() {
		conn.Close()
	}()

	// set initial timeout
	err := c.extendDeadline(conn)
	if err != nil {
		log.Println(err)
		return
	}

	bufReader := bufio.NewReader(conn)

	scanner := bufio.NewScanner(bufReader)
	w := bufio.NewWriter(conn)

	for {
		// get client data
		if !scanner.Scan() {
			if err := scanner.Err(); err != nil {
				select {
				case <-ctx.Done():
					// shutting down
					return
				default:
					if ne, ok := err.(net.Error); ok {
						switch {
						case ne.Timeout():
							// client timeout
							return
						case ne.Temporary():
							// try to scan again
							continue
						default:
							log.Println(err)
							return
						}
					} else {
						log.Println(err)
						return
					}
				}
			} else {
				return
			}
		}

		// update timeout
		c.extendDeadline(conn)
		if err != nil {
			log.Println(err)
			return
		}

		// parse command
		cmd, err := bytesToCommand(scanner.Bytes())
		if err != nil {
			if err := c.writeError(conn); err != nil {
				// can't return error message to client
				log.Println(err)
				return
			} else {
				// error has occured and succeeded to return error to client,
				// nothing is to do.
				continue
			}
		}

		// execute command
		if err := cmd.Execute(c, w); err != nil {
			if err == io.EOF {
				// normally finish
				return
			} else if err := c.writeError(conn); err != nil {
				// error has occured and can't return error message to client
				log.Println(err)
				return
			} else {
				// error has occured and succeeded to return error to client,
				// nothing is to do.
			}
		}

		// return response
		if err := w.Flush(); err != nil {
			if err == io.EOF {
				// normally finish
				return
			} else {
				log.Println("error on cmd %s write to conn: %s", cmd, err)
				return
			}
		}
	}
}

func (c *App) extendDeadline(conn net.Conn) error {
	if *c.timeout == InfiniteTimeout {
		return nil
	}
	d := time.Now().Add(*c.timeout)
	return conn.SetDeadline(d)
}

func (c *App) writeError(conn io.Writer) (err error) {
	_, err = conn.Write(respError)
	return
}

// bytesToCommand converts byte array to a MemdCommand and returns it.
func bytesToCommand(data []byte) (MemdCommand, error) {
	if len(data) == 0 {
		return nil, errors.New("No command")
	}

	fields := strings.Fields(string(data))
	switch name := strings.ToUpper(fields[0]); name {
	case "GET", "GETS":
		if len(fields) < 2 {
			return nil, errors.New("GET command needs key as second parameter")
		}
		return &MemdCommandGet{
			Name: name,
			Keys: fields[1:],
		}, nil
	case "QUIT":
		return MemdCommandQuit(0), nil
	case "VERSION":
		return MemdCommandVersion(1), nil
	default:
		return nil, errors.New(fmt.Sprintf("Unknown command: %s", name))
	}
}
