package procfile

import (
	"bufio"
	"io"
	"os"
	"strings"
)

func ParseFile(path string) ([]*Proc, error) {
	file, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer file.Close()

	return Parse(file), nil
}

func Parse(file io.Reader) []*Proc {
	var processes []*Proc
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

		processes = append(processes, New(parts[0], args))
	}

	return processes
}
