# GoBBit

GoBBit - Reddit's like Forum software written in Golang

**This is a personal project and it is not finished yet, DO NOT USE ON PRODUCTION.**

### Installation

First of all install and configure Golang, then:

1. Clone this repository (in $GOPATH/src) and rename the folder to GoBBit if it is not the name.
2. Install and configure MongoDB
3. Run the install script `install.sh` to install all the dependencies.
4. Compile using Makefile. For example: `make install` will compile and generate a binary called GoBBit in **$GOPATH/bin**.
5. Configure the forum using the configure_example.json file and renaming to config.json. (See: https://github.com/GoBBit/GoBBit/blob/master/config/config.go).

Then, you can run using `-c` param to provide a configuration file if the name is not "config.json" or it isn't in the actual directory.
Example: `bin/gobbit -c path/to/config.json`

### Forum UI

This repository only contains the backend (based on an API), you must download or create a theme for your forum. A simple theme can be found here: https://github.com/GoBBit/GoBBit-SimpleTheme. It is in development and only contains basic functionality.

You should use nginx in order to serve static files of the UI, or if you are testing it you can use the `-static` option to choose the directory of the static files. Example:

`
bin/gobbit -static /var/html
`


### Why?

I started to develop GoBBit because I wanted to learn Go lang, so I decided to create my own forum's software :)
