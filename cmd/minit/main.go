package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/xlgmokha/x/pkg/procfile"
)

var (
	pidMutex     sync.Mutex
	pids         []int
	procfilePath *string
)

func init() {
	procfilePath = flag.String("f", "Procfile", "path to Procfile")
	flag.Parse()
	log.SetFlags(0)
}

func addPid(pid int) {
	pidMutex.Lock()
	defer pidMutex.Unlock()
	pids = append(pids, pid)
}

func removePid(pid int) {
	pidMutex.Lock()
	defer pidMutex.Unlock()

	for i, p := range pids {
		if p == pid {
			pids = append(pids[:i], pids[i+1:]...)
			break
		}
	}
}

func forwardSignalToAll(sig os.Signal) {
	pidMutex.Lock()
	defer pidMutex.Unlock()

	signal := sig.(syscall.Signal)
	for _, pid := range pids {
		syscall.Kill(-pid, signal)
	}
}

func main() {
	var wg sync.WaitGroup
	var shutdown int32

	for _, path := range strings.Split(*procfilePath, ",") {
		procs, err := procfile.ParseFile(path)
		if err != nil {
			log.Fatalln(err)
		}

		for _, proc := range procs {
			wg.Add(1)
			go func(proc *procfile.Proc) {
				defer wg.Done()

				for atomic.LoadInt32(&shutdown) == 0 {
					cmd := proc.NewCommand()

					if cmd.Start() != nil {
						time.Sleep(2 * time.Second)
						continue
					}

					addPid(cmd.Process.Pid)
					cmd.Wait()
					removePid(cmd.Process.Pid)
					time.Sleep(time.Second)
				}
			}(proc)
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan)

	go func() {
		for sig := range sigChan {
			if sig == syscall.SIGINT || sig == syscall.SIGTERM {
				atomic.StoreInt32(&shutdown, 1)
			}

			forwardSignalToAll(sig)
		}
	}()

	wg.Wait()
}
