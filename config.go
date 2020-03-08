package pocket_manager

import (
	"encoding/json"
	"os"
)

type Config struct {
	Pocket struct {
		ConsumerKey string `json:"consumerKey"`
		AccessToken string `json:"accessToken"`
	} `json:"pocket"`
	Slack struct {
		PostUrl string `json:"postUrl"`
	} `json:"slack"`
}

func NewConfig() (*Config, error) {
	var config Config

	configFile, err := os.Open("./config.json")
	defer configFile.Close()
	if err != nil {
		return nil, err
	}

	jsonParser := json.NewDecoder(configFile)
	if err := jsonParser.Decode(&config); err != nil {
		return nil, err
	}

	return &config, nil
}
