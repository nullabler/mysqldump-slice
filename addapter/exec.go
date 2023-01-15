package addapter

import (
	"os"
	"os/exec"
)

type ExecInterface interface {
	Command(string) error
}

type Exec struct {
	sh string
}

func NewExec(sh string) *Exec {
	return &Exec{
		sh: sh,
	}
}

func (e *Exec) Command(call string) error {
	cmd := exec.Command(e.sh, "-c", call)
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		return err
	}

	return nil
}
