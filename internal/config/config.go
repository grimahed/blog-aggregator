package config

import (
"os"
"encoding/json"
)

const configFileName = ".gatorconfig.json"
type Config struct {
	Db_url 		  string `json:"db_url"`
	Current_user_name string `json:"current_user_name"`
}

func (c *Config) ReadFile() (Config, error) {
	stuff := Config{}

	filepath, err := getConfigFilePath()
	if err != nil {
		return Config{}, err
	}

	bytes, err := os.ReadFile(filepath)
	if err != nil {
		return Config{}, err
	}

	if err := json.Unmarshal(bytes, &stuff); err != nil {
		return Config{}, err
	}

	return stuff, nil
}

func  (c *Config) SetUser(name string) (Config, error) {
	c.Current_user_name = name
	if err := write(*c); err != nil {
		return Config{}, err
	}
	return *c, nil
}

func getConfigFilePath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", err
	}
	return homeDir + "/" + configFileName, nil
}


func write(c Config) error {
	filepath, err := getConfigFilePath()
	if err != nil {
		return err
	}

	f, err := os.OpenFile(filepath, os.O_WRONLY|os.O_TRUNC, 0600)
	if err != nil {
		return err
	}
	defer f.Close()

	enc := json.NewEncoder(f)
	if err := enc.Encode(&c); err != nil {
		return err
	}
	return nil
}
