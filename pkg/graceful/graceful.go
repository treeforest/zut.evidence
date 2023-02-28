package graceful

import (
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

// Stop graceful stop
func Stop(f func()) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill, syscall.SIGTERM, syscall.SIGINT, syscall.SIGKILL, syscall.SIGHUP, syscall.SIGQUIT)
	s := <-c
	f()
	fmt.Println("graceful stopped.")
	if i, ok := s.(syscall.Signal); ok {
		os.Exit(int(i))
	} else {
		os.Exit(0)
	}
}
