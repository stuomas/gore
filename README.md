# gore
Go remote run, an all-in-one Linux command for cross-compilation and execution of Go programs on a remote machine.

Born out of an engineer's laziness, `gore` seeks to combine three commands, `go build`, `scp`, and `ssh` with all their necessary arguments and environment variables into one goreous tool trying to mimic the `go run` command while targeting a remote machine. Unlike `go run`, `gore run` can be run directly from your working directory without specifying a source file.

## Installation

`go get github.com/stuomas/gore`

## Usage
Set the parameters in the configuration file and command `gore run <optional parameters> <optional file path>`

If you want to experiment, you can also set the parameters as command arguments:
```
  -arch
    Target architecture.
  -arm
    ARM version.
  -dir
    Target directory.
  -host
    Target IP or hostname.
  -user
    Username at target.
  -os 
    Target operating system.
```