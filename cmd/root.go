package cmd

import (
	"embed"
	"log"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var DistFS embed.FS

func Execute(distFS embed.FS) {
	DistFS = distFS
	var rootCmd = &cobra.Command{
		Use:   filepath.Base(os.Args[0]),
		Short: "frp-auth",
		Long:  "frp-auth",
	}
	rootCmd.AddCommand(InitCmd)
	rootCmd.AddCommand(CleanCmd)
	rootCmd.AddCommand(ServeCmd)
	rootCmd.AddCommand(StartCmd)
	rootCmd.AddCommand(StopCmd)
	rootCmd.AddCommand(StatusCmd)
	rootCmd.AddCommand(RestartCmd)
	rootCmd.AddCommand(AdminCmd)
	err := rootCmd.Execute()
	if err != nil {
		log.Fatal(err)
	}
}
