package cmd

import (
	"frp-auth/serve"

	"github.com/spf13/cobra"
)

var ServeCmd = &cobra.Command{
	Use:   "serve",
	Short: "Serves a API",
	Run: func(cmd *cobra.Command, args []string) {
		serve.Serve(DistFS)
	},
}
