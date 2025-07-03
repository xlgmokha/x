package main

import (
	"bufio"
	"os"
	"os/exec"
	"os/signal"
	"strings"
	"sync"
	"syscall"
)

func main() {
	file, _ := os.Open("Procfile")
	defer file.Close()

	var cmds []*exec.Cmd
	var wg sync.WaitGroup

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

		args := strings.Fields(strings.TrimSpace(parts[1]))
		if len(args) == 0 {
			continue
		}

		cmd := exec.Command(args[0], args[1:]...)
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
		cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}

		cmds = append(cmds, cmd)

		wg.Add(1)
		go func(c *exec.Cmd) {
			defer wg.Done()
			c.Start()
			c.Wait()
		}(cmd)
	}

	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)

	go func() {
		<-sigChan
		for _, cmd := range cmds {
			if cmd.Process != nil {
				syscall.Kill(-cmd.Process.Pid, syscall.SIGTERM)
			}
		}
	}()

	wg.Wait()
}
