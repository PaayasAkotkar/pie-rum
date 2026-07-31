// Package admin implements the pushing of all the command at one place
package admin

import (
	"pie-rum/cmd/config"
	"pie-rum/cmd/post"

	"github.com/spf13/cobra"
)

func Admin() []*cobra.Command {
	cmds := []*cobra.Command{config.Root(), post.Root()}
	return cmds
}
