# pas

## Abstract
The idea behind this project is based on my personal requirements for an effective
financial administration tool. I have clear business processes which I need.
Currently, my financials are planned with [LibreOffice Calc](https://www.libreoffice.org/discover/calc/). This is fine for
most cases. However, the devil is in the detail. Sometimes I forgot to update
some cells or I overlooked some dependencies on my accounts based business processes.
This is really annoying. This software is currently in the state of a proof
of concept. I'll just test my ideas and best practices. This is not a full time
project, I work on it in my free time.

## Build
You need to have installed at least Go in version 1.11, because I use the go
modules feature for a better package management. Clone the repository and run
`go test ./...` inside. If all tests are fine you can run `go build .` However,
because I work test driven, the final binary does nothing currently. In a later
state I'll assemble the pieces to a working daemon.