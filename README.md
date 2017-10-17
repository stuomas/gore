# gore 0.1
go remote run, an all-in-one linux command for Go cross-compilation and execution on a remote machine.

Born out of an engineer's laziness, *gore* seeks to combine three commands, *go build*, *scp*, and *ssh* with all their necessary arguments into a one goreous lump to imitate the *go run* command while targeting a remote machine.
There might well be a lot better way to accomplish this, but no obvious search results after half a minute of googling was enough reason for me to do my own program.
Myself I am using *gore* to code on my laptop and cross-compiling and running the binary on a Raspberry Pi.

0.1 not usable without modifications on anybody else but me (unless you happen to have all parameters the same)