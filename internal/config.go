package internal

import (
	"io/ioutil"
	"log"
	"os"

	"gopkg.in/yaml.v2"
	"github.com/leighmacdonald/steamid/v3/steamid"
)

type Account struct {
	SteamId	uint64 `yaml:"steam_id"`
	SteamName	string `yaml:"steam_name"`
	DiscordId	*uint64 `yaml:"discord_id"`
	DiscordName	*string `yaml:"discord_name"`
}

type Config struct {
	Dsn               string          `yaml:"dsn"`
	NodeId	int64	`yaml:"node_id"`
	Mode *string `yaml:"mode"`
	DefaultAccounts *[]Account `yaml:"default_accounts"`
	RequireLoginToSubmit *bool `yaml:"require_user_to_be_logged_in_to_submit_reports"`
	SteamApiKey *string `yaml:"steam_api_key"`
	AlertChannelsWebhook *string `yaml:"alert_channels_webhook"`
}

func InitConfig() (cfg Config, err error) {
	config := Config{}
	log.Print("Loading config file ...")

	pwd, err := os.Getwd()
	if err != nil {
		log.Fatal("Failed to get current working directory.")
		return config, err
	}
	log.Println(pwd)
	config_bytes, err := ioutil.ReadFile(pwd + "/config.yml")
	if err != nil {
		log.Fatal("Failed to read config file.")
		return config, err
	}

	yaml.Unmarshal(config_bytes, &config)
	if config.SteamApiKey != nil {
		if err := steamid.SetKey(*config.SteamApiKey); err != nil {
			log.Fatal("Invalid steamid")
			os.Exit(1)
		}
	}
	return config, nil
}