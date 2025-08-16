package config

import (
	"encoding/json"
	"fmt"
	"os"
)

var fileName string = ".gatorconfig.json"

type Config struct {
	DBurl    string `json:"db_url"`
	Username string `json:"current_user_name"`
}

func (c Config) SetUser(username string) error {
	c.Username = username

	err := write(c)

	if err != nil {
		return err
	}

	return nil
}

func (c Config) DeleteUser() error {
	c.Username = ""

	err := write(c)

	if err != nil {
		return err
	}

	return nil
}

func (c Config) String() string {
	return fmt.Sprintf("{\n\tDBurl: %v\n\tUsername: %v\n}", c.DBurl, c.Username)
}

func Read() (Config, error) {
	filePath, err := getConfigFilePath()

	if err != nil {
		return Config{}, err
	}

	data, err := os.ReadFile(filePath)

	if err != nil {
		return Config{}, fmt.Errorf("func Read: error with reading config file, %v", err)
	}

	config := Config{}

	err = json.Unmarshal(data, &config)

	if err != nil {
		return Config{}, fmt.Errorf("func Read: error with marshalling config file, %v", err)
	}

	return config, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()

	if err != nil {
		return "", fmt.Errorf("func getConfigFilePath: error with UserHomeDir, %v", err)
	}

	return homeDir + "/" + fileName, nil
}

func write(c Config) error {
	jsonConfig, err := json.Marshal(c)

	if err != nil {
		return fmt.Errorf("func write: error with marhsalling, %v", err)
	}

	filePath, err := getConfigFilePath()

	if err != nil {
		return err
	}

	err = os.WriteFile(filePath, jsonConfig, 0644)

	if err != nil {
		return fmt.Errorf("func write: error with writing data, %v", err)
	}

	return nil
}
