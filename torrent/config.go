package torrent

import (
	"github.com/anacrolix/torrent"
	"github.com/spf13/viper"
	"golang.org/x/time/rate"
	"torrent-fetcher/globals"
)

type Config struct {
	// Directory to store downloaded torrent data
	DataDir     string

	// Download traffic limit (in Megabytes)
	DownloadLim int

	// Port for running torrent client
	ListenPort  int
}

// create Config object for Torrent object
func newConfig() *torrent.ClientConfig {
	cfg := torrent.NewDefaultClientConfig()

	// setup limit as rate.Limit object
	if viper.GetInt(globals.KeyTorrentSpeedLimit) != 0 {
		toBytes := viper.GetInt(globals.KeyTorrentSpeedLimit) * 1024 * 1024

		bucket := toBytes / 16
		limit  := toBytes - bucket

		cfg.DownloadRateLimiter = rate.NewLimiter(rate.Limit(limit), bucket)
	}

	//cfg.Logger        = log.Discard
	cfg.DataDir       = viper.GetString(globals.KeyPathToDownloadData)
	cfg.ListenPort    = viper.GetInt(globals.KeyTorrentClientPort)
	cfg.HTTPUserAgent = "torrent-fetcher/1.0"
	cfg.UpnpID        = "golang/torrent-fetcher"
	cfg.Debug         = false
	cfg.NoUpload      = true

	return cfg
}
