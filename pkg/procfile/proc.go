package procfile

import (
	"os"
	"os/exec"
	"syscall"
)

type Proc struct {
	name string
	args []string
}

func New(name string, args []string) *Proc {
	return &Proc{
		name: name,
		args: args,
	}
}

func (p *Proc) NewCommand() *exec.Cmd {
	cmd := exec.Command(p.args[0], p.args[1:]...)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.SysProcAttr = &syscall.SysProcAttr{Setpgid: true}
	return cmd
}
