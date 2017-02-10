package plugins

type Plugin interface {
	StartMonitor()
	Chan() chan string
}
