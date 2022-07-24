package main

import (
	"errors"
	"os"
	"os/exec"
)

// RunCmd runs a command + arguments (cmd) with environment variables from env.
func RunCmd(cmd []string, env Environment) (returnCode int) {
	var (
		err         error
		command     *exec.Cmd
		commandName string
		commandArgs []string
	)
	if len(cmd) > 0 {
		commandName = cmd[0]
	}
	if len(cmd) > 1 {
		commandArgs = cmd[1:]
	}

	for k, v := range env {
		if _, ok := os.LookupEnv(k); ok {
			if err = os.Unsetenv(k); err != nil {
				return 111
			}
		}
		if v.NeedRemove {
			continue
		}
		if err = os.Setenv(k, v.Value); err != nil {
			return 111
		}
	}

	command = exec.Command(commandName, commandArgs...)
	command.Stdin = os.Stdin
	command.Stdout = os.Stdout
	command.Stderr = os.Stderr
	command.Env = os.Environ()

	err = command.Start()
	if err != nil {
		return 111
	}

	err = command.Wait()
	if err == nil {
		return 0
	}

	var exiterr *exec.ExitError
	if errors.As(err, &exiterr) {
		return exiterr.ExitCode()
	}

	return 1
}
