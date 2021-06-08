package main

import (
	"github.com/fsnotify/fsnotify"
	"github.com/spf13/pflag"
	"github.com/spf13/viper"
	"log"
	"os"
	"os/signal"
	"syscall"
	"torrent-fetcher/globals"
	"torrent-fetcher/torrent"
	"torrent-fetcher/watcher"
)

// -- >
func main() {
	pflag.String(globals.ConfigFlag, "", "config file path")
	pflag.Parse()

	if err := viper.BindPFlags(pflag.CommandLine); err != nil {
		log.Fatal(err)
	}

	viper.SetConfigFile(viper.GetString(globals.ConfigFlag))
	viper.SetConfigType(globals.ConfigType)

	if err := viper.ReadInConfig(); err != nil {
		log.Fatal(err)
	}
	if err := validateConfig(); err != nil {
		log.Fatal(err)
	}
	viper.WatchConfig()

	// -- >
	torrentObj, err := torrent.NewTorrent()
	if err != nil {
		log.Fatal(err)
	}

	watcherObj, err := watcher.NewWatcher()
	if err != nil {
		log.Fatal(err)
	}

	c := make(chan os.Signal)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)

	download := func(filePath string) {
		if err := torrentObj.DownloadFromFile(filePath); err != nil {
			log.Println(err)
		}
		if err := watcherObj.MoveFile(filePath); err != nil {
			log.Println(err)
		}
	}

	exit := func() {
		torrentObj.Close()
		watcherObj.Close()
		os.Exit(1)
	}

	// -- >
	viper.OnConfigChange(func(e fsnotify.Event) {exit()})
	go func() {<-c;exit()}()

	for _, file := range watcherObj.ExistFiles() {
		go download(file)
	}

	go func() {
		for {
			select {
			case event := <-watcherObj.Watch.Event:
				go download(event.Path)
			case err := <-watcherObj.Watch.Error:
				log.Fatal(err)
			}
		}
	}()
	log.Println("starting service")

	// start watcher
	if err := watcherObj.Run(); err != nil {
		log.Fatal(err)
	}
}
