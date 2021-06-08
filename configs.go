package main

import (
	"errors"
	"fmt"
	"github.com/spf13/viper"
	"golang.org/x/sys/unix"
	"os"
	"regexp"
	"strconv"
	"strings"
	"torrent-fetcher/globals"
)

/*
	Predefine some configs.
	Use values if in configs file values is not set.
*/
var (
	predefineTorrentFilePattern = ".*\\.torrent$"
	predefineTriggerActionsList = "create"
	predefineTorrentClientPort  = 42069
	predefineTorrentSpeedLimit  = 0
	predefineWatchDelayPeriod   = 1
)

// from watcher pkg
var possibleActionList = []string{"create", "remove", "rename", "write", "move", "chmod"}

// check income configs data and return errors
func validateConfig() error {
	var errList []string

	// check paths
	if err := _checkPath(globals.KeyPathToDownloadData); err != nil {
		errList = append(errList, err.Error())
	}
	if err := _checkPath(globals.KeyPathToMoveComplete); err != nil {
		errList = append(errList, err.Error())
	}
	if err := _checkPath(globals.KeyPathToWatchTorrent); err != nil {
		errList = append(errList, err.Error())
	}

	// check integers values
	if err := _checkInt(globals.KeyWatchDelayPeriod, predefineWatchDelayPeriod); err != nil {
		errList = append(errList, err.Error())
	}
	if err := _checkInt(globals.KeyTorrentClientPort, predefineTorrentClientPort); err != nil {
		errList = append(errList, err.Error())
	}
	if err := _checkInt(globals.KeyTorrentSpeedLimit, predefineTorrentSpeedLimit); err != nil {
		errList = append(errList, err.Error())
	}

	// check customs
	if err := _checkTriggers(globals.KeyTriggerActionsList); err != nil {
		errList = append(errList, err.Error())
	}
	if err := _checkPattern(globals.KeyTorrentFilePattern); err != nil {
		errList = append(errList, err.Error())
	}

	// -- >
	if len(errList) > 0 {
		return fmt.Errorf(strings.Join(errList, "\n"))
	}
	return nil
}

/*
	Internals
*/

func _checkInt(key string, predefine int) error {
	if viper.GetString(key) == "" {
		viper.Set(key, predefine)
		return nil
	}
	if _, err := strconv.Atoi(viper.GetString(key)); err != nil {
		return errors.New("config key: " + key + " must be integer")
	}
	return nil
}

func _checkPattern(key string) error {
	if viper.GetString(key) == "" {
		viper.Set(key, predefineTorrentFilePattern)
	}
	_, err := regexp.Compile(viper.GetString(key))
	return err
}

func _checkPath(key string) error {
	path := viper.GetString(key)

	if path == "" {
		return errors.New("config key: " + key + " is not set, but required")
	}

	info, err := os.Stat(path)
	if err != nil {
		return os.MkdirAll(path, os.ModePerm)
	}
	if !info.IsDir() {
		return errors.New(path + " is not a directory")
	}
	if err := unix.Access(path, unix.W_OK); err != nil {
		return errors.New("user doesn't have permission to write to " + path)
	}
	return nil
}

func _checkTriggers(key string) error {
	var incorrect []string

	if len(viper.GetStringSlice(key)) == 0 {
		viper.Set(key, []string{predefineTriggerActionsList})
	}

	for _, action := range viper.GetStringSlice(key) {
		exist := false
		for _, possible := range possibleActionList {
			if action == possible {
				exist = true
				break
			}
		}
		if !exist {
			incorrect = append(incorrect, action)
		}
	}

	// -- >
	switch len(incorrect) {
	case 0:
		return nil
	case 1:
		return errors.New(incorrect[0] + " is incorrect action, allowed: " + strings.Join(possibleActionList, ","))
	default:
		return errors.New(strings.Join(incorrect, "," + " are incorrect actions, allowed: " + strings.Join(possibleActionList, ",")))
	}
}
