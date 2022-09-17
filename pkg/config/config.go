package config

import "github.com/ilyakaznacheev/cleanenv"

func NewConfig(fileName string, cfg interface{}) error {
	err := cleanenv.ReadConfig(fileName, cfg)
	if err != nil {
		return err
	}

	return nil
}
