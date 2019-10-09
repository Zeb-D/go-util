package interact

import (
	"fmt"
	"os"
	"os/exec"
	"testing"
)

func TestStartProcess(t *testing.T) {
	// 1) os.StartProcess //
	/*********************/
	/* Linux: */
	env := os.Environ()
	procAttr := &os.ProcAttr{
		Env: env,
		Files: []*os.File{
			os.Stdin,
			os.Stdout,
			os.Stderr,
		},
	}
	// 1st example: list files
	pid, err := os.StartProcess("/bin/ls", []string{"ls", "-l"}, procAttr)
	if err != nil {
		fmt.Printf("Error %v starting process!", err) //
		os.Exit(1)
	}
	fmt.Printf("The process id is %v", pid)

	// 2nd example: show all processes
	pid, err = os.StartProcess("/bin/ps", []string{"-e", "-opid,ppid,comm"}, procAttr)
	if err != nil {
		fmt.Printf("Error %v starting process!", err) //
		os.Exit(1)
	}
	fmt.Printf("The process id is %v", pid)
}

func TestExec(t *testing.T) {
	// 2) exec.Run //
	/***************/
	// Linux:  OK, but not for ls ?
	cmd := exec.Command("ls", "-l")  // no error, but doesn't show anything ?
	// cmd := exec.Command("ls")        // no error, but doesn't show anything ?
	//cmd := exec.Command("gedit") // this opens a gedit-window
	err := cmd.Run()
	if err != nil {
		fmt.Printf("Error %v executing command!", err)
		os.Exit(1)
	}
	fmt.Printf("The command is %v", cmd)
}
