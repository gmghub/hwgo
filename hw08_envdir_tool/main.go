package main

import (
	"fmt"
	"os"
)

func usage() {
	fmt.Println("go-envdir\nUsage:\n  go-envdir /path/to/env/dir command arg1 arg2")
}

// EXIT CODES (man envdir)
// envdir exits 111 if it has trouble reading d, if it runs out of memory for environment variables,
// or if it cannot run child. Otherwise its exit code is the same as that of child.
func main() {
	if len(os.Args) < 3 {
		usage()
		os.Exit(1)
	}

	var (
		path string
		cmd  []string
	)

	path = os.Args[1]
	cmd = os.Args[2:]

	env, err := ReadDir(path)
	if err != nil {
		fmt.Println(err)
		os.Exit(111)
	}

	returnCode := RunCmd(cmd, env)
	os.Exit(returnCode)
}
