package main

import (
	"bufio"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

func main() {
	file, _ := os.Open("Procfile")
	defer file.Close()

	var cmds []*exec.Cmd
	var wg sync.WaitGroup
	var shutdown int32

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		parts := strings.SplitN(line, ":", 2)
		if len(parts) != 2 {
			continue
		}

		args := strings.Fields(os.ExpandEnv(strings.TrimSpace(parts[1])))
		if len(args) == 0 {
			continue
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		cmds = append(cmds, cmd)

		wg.Add(1)
		go func(args []string) {
			defer wg.Done()
			for atomic.LoadInt32(&shutdown) == 0 {
				cmd := exec.Command(args[0], args[1:]...)
				cmd.Stdout = os.Stdout
				cmd.Stderr = os.Stderr
				cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

				if cmd.Start() != nil {
					time.Sleep(2 * time.Second)
					continue
				}
				cmd.Wait()
				time.Sleep(time.Second)
			}
		}(args)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		atomic.StoreInt32(&shutdown, 1)

		for _, cmd := range cmds {
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
			}
		}
	}()

	wg.Wait()
}
