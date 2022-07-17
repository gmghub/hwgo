package main

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestReadDir(t *testing.T) {
	t.Run("empty dir name", func(t *testing.T) {
		dir := ""
		env, err := ReadDir(dir)
		assert.Nil(t, env)
		assert.ErrorIs(t, err, ErrReadingEnvDir)
	})

	t.Run("dir not exists", func(t *testing.T) {
		dir := "notexists"
		env, err := ReadDir(dir)
		assert.Nil(t, env)
		assert.ErrorIs(t, err, ErrReadingEnvDir)
	})

	t.Run("path is not a directory", func(t *testing.T) {
		dir := "testdata/echo.sh"
		env, err := ReadDir(dir)
		assert.Nil(t, env)
		assert.ErrorIs(t, err, ErrNotDirectory)
	})

	t.Run("dir exists and has keys", func(t *testing.T) {
		dir := "testdata/env"
		expected := Environment{
			"BAR":   EnvValue{Value: "bar", NeedRemove: false},
			"HELLO": EnvValue{Value: `"hello"`, NeedRemove: false},
			"FOO":   EnvValue{Value: "   foo\nwith new line", NeedRemove: false},
			"EMPTY": EnvValue{Value: "", NeedRemove: false},
			"UNSET": EnvValue{Value: "", NeedRemove: true},
		}
		env, err := ReadDir(dir)
		assert.Nil(t, err)
		for k, v := range expected {
			assert.Equal(t, v, env[k], "not valid env value key = %v, %v != %v", k, v, env[k])
		}
	})
}
