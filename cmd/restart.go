package cmd

import "github.com/spf13/cobra"

var RestartCmd = &cobra.Command{
	Use:   "restart",
	Short: "Restart the server",
	Run: func(cmd *cobra.Command, args []string) {
		stopService()
		startService()
	},
}
