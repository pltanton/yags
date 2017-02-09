package plugins

type Plugin interface {
	StartMonitor()
	StopMonitor()
	Chan() chan string
}
