// +build !windows
//
// Copyright 2016 HenryLee. All Rights Reserved.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//      http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package faygo

import (
	"context"
	"os"
	"os/signal"
	"syscall"
	"time"
)

func graceSignal() {
	WritePid(LogDir() + "app.pid")
	// subscribe to SIGINT signals
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	defer func() {
		os.Exit(0)
	}()
	sig := <-ch
	signal.Stop(ch)
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		Shutdown()
	case syscall.SIGUSR2:
		Reboot()
	}
}

// Reboot all the frame services gracefully.
// Notes: Windows system are not supported!
func Reboot(timeout ...time.Duration) {
	global.framesLock.Lock()
	defer global.framesLock.Unlock()
	defer CloseLog()
	Print("\x1b[46m[SYS]\x1b[0m rebooting services...")

	var (
		ppid     = os.Getppid()
		graceful = true
	)
	contextExec(timeout, "reboot", func(ctxTimeout context.Context) <-chan struct{} {
		endCh := make(chan struct{})
		go func() {
			defer close(endCh)

			var reboot = true

			if global.preCloseFunc != nil {
				if err := global.preCloseFunc(); err != nil {
					Errorf("[reboot-preClose] %s", err.Error())
					graceful = false
				}
			}

			// Starts a new process passing it the active listeners. It
			// doesn't fork, but starts a new process using the same environment and
			// arguments as when it was originally started. This allows for a newly
			// deployed binary to be started.
			_, err := grace.StartProcess()
			if err != nil {
				Errorf("[reboot-startNewProcess] %s", err.Error())
				reboot = false
			}

			// shut down
			graceful = shutdown(ctxTimeout, "reboot") && graceful
			if !reboot {
				if graceful {
					Fatalf("services reboot failed, but shut down gracefully!")
				} else {
					Fatalf("services reboot failed, and did not shut down gracefully!")
				}
				os.Exit(-1)
			}
		}()

		return endCh
	})

	// Close the parent if we inherited and it wasn't init that started us.
	if ppid != 1 {
		if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
			Errorf("[reboot-killOldProcess] %s", err.Error())
			graceful = false
		}
	}

	if graceful {
		Print("\x1b[46m[SYS]\x1b[0m services are rebooted gracefully.")
	} else {
		Print("\x1b[46m[SYS]\x1b[0m services are rebooted, but not gracefully.")
	}
}
