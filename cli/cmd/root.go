/*
Copyright © 2026 PAAYAS AKOTKAR <paayasakotkar@gmail.com>
*/
package cmd

import (
	"fmt"
	"log"
	"os"
	"pie-rum/cmd/admin"
	"pie-rum/cmd/keys"
	"pie-rum/cmd/policy"

	"github.com/spf13/cobra"
)

// Version verions-> 0.0.0 => major.minor.patch

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:     "pie-rum",
	Short:   "root pie-rum cli",
	Long:    `PIE RUM CLI tool serves the connection to the sdk system allowing the toggle, architecture & snapshot metadata management.`,
	Version: "0.1.0",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		if v, _ := cmd.Flags().GetBool("version"); v {
			fmt.Println("0.1.0")
			os.Exit(0)
		}
		p := policy.InitPolicy{ /* Initialize your data here */ }
		cmd.Flags().Set(keys.Init_Flag, string(p.Pack()))

	},
}

func init() {

	log.SetFlags(log.Lshortfile)
	rootCmd.PersistentFlags().BoolP("version", "v", false, "Print the version")
	rootCmd.PersistentFlags().String(keys.Init_Flag, "", "policy initialization")
	for _, c := range admin.Admin() {
		rootCmd.AddCommand(c)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
