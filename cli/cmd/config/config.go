// Package config
package config

import (
	"fmt"
	"pie-rum/cmd/keys"
	"pie-rum/cmd/policy"

	"github.com/spf13/cobra"
)

// just for testing

func Config() *cobra.Command {
	i := &cobra.Command{
		Use: "init",
	}
	i.AddCommand(Organization())
	return i
}

// Organization inits the log system
// note: for testing purpose we using the ip address
func Organization() *cobra.Command {
	cmd := &cobra.Command{
		Use:  "org [IP_ADDRESS]",
		Args: cobra.ExactArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			p, err := policy.UnPackInitPolicy(cmd)
			if err != nil {
				fmt.Printf("Error unpacking: %v\n", err)
				return
			}
			if p.CMDORG == nil {
				p.CMDORG = &policy.IOrganization{}
			}
			p.CMDORG.Key = args[0]
			p.CMDORG.Name = "organization"
			p.CMDORG.ShortHand = "init the org"

			cmd.Flags().Set(keys.Init_Flag, string(p.Pack()))

			fmt.Printf("Policy updated with key: %s\n", p.CMDORG.Key)
		},
	}
	return cmd
}
func Root() *cobra.Command {
	main := Config()
	main.AddCommand(Organization())
	return main
}
