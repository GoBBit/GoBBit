# GoBBit

GoBBit - Reddit's like Forum software written in Golang

**This is a personal project and it is not finished yet, DO NOT USE ON PRODUCTION.**

### Installation

1. Clone this repository and rename the folder to GoBBit if it is not the name.
2. Install and configure MongoDB
3. Run the install script `install.sh` to install all the dependencies.
4. Compile using Makefile `make`
5. Configure Enviroment variables to configure the database and others params. (See: https://github.com/GoBBit/GoBBit/blob/master/server/api.go#L35 and https://github.com/GoBBit/GoBBit/blob/master/db/db.go#L26)

### Forum UI

This repository only contains the backend (based on an API), you must download or create a theme for your forum. A simple theme can be found here: https://github.com/GoBBit/GoBBit-SimpleTheme. It is in development and only contains basic functionality.

You should use nginx in order to serve static files of the UI, or if you are testing it you can use the `-static` option to choose the directory of the static files. Example:

`
bin/gobbit -static /var/html
`


### Why?

I started to develop GoBBit because I wanted to learn Go lang, so I decided to create my own forum's software :)
