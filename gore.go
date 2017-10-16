package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

//func parseFlags
//func goBuild
//func runSCP
//func runSSH

func main() {
	// TODO: Make separate functions

	// Takes filename (without .go) as user input
	flagFile := flag.String("f", "", "File to be compiled and executed at target.")
	flag.Parse()

	// TODO: Add as flags and with defaults read from configuration file (toml?)
	osys := "linux"
	armv := "7"
	arch := "ARM"
	targ := "rasp"
	user := "pi"
	tdir := "/home/pi/go/bin/"

	// Start build process
	buildArg := "/home/tuomas/go/src/github.com/stuomas/piserve/" + *flagFile + ".go"
	cmdBuild := exec.Command("go", "build", buildArg)
	cmdBuild.Env = append(os.Environ(),
		"GOARM=" + armv,
		"GOOS=" + osys,
		"GOARCH=" + arch,
	)
	fmt.Printf("Cross-compiling for %s %s...", arch+armv, osys)
	errBuild := cmdBuild.Run()
	if errBuild != nil {
		fmt.Printf("Failure!\n")
	} else {
		fmt.Printf("Success!\n")
	}

	// Start scp process
	scpArgs := []string{*flagFile, user + "@" + targ + ":" + tdir}
	cmdScp := exec.Command("scp", scpArgs...)
	fmt.Printf("Copying binary to target %s...", tdir)
	errScp := cmdScp.Run()
	if errScp != nil {
		fmt.Printf("Failure!\n")
	} else {
		fmt.Printf("Success!\n")
	}

	// Execute via SSH at target
	sshArgs := []string{"-t", "pi@rasp", "'./main'"}
	cmdSSH := exec.Command("ssh", sshArgs...)
	fmt.Printf("Running freshly built binary at target...")
	errSSH := cmdSSH.Run()
	if errSSH != nil {
		fmt.Printf("Failure!\n")
	} else {
		fmt.Printf("Success!\n")
	}
}
