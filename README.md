[![logo](https://raw.github.com/evilsocket/shellz/master/logo.png)](https://asciinema.org/a/203726)

[![Build](https://img.shields.io/travis/evilsocket/shellz/master.svg?style=flat-square)](https://travis-ci.org/evilsocket/shellz) 
[![Go Report Card](https://goreportcard.com/badge/github.com/evilsocket/shellz)](https://goreportcard.com/report/github.com/evilsocket/shellz) 
[![License](https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square)](/LICENSE) 
[![GoDoc](https://godoc.org/github.com/evilsocket/shellz?status.svg)](https://godoc.org/github.com/evilsocket/shellz) 
[![Release](https://img.shields.io/github/release/evilsocket/shellz.svg?style=flat-square)](https://github.com/evilsocket/shellz/releases/latest) 

`shellz` is a small utility to keep track of your connection credentials, servers and run commands on multiple machines at once. It supports `ssh`, `telnet` with more shell types coming soon!

**WORK IN PROGRESS**

## Install

Make sure you have a correctly configured **Go >= 1.8** environment, that `$GOPATH/bin` is in `$PATH` and then:

    $ go get github.com/evilsocket/shellz
    $ cd $GOPATH/src/github.com/evilsocket/shellz
    $ make && sudo make install

This command will download shellz, install its dependencies, compile it and move the `shellz` executable to `/usr/local/bin`.

## How to Use

The tool will use the `~/.shellz` folder to load your identities and shells json files, running the command `shellz` the first time will create the folder and the `idents` and `shells` subfolders for you. Once both `~/.shellz/idents` and `~/.shellz/shells` folders have been created, you can start by creating your first identity json file, for instance let's create `~/.shellz/idents/default.json` with the following contents:

```json
{
    "name": "default",
    "username": "evilsocket",
    "key": "~/.ssh/id_rsa"
}
```

As you can see my `default` identity is using my SSH private key to log in the `evilsocket` user, alternatively you can specify a `"password"` field instead of a `"key"`.

Now let's create our first shell json file ( `~/.shellz/shells/media.json` ) that will use the `default` identity we just created to connect to our home media server (called `media.server` in our example):

```json
{
    "name": "media-server",
    "host": "media.server",
    "port": 22,
    "identity": "default"
}
```

By default, shells are considered `ssh`, but also the `telnet` protocol is supported:

```sh
cat ~/.shellz/shells/tnas.json
```

```json
{
    "name": "tnas",
    "host": "tnas.local",
    "port": 23,
    "identity": "admin-tnas",
    "type": "telnet"
}
```

Once you have your shell and identity files ready, you can use shellz to run a command on a single machine:

    shellz -run uptime -on media-server

or on multiple at once:

    shellz -run uptime -on "media-server, tnas"

or on all of them:

    shellz -run uptime

## License

Shellz was made with ♥  by [Simone Margaritelli](https://www.evilsocket.net/) and it's released under the GPL 3 license.
