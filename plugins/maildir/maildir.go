package maildir

import (
	"io/ioutil"
	"log"
	"path/filepath"
	"strconv"

	"github.com/howeyc/fsnotify"
	"github.com/spf13/viper"

	"github.com/pltanton/yags/utils"
)

// Maildir plugin structure
type Maildir struct {
	conf    *viper.Viper
	batName string
	out     chan string
}

// NewMaildir returns new instance of maildir plugin by given name
func NewMaildir(name string) Maildir {
	return Maildir{
		out:  make(chan string, 1),
		conf: viper.Sub("plugins." + name),
	}
}

// Chan returns a strings channel with maildir changing
func (m Maildir) Chan() chan string {
	return m.out
}

// StartMonitor starts monitoring for battery changing events
func (m Maildir) StartMonitor() {
	m.out <- m.formatMessage()
	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}
	defer watcher.Close()

	err = watcher.Watch(filepath.Join(m.conf.GetString("dir"), "new"))
	if err != nil {
		log.Fatal(err)
	}

	for {
		select {
		case <-watcher.Event:
			m.out <- m.formatMessage()
		case err := <-watcher.Error:
			log.Fatal(err)
		}
	}
}

// formatMessage formats message for printing
func (m Maildir) formatMessage() string {
	dir := m.conf.GetString("dir")
	files, err := ioutil.ReadDir(filepath.Join(dir, "new"))
	if err != nil {
		log.Fatal(err)
	}
	unreaded := len(files)
	var state string
	if unreaded == 0 {
		state = "empty"
	} else {
		state = "filled"
	}
	return utils.ReplaceVar(m.conf.GetString(state), "new", strconv.Itoa(unreaded))
}
