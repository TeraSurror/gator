package config

import (
	"encoding/json"
	"log"
	"os"
	"path/filepath"
)

const configfile = ".gatorconfig.json"

type Config struct {
	DbUrl           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(userName string) error {
	cfg.CurrentUserName = userName
	return write(*cfg)
}

func Read() (Config, error) {
	var config Config

	configFilePath, err := getConfigFilePath()
	if err != nil {
		log.Printf("Could not get config file: %v\n", err)
		return config, err
	}

	configFile, err := os.Open(configFilePath)
	if err != nil {
		log.Printf("Could not open config file: %v\n", err)
		return config, nil
	}

	jsonParser := json.NewDecoder(configFile)
	if err = jsonParser.Decode(&config); err != nil {
		log.Printf("Could not decode file: %v\n", err)
		return config, err
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Printf("cannot get home dir: %v\n", err)
		return "", err
	}

	path := filepath.Join(homeDir, configfile)
	return path, nil
}

func write(cfg Config) error {
	fullPath, err := getConfigFilePath()
	if err != nil {
		log.Printf("Could not get config file: %v\n", err)
		return err
	}

	file, err := os.Create(fullPath)
	if err != nil {
		log.Printf("Could not create config file: %v\n", err)
		return err
	}
	defer file.Close()

	encoder := json.NewEncoder(file)
	err = encoder.Encode(cfg)
	if err != nil {
		log.Printf("Could not encode Config to json file: %v\n", err)
		return err
	}

	return nil
}
