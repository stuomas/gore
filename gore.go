package main

import (
	"flag"
	"bufio"
	"fmt"
	"errors"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"io/ioutil"
	"github.com/BurntSushi/toml"
)

//Configuration file structure
type Configuration struct {
	GOOS, GOARCH, GOARM, USERNAME, HOSTNAME, DIRECTORY string
}

func askConfig() ([]string) {
	configName := []string{"GOOS", "GOARCH", "GOARM", "USERNAME", "HOSTNAME", "DIRECTORY"}
	var configValue []string
	fmt.Printf("It seems that you just started using gore. Please generate a configuration file for you use by answering the prompts below:\n")
	prompt := bufio.NewReader(os.Stdin)
	fmt.Printf("Set the environment variables to match your target system:\n")

	for i := 0; i < len(configName); i++ {
		fmt.Print("	" + configName[i] + "=")
		configValue[i], _ = prompt.ReadString('\n')
	}
	return configValue
}

//Write initial config file and folder
func writeConfig(defConf []string) {
	confDir := filepath.Join(os.Getenv("HOME"), ".config")
	if err := os.Mkdir(confDir, 0755); err != nil {
		fmt.Println("Error creating folder.")
		os.Exit(1)
	}
	confSlice := []byte(fmt.Sprintf("GOOS=%s\nGOARCH=%s\nGOARM=%s\nUSERNAME=%s\nHOSTNAME=%s\nDIRECTORY=%s", defConf[0], defConf[1], defConf[2], defConf[3], defConf[4], defConf[5]))
	if err := ioutil.WriteFile(confDir, confSlice, 0755); err != nil {
		fmt.Println("Error writing configuration file.")
		os.Exit(1)
	}
}

//Read configuration file for environment variables and connection parameters
func readConfig() (Configuration, error) {
	var config Configuration
	confPath := filepath.Join(os.Getenv("HOME"), ".config")
	if _, err := os.Stat(confPath + "/config.toml"); !os.IsNotExist(err) {
		return config, errors.New("File exists.")
	} 
	if _, err := toml.DecodeFile(confPath+"/gore/config.toml", &config); err != nil {
		//fmt.Println(confPath + "/gore/config.toml")
		fmt.Printf("Error reading configuration file. Please check that %s/config.toml exists and syntax is correct.", confPath)
		os.Exit(1)
	}
	return config, nil
}

//Build the input source file
func runGoBuild(args []string, osys string, arch string, armv string) {
	cmdBuild := exec.Command("go", args...)
	cmdBuild.Env = append(os.Environ(),
		"GOARM="+armv,
		"GOOS="+osys,
		"GOARCH="+arch,
	)
	fmt.Printf("Cross-compiling for %s %s...", arch+armv, osys)
	if stdoutStderr, err := cmdBuild.CombinedOutput(); err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("\x1b[32;1mSuccess!\x1b[0m\n")
	}
}

//Copy binary to target with scp
func runSCP(args []string) {
	cmdSCP := exec.Command("scp", args...)
	fmt.Printf("Copying binary to target...")
	if stdoutStderr, err := cmdSCP.CombinedOutput(); err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("\x1b[32;1mSuccess!\x1b[0m\n")
	}
}

//Execute binary at target via SSH
func runSSH(args []string) {
	cmdSSH := exec.Command("ssh", args...) //TODO: password prompt not working, only key-based authentication works
	fmt.Printf("Running binary at target... Press ^C to stop.")
	if stdoutStderr, err := cmdSSH.CombinedOutput(); err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	}
}

func main() {
	config, missing := readConfig(); 
	if missing != nil {
		writeConfig(askConfig())
	}
	flagOsys := flag.String("os", config.GOOS, "Target operating system.")
	flagArch := flag.String("arch", config.GOARCH, "Target architecture.")
	flagArmv := flag.String("arm", config.GOARM, "ARM version.")
	flagTarg := flag.String("host", config.HOSTNAME, "Target IP or hostname.")
	flagUser := flag.String("user", config.USERNAME, "Username at target.")
	flagTdir := flag.String("dir", config.DIRECTORY, "Target directory.")
	flag.Parse()

	var sourcePath, localPath, packageName string
	buildArgs := []string{"build"}

	if flag.Arg(0) != "run" {
		fmt.Printf("\nSyntax: gore run <optional parameters> <optional path>\n\nAvailable parameters (set preferably in config.toml):\n")
		flag.Usage()
		os.Exit(1)
	}

	switch len(flag.Args()) { //Non-flag arguments, starts counting from AFTER last flag
	case 1:
		var err error
		if sourcePath, err = filepath.Abs(filepath.Dir(os.Args[0])); err != nil {
			fmt.Println("Could not read path.")
			os.Exit(1)
		}
		packageName = filepath.Base(sourcePath)
	case 2:
		sourcePath = flag.Arg(1)
		packageName = strings.TrimSuffix(filepath.Base(sourcePath), ".go")
		buildArgs = append(buildArgs, packageName)
		localPath = strings.TrimSuffix(sourcePath, filepath.Base(sourcePath))
		if !strings.HasSuffix(sourcePath, ".go") {
			fmt.Println("Absolute path needs to target a source file.")
			os.Exit(1)
		}
	default:
		fmt.Println("Check arguments.")
		os.Exit(1)
	}

	runGoBuild(buildArgs, *flagOsys, *flagArch, *flagArmv)
	runSCP([]string{localPath + packageName, *flagUser + "@" + *flagTarg + ":" + *flagTdir})
	runSSH([]string{"-t", "-t", *flagUser + "@" + *flagTarg, "'" + *flagTdir + packageName + "'"})
}
