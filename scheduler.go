package main

import (
	"log"
	"os/exec"
	"syscall"
)

//RunTask run task
func (run *Run) RunTask(c *Cron) {
	log.Printf("start task %v ...", c)
	for {
		select {
		case <-run.Ctx.Done():
			return
		case <-c.tk.C:
			ExecShellScript(c.CMD, c.Args)
		}
	}

}

// ExecShellScript runs shell command.
func ExecShellScript(command string, args []string) (int, []byte, error) {
	var waitStatus syscall.WaitStatus
	cmd := exec.Command(command, args...)
	out, err := cmd.CombinedOutput()
	log.Printf("command %s %v, output: %s", command, args, string(out))
	if err != nil {
		return -1, nil, err
	}
	waitStatus = cmd.ProcessState.Sys().(syscall.WaitStatus)

	return int(waitStatus), out, nil
}
