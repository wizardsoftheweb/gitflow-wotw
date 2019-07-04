package main

import (
	"errors"
	"fmt"
	"os/exec"
	"strings"
	"sync"
	"time"

	"github.com/sirupsen/logrus"
)

func main() {
	fmt.Println("rad")
	BootstrapLogger(10)
	importChannel := make(chan bool)
	parseChannel := make(chan []string)
	storeChannel := make(chan []string)
	go ConfigGrabber(importChannel, parseChannel)
	go ConfigParser(parseChannel, storeChannel)
	go ConfigStorer(storeChannel)
	for {
		select {
		case importChannel <- true:
			logrus.Info("Kicking off data update")
		}

		time.Sleep(30000 * time.Millisecond)
	}
}

var ConfigCache ConfigStore

func shell(args ...string) (string, error) {
	process := exec.Command(args[0], args[1:]...)
	combined, err := process.CombinedOutput()
	if nil != err {
		logrus.Fatal(err)
	}
	return string(combined), err
}

func pollConfig(parseChannel chan []string) {
	result, _ := shell("git", "config", "--global", "--list")
	parseChannel <- strings.Split(result, "\n")
}

func ConfigGrabber(importChannel chan bool, parseChannel chan []string) {
	for {

		select {
		case importAlert := <-importChannel:

			if importAlert {
				logrus.Info("Grabbing new config")
				pollConfig(parseChannel)
			}
		}
		time.Sleep(10000 * time.Millisecond)
	}
}

func explodeLine(line string) []string {
	exploded := strings.Split(line, "=")
	value := strings.TrimSpace(strings.Join(exploded[1:], ""))
	exploded = strings.Split(strings.Join(exploded[:len(exploded)-1], ""), ".")
	section := strings.TrimSpace(exploded[0])
	key := strings.TrimSpace(exploded[len(exploded)-1])
	subsection := strings.TrimSpace(strings.Join(exploded[1:], ""))
	return []string{section, subsection, key, value}
}

func ConfigParser(parseChannel chan []string, storeChannel chan []string) {
	for {
		select {
		case freshConfig := <-parseChannel:
			logrus.Info("Parsing new config")
			for _, line := range freshConfig {
				storeChannel <- explodeLine(line)
			}
		}
		time.Sleep(10000 * time.Millisecond)
	}
}

func storeData(line []string) error {
	ConfigCache.CreateUpdate(Option{
		Section:    line[0],
		Subsection: line[1],
		Key:        line[2],
		Value:      line[3],
	})
	return nil
}

func ConfigStorer(storeChannel chan []string) {
	for {
		select {
		case freshData := <-storeChannel:
			logrus.Info("Storing new config")
			storeData(freshData)
		}
	}
}

type Option struct {
	Section    string
	Subsection string
	Key        string
	Value      string
}

func (o *Option) IdString() string {
	return strings.Join(
		[]string{
			o.Section,
			o.Subsection,
			o.Key,
		},
		" ",
	)
}

type ConfigStore struct {
	Options map[string]Option
	lock    sync.RWMutex
}

func GenerateKey(option Option) string {
	return option.IdString()
}

func (c *ConfigStore) CreateUpdate(option Option) {
	c.lock.Lock()
	defer c.lock.Unlock()
	if nil == c.Options {
		c.Options = make(map[string]Option)
	}
	key := GenerateKey(option)
	c.Options[key] = option
}

func (c *ConfigStore) Read(idString string) (*Option, error) {
	c.lock.RLock()
	defer c.lock.RUnlock()
	value, ok := c.Options[idString]
	if ok {
		return &value, nil
	}
	return nil, errors.New("Hash DNE")
}

func (c *ConfigStore) Delete(idString string) {
	c.lock.Lock()
	defer c.lock.Unlock()
	delete(c.Options, idString)
}
