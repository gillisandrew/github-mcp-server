package main

import (
	"fmt"
	"os"

	"github.com/github/github-mcp-server/pkg/cmd"
)

var version = "version"
var commit = "commit"
var date = "date"

func main() {
	rootCmd := cmd.GetRootCmd()
	rootCmd.Version = fmt.Sprintf("%s (%s) %s", version, commit, date)
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
