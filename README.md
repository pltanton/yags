YetAnotherGoStatus
==================

This program provides a simple configurable satatusline generator. It passes
the formatted status line each time when callbacks from plugging is received.

Unlike `conky`, `yags` prints satatusline only when status really changed and
not overloads disk with useless executions. By the way it is possible to
configure plugins to implement the conky behavior of execution any command
with constant pause.

## Installation

`go get github.com/pltanton/yags`

## Usage

If you use `conky` or `dzen` you could pass `yags` output to them through pipe,
for example:

```
yags /path/to/config.yml | dzen2 -x '866' -w '496' -ta 'r'

```

## Plugins

_Plugins_ is a subprograms which provides callbacks when related to them action
is appears and provides formatted result of that action to `yags` output. There
are several implemented plugins:

- [x] battery -- uses upower dbus messages to monitor battery device
- [x] kbdd -- uses `kbdd` daemon to watch for keyboard layout
- [x] timer -- conky like plugin to execute any shell command with pause
- [x] time -- alias for timer, but uses Go library to display current date/time
- [x] maildir -- monitors maildir changing for new mails
- [x] volume -- uses alsalib to monitor volume changing and `pamixer` program
  to fetch an volume and a mute state, would be overwritten with pulselib in
  future
- [x] network -- monitors network without networkmanager
- [ ] cpu -- alias for timer for cpu monitoring
- [ ] ram -- alias for timer for ram monitoring

## Configuration

You can configure `yags` by an any configuration file format, which
[viper](eat, multi-panel tabbed file manager for Linux desktops. Developed to provide a stable, efficient and highly customizable file manager, some of its features include: in-built VFS, HAL-based device manager, customizable menu system and bash integration.)
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

Each plugin in `plugin` section should contain `type` key to specifi plugin
type.

Some plugin has variables for interpolation and own formats. It acts just lkie
global `format`, but locally for plugin. You should first specify format for
plugin to insert it final output.

Available plugin types:

#### battery

Available configuration keys:

| Key    | Description                 | Default value |
| ---    | ---                         | ---           |
| name   | UPower name of battery      | BAT0          |
| high   | format for lvl > 75         | {lvl}         |
| medium | format for lvl > 35         | {lvl}         |
| low    | format for lvl > 12         | {lvl}         |
| empty  | format for lvl <= 12        | {lvl}         |
| ac     | format for conneted adapted | {lvl}         |

Available variables for interpolation

| Variable name | Description               |
| ---           | ---                       |
| lvl           | battery level in precents |

#### kbdd

TBD

#### timer

TBD

#### time

TBD

#### maildir

TBD

#### volume

TBD

#### network

TBD
