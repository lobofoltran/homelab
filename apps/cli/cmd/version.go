package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Show the version of the CLI",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("homelab v0.1.0")
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
