package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/require"
)

func TestCopy(t *testing.T) {
	var (
		in     = "testdata/input.txt"
		out    = "out.txt"
		offset int64
		limit  int64
	)

	t.Run("no input filename", func(t *testing.T) {
		err := Copy("", out, offset, limit)
		require.ErrorIs(t, err, ErrInputFileNotSpecified)
	})

	t.Run("no output filename", func(t *testing.T) {
		err := Copy(in, "", offset, limit)
		require.ErrorIs(t, err, ErrOutputFileNotSpecified)
	})

	t.Run("unsupported input file", func(t *testing.T) {
		err := Copy("/dev/urandom", out, offset, limit)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("unsupported output file", func(t *testing.T) {
		err := Copy(in, "/dev/null", offset, limit)
		require.ErrorIs(t, err, ErrUnsupportedFile)
	})

	t.Run("offset exceeds file size", func(t *testing.T) {
		info, _ := os.Stat(in)
		err := Copy(in, out, info.Size()+1, limit)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)
	})

	t.Run("limit exceeds file size", func(t *testing.T) {
		info, _ := os.Stat(in)
		fsize := info.Size()
		err := Copy(in, out, offset, fsize+10)
		require.Nil(t, err)
		info, _ = os.Stat(out)
		require.Equal(t, fsize, info.Size())
	})
}
