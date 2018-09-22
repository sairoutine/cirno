package cirno

import (
	"io"
	"strconv"

	"github.com/rs/xid"
)

var (
	memdValHeader     = []byte("VALUE ")
	memdVersionHeader = []byte("VERSION ")
	memdSep           = []byte("\r\n")
	memdSpc           = []byte(" ")
)

type MemdCommand interface {
	Execute(*App, io.Writer) error
}

type MemdCommandGet struct {
	Name string
	Keys []string
}

// Execute generates new ID.
func (cmd *MemdCommandGet) Execute(app *App, w io.Writer) error {
	values := make([]string, len(cmd.Keys))
	for i := range cmd.Keys {
		values[i] = xid.New().String()
	}

	// WriteTo writes content of MemdValue to io.Writer.
	// Its format is compatible to memcached protocol.
	for i, key := range cmd.Keys {
		w.Write(memdValHeader)
		io.WriteString(w, key)
		w.Write(memdSpc)
		io.WriteString(w, strconv.Itoa(0)) // flag
		w.Write(memdSpc)
		io.WriteString(w, strconv.Itoa(len(values[i])))
		w.Write(memdSep)
		io.WriteString(w, values[i])
		w.Write(memdSep)
	}

	return nil
}

// MemdCommandQuit defines QUIT command.
type MemdCommandQuit int

// Execute disconnect by server.
func (cmd MemdCommandQuit) Execute(app *App, conn io.Writer) error {
	return io.EOF
}

// MemdCommandVersion defines VERSION command.
type MemdCommandVersion int

// Execute writes Version number.
func (cmd MemdCommandVersion) Execute(app *App, w io.Writer) error {
	w.Write(memdVersionHeader)
	io.WriteString(w, Version)
	w.Write(memdSep)
	return nil
}
