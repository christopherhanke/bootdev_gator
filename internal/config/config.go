package config

import (
	"encoding/json"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	DBURL           string `json:"db_url"`
	CurrentUserName string `json:"current_user_name"`
}

func (cfg *Config) SetUser(username string) error {

	cfg.CurrentUserName = username
	return write(*cfg)
}

func Read() (Config, error) {
	// read the JSON config file and return Config struct
	path, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}
	file, err := os.Open(path)
	if err != nil {
		return Config{}, err
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		return Config{}, err
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		return Config{}, err
	}
	return config, nil
}

func getConfigFilePath() (string, error) {
	//get Filepath and return Path as string
	home, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}

	return home + "/" + configFileName, nil
}

func write(cfg Config) error {
	//writing Config struct to JSON config file
	path, err := getConfigFilePath()
	if err != nil {
		return err
	}
	file, err := os.Create(path)
	if err != nil {
		return err
	}
	defer file.Close()

	if data, err := json.Marshal(cfg); err == nil {
		file.Write(data)
	} else {
		return err
	}
	return nil
}
