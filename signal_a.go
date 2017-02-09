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

package thinkgo

import (
	"context"
	"os"
	"os/signal"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func graceSignal() {
	ch := make(chan os.Signal)
	signal.Notify(ch, syscall.SIGINT, syscall.SIGTERM, syscall.SIGUSR2)
	defer signal.Stop(ch)
	sig := <-ch
	switch sig {
	case syscall.SIGINT, syscall.SIGTERM:
		Shutdown()
		return
	case syscall.SIGUSR2:
		reboot()
	}
}

// reboot all the frame services gracefully.
func reboot(timeout ...time.Duration) {
	Print("\x1b[46m[SYS]\x1b[0m rebooting servers...")
	defer CloseLog()
	global.framesLock.Lock()
	defer global.framesLock.Unlock()

	ppid := os.Getppid()

	// Starts a new process passing it the active listeners. It
	// doesn't fork, but starts a new process using the same environment and
	// arguments as when it was originally started. This allows for a newly
	// deployed binary to be started.
	_, err := grace.StartProcess()
	if err != nil {
		Error(err.Error())
		Print("\x1b[46m[SYS]\x1b[0m reboot servers failed, so close parent.")
		return
	}

	// Shut down gracefully, but wait no longer than d before halting
	var d = SHUTDOWN_TIMEOUT
	if len(timeout) > 0 {
		d = timeout[0]
	}
	ctxTimeout, _ := context.WithTimeout(context.Background(), d)
	count := new(sync.WaitGroup)
	var flag int32 = 1
	for _, frame := range global.frames {
		count.Add(1)
		go func(fm *Framework) {
			graceful := fm.shutdown(ctxTimeout)
			if !graceful {
				atomic.StoreInt32(&flag, 0)
			}
			count.Done()
		}(frame)
	}
	count.Wait()
	if global.finalizer != nil {
		if err := global.finalizer(ctxTimeout); err != nil {
			flag = 0
			Error("[finalizer]", err.Error())
		}
	}

	// Close the parent if we inherited and it wasn't init that started us.
	if ppid != 1 {
		if err := syscall.Kill(ppid, syscall.SIGTERM); err != nil {
			Error("failed to close parent: %s", err.Error())
			Print("\x1b[46m[SYS]\x1b[0m servers reboot failed, so close parent.")
			return
		}
	}

	if flag == 1 {
		Print("\x1b[46m[SYS]\x1b[0m servers are rebooted gracefully.")
	} else {
		Print("\x1b[46m[SYS]\x1b[0m servers are rebooted, but not gracefully.")
	}
}
