// Package life will control the life of the program
// If you have some cleaning operations that the program needs to handle
package life

import (
	"os"
	"os/signal"
	"syscall"
)

// Listener is a type of func()
type Listener func()

var (
	exitListeners    []Listener
	restartListeners []Listener
)

// WhenExit will register a set of destructor functions
func WhenExit(liss ...Listener) {
	exitListeners = append(exitListeners, liss...)
}

// WhenRestart ...
func WhenRestart(liss ...Listener) {
	restartListeners = append(restartListeners, liss...)
}

// Start should be used at the end of main() function because it will block the program waitting for signals
func Start() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, syscall.SIGHUP, syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT)
	for {
		switch <-c {
		case syscall.SIGQUIT, syscall.SIGTERM, syscall.SIGINT:
			for _, lis := range exitListeners {
				lis()
			}
			return
		case syscall.SIGHUP:
			for _, lis := range restartListeners {
				lis()
			}
		default:
			// Others
			return
		}
	}
}
