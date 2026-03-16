package main

import (
	"fmt"
	"os"

	"github.com/alecthomas/kong"
	"github.com/anthropics/docs7/api"
	"github.com/anthropics/docs7/cmd"
)

var cli struct {
	Search cmd.SearchCmd `cmd:"" help:"Find a library by name"`
	Docs   cmd.DocsCmd   `cmd:"" help:"Get documentation source URLs for a library"`
	Read   cmd.ReadCmd   `cmd:"" help:"Read a document URL (github:// or https://)"`
	Login  cmd.LoginCmd  `cmd:"" help:"Log in to Context7 (opens browser)"`
	Logout cmd.LogoutCmd `cmd:"" help:"Log out and clear stored tokens"`
	Whoami cmd.WhoamiCmd `cmd:"" help:"Show current login status"`
}

func main() {
	ctx := kong.Parse(&cli,
		kong.Name("docs7"),
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
