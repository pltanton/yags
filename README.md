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
- [x] volume -- uses alsalib to monitor volume changing and `pamixer` program
  to fetch an volume and a mute state, would be overwritten with pulselib in
  future
- [ ] network -- monitors network without networkmanager
- [ ] cpu -- alias for timer for cpu monitoring
- [ ] ram -- alias for timer for ram monitoring

## Configuration

TBD
