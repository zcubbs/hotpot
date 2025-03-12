package os

import "github.com/dustin/go-humanize"

func StringToBytes(s string) (uint64, error) {
	return humanize.ParseBytes(s)
}

func BytesToString(b uint64) string {
	return humanize.Bytes(b)
}
