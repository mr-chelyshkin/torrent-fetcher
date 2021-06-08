package watcher

import (
	"fmt"
	"github.com/radovskyb/watcher"
	"github.com/spf13/viper"
	"strings"
	"torrent-fetcher/globals"
)

type Config struct {
	// Regexp pattern for math files while directories listening
	FilePattern    string

	// When torrent data was downloaded torrent file will put to thi directory, if empty the file will be deleted
	DoneDir        string

	// Delay between check actions in directories
	WatchDelay     int

	// Trigger by actions: prepared list in _actions variable
	TriggerActions []string

	// Path list for watching actions
	WatchPathList  []string

	// internal: transformed data from TriggerActions to []watcher.Op list object
	triggerActions []watcher.Op
}

// key/value matcher between Config.TriggerActions and watcher.Op
var _actions = map[string]watcher.Op{
	"create": watcher.Create,
	"remove": watcher.Remove,
	"rename": watcher.Rename,
	"write":  watcher.Write,
	"cmod":   watcher.Chmod,
	"move":   watcher.Move,
}

// create Config object for Watcher object
func newConfig() (*Config, error) {
	cfg := Config{
		WatchPathList:  []string{viper.GetString(globals.KeyPathToWatchTorrent)},
		TriggerActions: viper.GetStringSlice(globals.KeyTriggerActionsList),
		FilePattern:    viper.GetString(globals.KeyTorrentFilePattern),
		DoneDir:        viper.GetString(globals.KeyPathToMoveComplete),
		WatchDelay:     viper.GetInt(globals.KeyWatchDelayPeriod),
	}

	// convert objects between Config.TriggerActions and watcher.Op
	for _, action := range cfg.TriggerActions {
		if val, ok := _actions[strings.ToLower(action)]; ok {
			cfg.triggerActions = append(cfg.triggerActions, val)
		} else {
			return nil, fmt.Errorf("incorrect trigger action")
		}
	}

	return &cfg, nil
}
