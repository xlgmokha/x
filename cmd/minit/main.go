package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"

	"github.com/xlgmokha/x/pkg/procfile"
)

var procfilePath *string

func init() {
	procfilePath = flag.String("f", "Procfile", "path to Procfile")
	flag.Parse()
	log.SetFlags(0)
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
				var cmd *exec.Cmd

				for atomic.LoadInt32(&shutdown) == 0 {
					cmd = proc.NewCommand()

					if cmd.Start() != nil {
						time.Sleep(2 * time.Second)
						continue
					}
					cmd.Wait()
					time.Sleep(time.Second)
				}

				if cmd != nil && cmd.Process != nil {
					syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
				}
			}(proc)
		}
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		atomic.StoreInt32(&shutdown, 1)
	}()

	wg.Wait()
}
