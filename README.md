YetAnotherGoStatus
==================

This program provides a simple configurable satatusline generator. It passes
the formatted status line each time when callbacks from plugging is received.

Unlike `conky`, `yags` prints satatusline only when status really changed and
not overloads disk with useless executions. By the way it is possible to
configure plugins to implement the conky behavior of execution any command
with constant pause.

Now it uses brand new golang's plugin system powered by `plugin` package. That
mean, that this program would only work with go version newer or equal 1.8.

## Installation

`go get github.com/pltanton/yags`

## Usage

If you use `conky` or `dzen` you could pass `yags` output to them through pipe,
for example:

```
yags /path/to/config.yml | dzen2 -x '866' -w '496' -ta 'r'

```

## Plugins

_Plugins_ is a go modules which provides callbacks when related to them action
is appears and provides formatted result of that action to `yags` output. There
are several implemented plugins:

- [x] [battery](https://github.com/pltanton/yags-battery) -- uses upower dbus
  messages to monitor battery device
- [x] [kbdd](https://github.com/pltanton/yags-kbdd) -- uses `kbdd` daemon to
  watch for keyboard layout
- [x] [timer](https://github.com/pltanton/yags-timer) -- conky like plugin to
  execute any shell command with pause it includes predefined plugins for
  **WiFi** and **time** monitoring
- [x] [maildir](https://github.com/pltanton/yags-maildir) -- monitors _new_
  maildir dirrectory changing for new mails
- [x] [volume](https://github.com/pltanton/yags-volume) -- uses alsalib to
  monitor volume changing and `pamixer` program to fetch an volume and a mute
  state, would be overwritten with pulselib in future
- [ ] cpu -- alias for timer for cpu monitoring
- [ ] ram -- alias for timer for ram monitoring

### How to use

After you download a plugin via `go get` or any other way, you should build it
as go plugin by `go build -buildmode=plugin` as documented in official golang
[docpage](https://tip.golang.org/pkg/plugin/). As a result of building you
should have a binary file ended by `.so` in working directory. Then you should
pass path to that file by `path` plugin's configuration attribute (you can see
it in [`config.example.yml`](https://github.com/pltanton/yags/blob/master/config.example.yml)).

### How to implement custom plugin

TBD.

## Configuration

You can configure `yags` by an any configuration file format, which
[viper](https://github.com/spf13/viper)
supports and pass configuration file path as a first argument to `yags`
command.

Exapmle of configuration you can find in root of this repository or in my
[dotfiles](https://github.com/pltanton/dotfiles/tree/master/config/yags)
repository.

At the root of configuration file it few basic configuration fields:

| Key     | Description                                                                                                                                                             | Default value |
| ---     | ---                                                                                                                                                                     | ---           |
| varSeps | symbols pair to wrap variables                                                                                                                                          | {}            |
| format  | string which would be formatted, you should use wrapped plugins names by `varSeps` to display plugin's output. Only if plugin passed to `format` it would be triggered. |               |
| plugins | subtree of configuration, where each plugin should be described                                                                                                         |               |

### Plugins

Each plugin in `plugin` section should contain `path` key to specify go's
plugin path to include.

`YAGS` will pass configuration subtree related to `plugin` into plugin, for
information about configuration plugin in specific plugin page.
