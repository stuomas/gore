package main

import (
	"bufio"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/BurntSushi/toml"
)

//Command flags
var flagOsys, flagArch, flagArmv, flagTarg, flagUser, flagTdir string

//ErrFileNotExist -> File doesn't exist
var ErrFileNotExist error = errors.New("file does not exist")

//ErrFileError -> Undefined problem with file
var ErrFileError error = errors.New("cannot read file")

//Configuration file structure
type Configuration struct {
	GOOS, GOARCH, GOARM, USERNAME, HOSTNAME, DIRECTORY string
}

func parseFlags(config Configuration) {
	flag.StringVar(&flagOsys, "os", config.GOOS, "Target operating system.")
	flag.StringVar(&flagArch, "arch", config.GOARCH, "Target architecture.")
	flag.StringVar(&flagArmv, "arm", config.GOARM, "ARM version.")
	flag.StringVar(&flagTarg, "host", config.HOSTNAME, "Target IP or hostname.")
	flag.StringVar(&flagUser, "user", config.USERNAME, "Username at target.")
	flag.StringVar(&flagTdir, "dir", config.DIRECTORY, "Target directory.")
	flag.Parse()
}

//Prompt user for filling in the config file interactively
func askConfig() []string {
	configName := []string{"GOOS", "GOARCH", "GOARM", "USERNAME", "HOSTNAME", "DIRECTORY"}
	var configValue []string
	fmt.Printf("Generate a configuration file by filling in the settings below:\n\n")
	prompt := bufio.NewReader(os.Stdin)
	for i := 0; i < len(configName); i++ {
		fmt.Print("  " + configName[i] + "=")
		value, _ := prompt.ReadString('\n')
		value = strings.Trim(strings.TrimRight(value, "\n"), "\"")
		configValue = append(configValue, "\""+value+"\"")
	}
	return configValue
}

//Write initial config file and folder
func writeConfig(defConf []string, confDir string) (string, error) {
	if err := os.MkdirAll(confDir+"/gore", 0755); err != nil {
		return "Error creating folder.", err
	}
	confSlice := []byte(fmt.Sprintf("GOOS=%s\nGOARCH=%s\nGOARM=%s\nUSERNAME=%s\nHOSTNAME=%s\nDIRECTORY=%s", defConf[0], defConf[1], defConf[2], defConf[3], defConf[4], defConf[5]))
	if err := ioutil.WriteFile(confDir+"/gore/config.toml", confSlice, 0755); err != nil {
		return "Error writing configuration file.", err
	}
	return fmt.Sprintf("\nConfiguration file created in %s/gore/config.toml. You can re-do the configuration with 'gore config'. Enjoy gore!\n\n", confDir), nil
}

//Read configuration file for environment variables and connection parameters
func readConfig(confDir string) (Configuration, error) {
	var config Configuration
	if _, err := os.Stat(confDir + "/gore/config.toml"); !os.IsNotExist(err) { //OK -> read
		if _, err := toml.DecodeFile(confDir+"/gore/config.toml", &config); err != nil {
			return config, ErrFileError
		}
	} else {
		return config, ErrFileNotExist //NOT OK -> write
	}
	return config, nil
}

//Build the input source file
func runGoBuild(args []string) {
	cmdBuild := exec.Command("go", args...)
	cmdBuild.Env = append(os.Environ(),
		"GOARM="+flagArmv,
		"GOOS="+flagOsys,
		"GOARCH="+flagArch,
	)
	fmt.Printf("Cross-compiling for %s %s...", flagArch+flagArmv, flagOsys)
	if stdoutStderr, err := cmdBuild.CombinedOutput(); err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	} else {
		fmt.Printf("\x1b[32;1mSuccess!\x1b[0m\n")
	}
}

//Copy binary to target with scp
func runSCP(args []string) {
	args = append(args, flagUser+"@"+flagTarg+":"+flagTdir)
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
func runSSH(args []string, packageName string) {
	args = append(args, flagUser+"@"+flagTarg, "'"+flagTdir+packageName+"'")
	cmdSSH := exec.Command("ssh", args...)
	fmt.Printf("Running binary at target... \x1b[32;1mPress ^C to stop.\x1b[0m")
	if stdoutStderr, err := cmdSSH.CombinedOutput(); err != nil {
		fmt.Printf("\x1b[31;1mFailed!\x1b[0m\n \u2937 %s", stdoutStderr)
		os.Exit(1)
	}
}

func main() {
	confDir := filepath.Join(os.Getenv("HOME"), ".config")
	config, fileErr := readConfig(confDir)
	if fileErr == ErrFileNotExist {
		if returnMsg, err := writeConfig(askConfig(), confDir); err != nil {
			fmt.Printf(returnMsg)
			os.Exit(1)
		} else {
			fmt.Printf(returnMsg)
			config, err = readConfig(confDir)
		}
	}
	if fileErr == ErrFileError {
		fmt.Println("Cannot read configuration file.")
		os.Exit(1)
	}

	parseFlags(config)

	var sourcePath, localPath, packageName string
	buildArgs := []string{"build"}

	if flag.Arg(0) == "config" && fileErr != ErrFileNotExist {
		returnMsg, _ := writeConfig(askConfig(), confDir)
		fmt.Printf(returnMsg)
		os.Exit(1)
	} else if flag.Arg(0) == "config" && fileErr == ErrFileNotExist {
		os.Exit(1)
	} else if flag.Arg(0) != "run" {
		fmt.Printf("\nSyntax: gore <optional parameters> run <optional path>\n	gore config to re-configure\nAvailable parameters:\n")
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

	runGoBuild(buildArgs)
	runSCP([]string{localPath + packageName})
	runSSH([]string{"-t", "-t"}, packageName)
}
