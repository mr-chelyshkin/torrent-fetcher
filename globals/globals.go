package globals

var ConfigFlag = "configs"
var ConfigType = "yaml"

/*
	Configs keys.
	Using for getting configs from viper in all pkgs.
*/
var (
	KeyPathToDownloadData = "downloadTo"
	KeyPathToMoveComplete = "completeTo"
	KeyPathToWatchTorrent = "watchTo"

	KeyTriggerActionsList = "triggers"
	KeyTorrentFilePattern = "pattern"

	KeyTorrentSpeedLimit  = "limit_mb"
	KeyTorrentClientPort  = "port"
	KeyWatchDelayPeriod   = "delay_sec"
)
