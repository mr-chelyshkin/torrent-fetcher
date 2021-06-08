package watcher

import (
	"github.com/radovskyb/watcher"
	"github.com/spf13/viper"
	"os"
	"path"
	"path/filepath"
	"regexp"
	"time"
	"torrent-fetcher/globals"
)

type Watcher struct {
	Watch *watcher.Watcher
}

// create Watcher object
func NewWatcher() (*Watcher, error) {
	cfg, err := newConfig()
	if err != nil {
		return nil, err
	}

	filePattern := regexp.MustCompile(cfg.FilePattern)
	watchObject := watcher.New()

	watchObject.SetMaxEvents(1)
	watchObject.FilterOps(cfg.triggerActions...)
	watchObject.AddFilterHook(watcher.RegexFilterHook(filePattern, false))

	for _, item := range cfg.WatchPathList {
		if err := watchObject.Add(item); err != nil {
			return nil, err
		}
	}

	return &Watcher{
		Watch: watchObject,
	}, nil
}

// find and return all files in watched directories match by pattern
func (d Watcher) ExistFiles() (res []string) {
	filePattern := regexp.MustCompile(viper.GetString(globals.KeyTorrentFilePattern))

	for filePath, file := range d.Watch.WatchedFiles() {
		if filePattern.MatchString(file.Name()) {
			res = append(res, filePath)
		}
	}

	return
}

// move torrent file to Config.DoneDir if torrent is downloaded
func (d Watcher) MoveFile(filePath string) error {
	if viper.GetString(globals.KeyPathToMoveComplete) == "" {
		return d.Watch.Remove(filePath)
	}
	return os.Rename(filePath, path.Join(viper.GetString(globals.KeyPathToMoveComplete), filepath.Base(filePath)))
}

// gracefully close Watcher object process
func (d Watcher) Close() {
	d.Watch.Close()
}

// start watching
func (d Watcher) Run() error {
	timeDelay := time.Duration(viper.GetInt(globals.KeyWatchDelayPeriod))
	if err := d.Watch.Start(time.Second * timeDelay); err != nil {
		return err
	}

	return nil
}