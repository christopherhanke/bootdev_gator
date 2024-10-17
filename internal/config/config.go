package config

import (
	"encoding/json"
	"fmt"
	"io"
	"os"
)

const configFileName = ".gatorconfig.json"

type Config struct {
	Db_url            string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func Read() Config {
	// read the JSON config file and return Config struct
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Error get file path", err)
		return Config{}
	}
	file, err := os.Open(path)
	if err != nil {
		fmt.Println("Error open file", err)
		return Config{}
	}
	defer file.Close()
	data, err := io.ReadAll(file)
	if err != nil {
		fmt.Println("Error read file", err)
		return Config{}
	}
	var config Config
	err = json.Unmarshal(data, &config)
	if err != nil {
		fmt.Println("Error unmarshal", err)
		return Config{}
	}
	write(config)
	return config
}

func (cfg Config) SetUser(username string) {
	//writing Config struct to JSON config file after setting current_user_name
	path, err := getConfigFilePath()
	if err != nil {
		fmt.Println("Error get file path", err)
		return
	}
	file, err := os.Create(path)
	if err != nil {
		fmt.Println("Error open file", err)
		return
	}
	defer file.Close()

	cfg.Current_user_name = username

	if data, err := json.Marshal(cfg); err == nil {
		file.Write(data)
	} else {
		fmt.Println("Error writing file", err)
	}

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
	fmt.Println(cfg)
	return nil
}
