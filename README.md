# gore
Go remote run, an all-in-one Linux tool for cross-compilation and execution of Go programs on a remote system.

Born out of engineer's laziness, `gore` seeks to combine three commands, `go build`, `scp`, and `ssh` with all their necessary arguments and environment variables into one goreous tool trying to mimic—and of course, improve—the `go run` command while targeting a remote system. Unlike `go run`, `gore run` can be run directly from your working directory without specifying a source file. It is useful e.g. when prototyping on a Raspberry Pi or similar board or headless system, where you might not want to set up a separate programming environment.

## Installation
`go get github.com/stuomas/gore`

## Syntax
`gore run <optional file path>`

When you run `gore` for the first time, it asks you to set up a configuration file interactively. The file is created in `$XDG_CONFIG_HOME/gore/config.toml`, where you can freely change the settings. You can force interactive re-configuration with `gore config`. For experimentation, you can also set the parameters as arguments before the run command:
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

## Notes
As of now, `gore` does not prefer to be bothered about passwords, so you should have key-based authentication set up to your remote system!