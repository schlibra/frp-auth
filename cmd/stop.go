package cmd

import (
	"fmt"
	"os"
	"syscall"
	"time"

	"github.com/spf13/cobra"
)

func stopService() {
	pid, err := readPid()
	if err != nil {
		fmt.Println("No running service(pid file not exist)")
		return
	}
	if !isProcessRunning(pid) {
		fmt.Println("Process not running")
		_ = os.Remove(PidFile)
		return
	}
	proc, _ := os.FindProcess(pid)
	if err := proc.Signal(syscall.SIGTERM); err != nil {
		fmt.Printf("Failed to stop: %v\n", err)
		return
	}
	fmt.Printf("Stopping service (PID: %d)\n", pid)
	for range 30 {
		if !isProcessRunning(pid) {
			break
		}
		time.Sleep(200 * time.Millisecond)
	}
	if isProcessRunning(pid) {
		_ = proc.Signal(syscall.SIGKILL)
		fmt.Println("Timeout, force kill")
	} else {
		fmt.Println("Service stopped")
	}
	_ = os.Remove(PidFile)
}

var StopCmd = &cobra.Command{
	Use:   "stop",
	Short: "stop command",
	Run: func(cmd *cobra.Command, args []string) {
		stopService()
	},
}
