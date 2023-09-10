package config

import "github.com/ilyakaznacheev/cleanenv"

type Config struct {
	BotToken  string `yaml:"botToken"`
	DebugMode bool   `yaml:"debugMode"`

	Database struct {
		Address        string `yaml:"address"`
		Database       string `yaml:"database"`
		User           string `yaml:"user"`
		Password       string `yaml:"password"`
		SSLMode        string `yaml:"sslMode"`
		InitMigrations bool   `yaml:"initMigrations"`
	} `yaml:"db"`
}

func NewConfig(fileName string, cfg interface{}) error {
	err := cleanenv.ReadConfig(fileName, cfg)
	if err != nil {
		return err
	}

	return nil
}
