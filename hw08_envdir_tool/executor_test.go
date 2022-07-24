package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestRunCmd(t *testing.T) {
	var (
		env  Environment
		cmd  []string
		code int
	)

	t.Run("empty command", func(t *testing.T) {
		cmd = []string{}
		code = RunCmd(cmd, env)
		require.Equal(t, 111, code)
	})

	t.Run("command name only", func(t *testing.T) {
		cmd = []string{"echo"}
		code = RunCmd(cmd, env)
		require.Equal(t, 0, code)
	})

	t.Run("command with args", func(t *testing.T) {
		cmd = []string{"echo", "123", "321"}
		code = RunCmd(cmd, env)
		require.Equal(t, 0, code)
	})

	t.Run("command with env vars", func(t *testing.T) {
		env = Environment{
			"HELLO": EnvValue{Value: "hello"},
			"BAR":   EnvValue{Value: "bar"},
			"FOO":   EnvValue{Value: "foo"},
		}
		cmd = []string{"testdata/echo.sh"}
		code = RunCmd(cmd, env)
		require.Equal(t, 0, code)
		osEnv := os.Environ()
		require.Contains(t, osEnv, "HELLO=hello")
		require.Contains(t, osEnv, "BAR=bar")
		require.Contains(t, osEnv, "FOO=foo")
	})
}
