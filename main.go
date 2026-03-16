package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/ethan-huo/ctx/api"
	"github.com/ethan-huo/ctx/cmd"
)

var cli struct {
	Search cmd.SearchCmd `cmd:"" help:"Find a library by name"`
	Docs   cmd.DocsCmd   `cmd:"" help:"Get documentation source URLs for a library"`
	Read   cmd.ReadCmd   `cmd:"" help:"Read a document URL (github:// or https://)"`
	Auth   cmd.AuthCmd   `cmd:"" help:"Manage authentication"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("ctx"),
		kong.Description("Library documentation finder — ctx7 index + full document URLs"),
		kong.UsageOnError(),
	)

	client := api.NewClient()
	err := ctx.Run(client)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		os.Exit(1)
	}
}
