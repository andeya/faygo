// Package gracehttp provides easy to use graceful restart
// functionality for HTTP server.
// modified by henrylee 2016.10.29

package gracehttp

import (
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
)

func (a *app) signalHandler(wg *sync.WaitGroup) {
	ch := make(chan os.Signal, 10)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	for {
		sig := <-ch
		switch sig {
		case syscall.SIGINT, syscall.SIGTERM:
			// this ensures a subsequent INT/TERM will trigger standard go behaviour of
			// terminating.
			signal.Stop(ch)
			if err := a.terminateFunc(); err != nil {
				a.errors <- err
			}
			a.term(wg)
			return
		case syscall.SIGUSR2:
			if err := a.terminateFunc(); err != nil {
				signal.Stop(ch)
				a.term(wg)
				return
			}
			// we only return here if there's an error, otherwise the new process
			// will send us a TERM when it's ready to trigger the actual shutdown.
			if _, err := a.net.StartProcess(); err != nil {
				a.errors <- err
			}
		}
	}
}

// ServeWithTerminateFunc will serve the given http.Servers and will monitor for signals
// allowing for graceful termination (SIGTERM) or restart (SIGUSR2).
func ServeWithTerminateFunc(terminateFunc func() error, servers ...*http.Server) error {
	a := newApp(servers)
	if terminateFunc == nil {
		a.terminateFunc = func() error {
			return nil
		}
	} else {
		a.terminateFunc = terminateFunc
	}

	// Acquire Listeners
	if err := a.listen(); err != nil {
		return err
	}

	// Some useful logging.
	if *verbose {
		if didInherit {
			if ppid == 1 {
				log.Printf("Listening on init activated %s", pprintAddr(a.listeners))
			} else {
				const msg = "Graceful handoff of %s with new pid %d and old pid %d"
				log.Printf(msg, pprintAddr(a.listeners), os.Getpid(), ppid)
			}
		} else {
			const msg = "Serving %s with pid %d"
			log.Printf(msg, pprintAddr(a.listeners), os.Getpid())
		}
	}

	// Start serving.
	a.serve()

	// Close the parent if we inherited and it wasn't init that started us.
	if didInherit && ppid != 1 {
		if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
			return fmt.Errorf("failed to close parent: %s", err)
		}
	}

	waitdone := make(chan struct{})
	go func() {
		defer close(waitdone)
		a.wait()
	}()

	select {
	case err := <-a.errors:
		if err == nil {
			panic("unexpected nil error")
		}
		return err
	case <-waitdone:
		if *verbose {
			log.Printf("Exiting pid %d.", os.Getpid())
		}
		return nil
	}
}
