package main

import (
	"flag"
	"log"
	"os"
)

var (
	from, to      string
	limit, offset int64
)

func init() {
	flag.StringVar(&from, "from", "", "file to read from")
	flag.StringVar(&to, "to", "", "file to write to")
	flag.Int64Var(&limit, "limit", 0, "limit of bytes to copy")
	flag.Int64Var(&offset, "offset", 0, "offset in input file")
}

func main() {
	flag.Parse()

	logInfo := log.New(os.Stdout, "INFO: ", log.Ldate|log.Ltime)
	logError := log.New(os.Stderr, "ERROR: ", log.Ldate|log.Ltime)
	logInfo.Printf("copying from=%s to=%s offset=%d limit=%d", from, to, offset, limit)

	err := Copy(from, to, offset, limit)
	if err != nil {
		logError.Print(err)
		os.Exit(1)
	}
	logInfo.Print("success")
	os.Exit(0)
}
