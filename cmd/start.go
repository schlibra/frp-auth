package cmd

import (
	"fmt"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"syscall"

	"github.com/spf13/cobra"
)

const PidFile = "frp-auth.pid"
const LogFile = "frp-auth.log"

func readPid() (int, error) {
	data, err := os.ReadFile(PidFile)
	if err != nil {
		return 0, nil
	}
	return strconv.Atoi(strings.TrimSpace(string(data)))
}
func isProcessRunning(pid int) bool {
	proc, err := os.FindProcess(pid)
	if err != nil {
		return false
	}
	return proc.Signal(syscall.Signal(0)) == nil
}

func startService() {
	if pid, err := readPid(); err == nil {
		if isProcessRunning(pid) {
			fmt.Printf("Frp-Auth already running(PID: %d)!\n", pid)
			return
		}
		_ = os.Remove(PidFile)
	}
	logOutput, err := os.OpenFile(LogFile, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0644)
	if err != nil {
		fmt.Printf("Failed to open log file: %v\n", err)
		return
	}
	_cmd := exec.Command(os.Args[0], "serve")
	_cmd.Stdout = logOutput
	_cmd.Stderr = logOutput
	err = _cmd.Start()
	if err != nil {
		fmt.Printf("Failed to start server: %v\n", err)
		return
	}
	pid := _cmd.Process.Pid
	if err := os.WriteFile(PidFile, []byte(fmt.Sprintf("%d", pid)), 0644); err != nil {
		fmt.Printf("Failed to write pid to file: %v\n", err)
		return
	}
	fmt.Printf("Frp-Auth started with PID: %d\n", pid)
}

var StartCmd = &cobra.Command{
	Use:   "start",
	Short: "Start the server",
	Run: func(cmd *cobra.Command, args []string) {
		startService()
	},
}
