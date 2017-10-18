package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

//Build the input source file
func runGoBuild(args []string, osys string, arch string, armv string) {
	cmdBuild := exec.Command("go", args...)
	cmdBuild.Env = append(os.Environ(),
		"GOARM="+armv,
		"GOOS="+osys,
		"GOARCH="+arch,
	)
	fmt.Printf("Cross-compiling for %s %s...", arch+armv, osys)
	stdoutStderr, err := cmdBuild.CombinedOutput()
	if err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("\x1b[32;1mSuccess!\x1b[0m\n")
	}
}

//Copy binary to target with scp
func runSCP(args []string) {
	cmdScp := exec.Command("scp", args...)
	fmt.Printf("Copying binary to target...")
	stdoutStderr, err := cmdScp.CombinedOutput()
	if err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("\x1b[32;1mSuccess!\x1b[0m\n")
	}
}

//Execute binary at target via SSH
func runSSH(args []string) {
	//TODO: password prompt not working, only key-based authentication works
	cmdSSH := exec.Command("ssh", args...)
	fmt.Printf("Running freshly built binary at target...")
	stdoutStderr, err := cmdSSH.CombinedOutput()
	if err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("\x1b[32;1mSuccess!\x1b[0m\n") //TODO: Not executed
	}
}

func main() {
	//TODO: Defaults read from configuration file (toml?)
	osysDefault := "linux"
	armvDefault := "7"
	archDefault := "arm"
	targDefault := "rasp"
	userDefault := "pi"
	tdirDefault := "/home/pi/go/bin/"

	//Flags set by user
	flagOsys := flag.String("os", osysDefault, "Target operating system.")
	flagArch := flag.String("arch", archDefault, "Target architecture.")
	flagArmv := flag.String("arm", armvDefault, "ARM version.")
	flagTarg := flag.String("host", targDefault, "Target IP or hostname.")
	flagUser := flag.String("user", userDefault, "Username at target.")
	flagTdir := flag.String("dir", tdirDefault, "Target directory.")
	flag.Parse()

	var sourcePath string
	var localPath string
	var packageName string
	var buildArgs []string

	if len(flag.Args()) > 0 {
		sourcePath = flag.Arg(0)
		packageName = strings.TrimSuffix(filepath.Base(sourcePath), ".go")
		buildArgs = []string{"build", sourcePath}
		localPath = strings.TrimSuffix(sourcePath, filepath.Base(sourcePath))
		if !strings.HasSuffix(sourcePath, ".go") {
			fmt.Println("Absolute path needs to target a source file.")
			os.Exit(1)
		}
	} else {
		localPath = ""
		sourcePath, _ = filepath.Abs(filepath.Dir(os.Args[0]))
		packageName = filepath.Base(sourcePath)
		buildArgs = []string{"build"}
	}

	runGoBuild(buildArgs, *flagOsys, *flagArch, *flagArmv)
	runSCP([]string{localPath + packageName, *flagUser + "@" + *flagTarg + ":" + *flagTdir})
	runSSH([]string{"-t", "-t", *flagUser + "@" + *flagTarg, "'" + *flagTdir + packageName + "'"})
}
