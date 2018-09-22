package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"syscall"
	"testing"
	"time"

	"github.com/bradfitz/gomemcache/memcache"
	"github.com/lestrrat-go/tcptest"
)

var mc *memcache.Client

func TestMain(m *testing.M) {
	var cmd *exec.Cmd
	memd := func(port int) {
		cmd = exec.Command("go", "run", "cmd/cirno_server/main.go", "-p", fmt.Sprintf("%d", port))
		cmd.SysProcAttr = &syscall.SysProcAttr{
			Setpgid: true,
		}
		cmd.Run()
	}

	server, err := tcptest.Start(memd, 30*time.Second)
	if err != nil {
		log.Fatalf("Failed to start cirno server: %s", err)
	}

	log.Printf("cirno server started on port %d", server.Port())
	defer func() {
		if cmd != nil && cmd.Process != nil {
			cmd.Process.Signal(syscall.SIGTERM)
		}
	}()
	mc = memcache.New(fmt.Sprintf("localhost:%s", server.Port()))

	// execute test for cirno server
	code := m.Run()

	// Then when you're done, you need to kill it
	cmd.Process.Signal(syscall.SIGTERM)

	// And wait
	server.Wait()
	os.Exit(code)
}

func BenchmarkSingle(b *testing.B) {
	for i := 0; i < b.N; i++ {
		DoAnyThing()
	}
}

func BenchmarkParallel(b *testing.B) {
	b.SetParallelism(5)
	b.ResetTimer()
	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			DoAnyThing()
		}
	})
}

func DoAnyThing() {
	mc.Get("id")
}
