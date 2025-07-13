package Lightshow

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"os"
)

func loadConfiguration[N any](path string) (N, error) {
	configFile, err := os.Open(path)

	var config N
	if err != nil {
		return config, fmt.Errorf("config file for auto mode could not be accessed, %s", err)
	}

	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	jsonParser := json.NewDecoder(configFile)

	err = jsonParser.Decode(&config)
	if err != nil {
		return config, fmt.Errorf("invalid JSON format of auto mode config file, %s", err)
	}

	return config, nil
}

func storeConfiguration[N any](config *N, basePath string, backupPath string) {
	_ = os.Remove(backupPath)
	err := copyFile(basePath, backupPath)
	if err != nil {
		log.Fatalf("Backup of config file couldn't be creaed! %s", err)
	}
	_ = os.Remove(basePath)

	configFile, err := os.OpenFile(basePath, os.O_CREATE|os.O_WRONLY, os.ModePerm)

	if err != nil {
		log.Fatalf("Config file for auto mode could not be opened for writing! %s", err)
	}

	defer func(configFile *os.File) {
		_ = configFile.Close()
	}(configFile)

	jsonParser := json.NewEncoder(configFile)
	err = jsonParser.Encode(config)
	if err != nil {
		log.Fatalf("Configuration for auto mode could be JSON encoded! %s", err)
	}
}

func copyFile(src, dst string) error {
	sourceFileStat, err := os.Stat(src)
	if err != nil {
		return err
	}

	if !sourceFileStat.Mode().IsRegular() {
		return fmt.Errorf("%s is not a regular file", src)
	}

	source, err := os.Open(src)
	if err != nil {
		return err
	}
	defer func(source *os.File) {
		_ = source.Close()
	}(source)

	destination, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer func(destination *os.File) {
		_ = destination.Close()
	}(destination)
	_, err = io.Copy(destination, source)
	return err
}
