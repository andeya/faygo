// +build windows
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
	"os"
	"os/signal"
	"time"
)

func graceSignal() {
	WritePid(LogDir() + "app.pid")
	// subscribe to SIGINT signals
	ch := make(chan os.Signal)
	signal.Notify(ch, os.Interrupt, os.Kill)
	defer func() {
		os.Exit(0)
	}()
	<-ch // wait for SIGINT
	signal.Stop(ch)
	Shutdown()
}

// Reboot all the frame services gracefully.
// Notes: Windows system are not supported!
func Reboot(timeout ...time.Duration) {
	Print("\x1b[46m[SYS]\x1b[0m windows system doesn't support reboot! call Shutdown() is recommended.")
}
