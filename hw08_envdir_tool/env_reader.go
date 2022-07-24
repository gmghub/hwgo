package main

import (
	"bufio"
	"errors"
	"io/ioutil"
	"os"
	"path"
	"strings"
)

var (
	ErrReadingEnvDir = errors.New("failed to read envdir")
	ErrNotDirectory  = errors.New("path is not a directory")
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	d, err := os.Stat(dir)
	if err != nil {
		return nil, ErrReadingEnvDir
	}
	if !d.IsDir() {
		return nil, ErrNotDirectory
	}

	files, err := ioutil.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	for _, f := range files {
		if !f.Mode().IsRegular() {
			continue
		}
		name := f.Name()
		if strings.Contains(name, "=") {
			continue
		}
		fd, err := os.Open(path.Join(dir, name))
		if err != nil {
			continue
		}

		scanner := bufio.NewScanner(fd)
		if !scanner.Scan() {
			env[name] = EnvValue{Value: "", NeedRemove: true}
		} else {
			value := scanner.Text()
			value = strings.ReplaceAll(value, "\000", "\n")
			value = strings.TrimRight(value, " \t\n")
			env[name] = EnvValue{Value: value, NeedRemove: false}
		}
		fd.Close()
	}
	return env, nil
}
