package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var StatusCmd = &cobra.Command{
	Use:   "status",
	Short: "Show status of Frp-Auth service",
	Run: func(cmd *cobra.Command, args []string) {
		pid, err := readPid()
		if err != nil || !isProcessRunning(pid) {
			fmt.Println("Service not running")
			return
		}
		fmt.Printf("Frp-Auth service is running(PID: %d)!\n", pid)
	},
}
