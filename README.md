# gore
Go remote run, an all-in-one Linux command for cross-compilation and execution of Go programs on a remote machine.

Born out of an engineer's laziness, `gore` seeks to combine three commands, `go build`, `scp`, and `ssh` with all their necessary arguments and environment variables into a one goreous lump to imitate the `go run` command while targeting a remote machine. Myself I am using `gore` to code on my laptop and running the binary on a Raspberry Pi.

## Installation

`go get github.com/stuomas/gore`

## Usage
Set the parameters in the configuration file and command `gore /path/to/directory`
Optionally, set the parameters as arguments:
```
  -help
    Get help.
  -arch
    Target architecture.
  -dir
    Target directory.
  -os 
    Target operating system.
  -host
    Target IP or hostname.
  -user
    Username at target.
  -arm
    ARM version.
```