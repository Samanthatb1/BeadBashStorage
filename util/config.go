// Handles ENV variables and config

package util

import "github.com/spf13/viper"

// Config struct stores all configuration for the project
// the values are read by viper from a config file or from env variables
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
}

// Reads configs from config file or env variables
func LoadConfig(path string) (config Config, err error){
	// Viper reads config variables from config file
	viper.AddConfigPath(path)
	viper.SetConfigName("app") // file -> "app.env"
	viper.SetConfigType("env")

	// Viper reads config variables from env variables
	viper.AutomaticEnv() // Will automatically overwrite any variables already set with thier updated env var

	// Begin reading variables from file
	err = viper.ReadInConfig()
	if err != nil {return}

	err = viper.Unmarshal(&config)
	return
}