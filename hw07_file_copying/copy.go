package main

import (
	"errors"
	"io"
	"os"

	pb "github.com/cheggaaa/pb/v3"
)

var (
	ErrUnsupportedFile        = errors.New("unsupported file")
	ErrOffsetExceedsFileSize  = errors.New("offset exceeds file size")
	ErrInputFileNotSpecified  = errors.New("input file not specified")
	ErrOutputFileNotSpecified = errors.New("output file not specified")
	ErrPartiallyCopied        = errors.New("file partially copied")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	var err error
	if fromPath == "" {
		return ErrInputFileNotSpecified
	}
	if toPath == "" {
		return ErrOutputFileNotSpecified
	}

	inFileInfo, err := os.Stat(fromPath)
	if err != nil || !inFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}
	if offset > inFileInfo.Size() {
		return ErrOffsetExceedsFileSize
	}

	outFileInfo, err := os.Stat(toPath)
	if err == nil && !outFileInfo.Mode().IsRegular() {
		return ErrUnsupportedFile
	}

	inFile, err := os.OpenFile(fromPath, os.O_RDONLY, 0)
	if err != nil {
		return err
	}
	defer inFile.Close()

	outFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer outFile.Close()

	var (
		totalToWrite int64
		totalWritten int64
		blockSize    int64
	)

	blockSize = 4096
	totalToWrite = inFileInfo.Size()
	if offset != 0 {
		totalToWrite -= offset
		if _, err = inFile.Seek(offset, 0); err != nil {
			return err
		}
	}
	if limit != 0 && limit < totalToWrite {
		totalToWrite = limit
	}

	bar := pb.Full.Start64(totalToWrite)
	reader := io.LimitReader(inFile, totalToWrite)
	barReader := bar.NewProxyReader(reader)

	for {
		written, err := io.CopyN(outFile, barReader, blockSize)
		totalWritten += written
		if errors.Is(err, io.EOF) {
			break
		}
	}

	bar.Finish()

	if totalWritten != totalToWrite {
		return ErrPartiallyCopied
	}

	return nil
}
