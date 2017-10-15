package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

func main() {
	// Takes filename (without .go) as user input
	flagFile := flag.String("f", "", "File to be compiled and executed at target.")
	flag.Parse()

	oss := "Linux"
	arm := "7"
	arch := "ARM"

	// Start build process
	buildArg := "/home/tuomas/go/src/github.com/stuomas/piserve/" + *flagFile + ".go"
	cmdBuild := exec.Command("go", "build", buildArg)
	cmdBuild.Env = append(os.Environ(),
		"GOARM=7",
		"GOOS=linux",
		"GOARCH=arm",
	)
	fmt.Printf("Cross-compiling for %s %s...", arch+arm, oss)
	errBuild := cmdBuild.Run()
	if errBuild != nil {
		fmt.Printf("Failure!\n")
	} else {
		fmt.Printf("Success!\n")
	}

	// Start scp process
	targetLoc := "pi@rasp"
	targetDir := "/home/pi/go/bin/"
	scpArgs := []string{*flagFile, targetLoc + ":" + targetDir}
	cmdScp := exec.Command("scp", scpArgs...)
	fmt.Printf("Copying binary to target %s...", targetDir)
	errScp := cmdScp.Run()
	if errScp != nil {
		fmt.Printf("Failure!\n")
	} else {
		fmt.Printf("Success!\n")
	}

	// Execute via SSH at target
	sshArgs := []string{"-t", "pi@rasp", "'./main'"}
	cmdSSH := exec.Command("ssh", sshArgs...)
	fmt.Printf("Running freshly built binary at target %s...", targetLoc)
	errSSH := cmdSSH.Run()
	if errSSH != nil {
		fmt.Printf("Failure!\n")
	} else {
		fmt.Printf("Success!\n")
	}
}
