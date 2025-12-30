# crlancli
crlancli is a basic cli to interact with networked creality printers
# installation
first make sure you have [go](https://go.dev/dl/) installed, then run `go install github.com/axellse/crlancli`.
# contribute model names
if you have a creality printer and its model name is not added in the software, run `crlancli scan` and open an issue/pr with the model ID (eg. F003) and the model (eg. CR-10 SE)
# usage
run `crlancli scan` to view printers on the network and use the ip address shown with other commands. Also see `crlancli help` for a list of commands.
