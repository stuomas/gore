package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
)

//Build the input source file
func runGoBuild(file string, osys string, arch string, armv string) {
	cmdBuild := exec.Command("go", "build", file)
	cmdBuild.Env = append(os.Environ(), 
		"GOARM=" + armv,
		"GOOS=" + osys,
		"GOARCH=" + arch,
	)
	fmt.Printf("Cross-compiling for %s %s...", arch + armv, osys)
	stdoutStderr, errBuild := cmdBuild.CombinedOutput()
	if errBuild != nil {
		fmt.Printf("Failed!\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("Success!\n \u2937 %s", stdoutStderr)
	}
}

//Copy binary to target with scp
func runSCP(args []string) {
	cmdScp := exec.Command("scp", args...)
	fmt.Printf("Copying binary to target...")
	stdoutStderr, errScp := cmdScp.CombinedOutput()
	if errScp != nil {
		fmt.Printf("Failed!\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("Success!\n \u2937 %s", stdoutStderr)
	}
}

//Execute binary at target via SSH
func runSSH(args []string) {
	//TODO: password prompt not working, only key-based authentication works 
	cmdSSH := exec.Command("ssh", args...)
	fmt.Printf("Running freshly built binary at target...")
	stdoutStderr, errSSH := cmdSSH.CombinedOutput()
	if errSSH != nil {
		fmt.Printf("Failed!\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("Success!\n \u2937 %s", stdoutStderr)
	}
}

func main() {
	//TODO: Defaults read from configuration file (toml?)
	osysDefault := "linux"
	armvDefault := "7"
	archDefault := "ARM"
	targDefault := "rasp"
	userDefault := "pi"
	tdirDefault := "/home/pi/go/bin/" //$GOBIN?

	//User input arguments
	flagFile := flag.String("f", "", "Source file to be compiled and executed at target.")
	flagOsys := flag.String("o", osysDefault, "Target operating system.")
	flagArch := flag.String("a", archDefault, "Target architecture.")
	flagArmv := flag.String("v", armvDefault, "ARM version.")
	flagTarg := flag.String("t", targDefault, "Target IP or hostname.")
	flagUser := flag.String("u", userDefault, "Username at target.")
	flagTdir := flag.String("d", tdirDefault, "Target directory.")

	flag.Parse()

	if *flagFile == "" {
		fmt.Println("Please specify a source file.")
		os.Exit(1)
	}

	runGoBuild(*flagFile + ".go", *flagOsys, *flagArch, *flagArmv)
	runSCP([]string{*flagFile, *flagUser + "@" + *flagTarg + ":" + *flagTdir})
	runSSH([]string{"-t", *flagUser + "@" + *flagTarg, "'./" + *flagFile + "'"})	
}