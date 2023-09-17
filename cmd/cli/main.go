package main

import "github.com/zcubbs/hotpot/cmd/cli/cmd"

var (
	Version = "dev"
	Commit  = "none"
	Date    = "unknown"
)

func init() {
	cmd.Version = Version
	cmd.Commit = Commit
	cmd.Date = Date
}

func main() {
	cmd.Execute()
}
