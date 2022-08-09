// Handles ENV variables and config

package util

import (
	"github.com/spf13/viper"
)

// Config struct stores all configuration for the project
// the values are read by viper from a config file or from env variables
type Config struct {
	DBDriver string `mapstructure:"DB_DRIVER"`
	DBSource string `mapstructure:"DB_SOURCE"`
	ServerAddress string `mapstructure:"SERVER_ADDRESS"`
	MigrationURL   string  `mapstructure:"MIGRATION_URL"`
}

// Reads configs from config file or env variables
func LoadConfig(path string) (config Config, err error){
	// Default variables
	viper.SetDefault("DB_DRIVER", "postgres")
	viper.SetDefault("DB_SOURCE", "postgresql://root:secret@postgres:5432/BB-DB?sslmode=disable")
	viper.SetDefault("SERVER_ADDRESS", "0.0.0.0:8080")
	viper.SetDefault("MIGRATION_URL", "file://db/migration")

	// Viper reads config variables from env variables
	viper.AutomaticEnv() // Will automatically overwrite any variables already set with thier updated env var
	
	if (viper.Get("ENVIRONMENT") != "production"){
		// Viper reads config variables from config file
		viper.AddConfigPath(path)
		viper.SetConfigName("app") // file -> "app.env"
		viper.SetConfigType("env")


		// Begin reading variables from file
		err = viper.ReadInConfig()
		if err != nil {return}
	}

	err = viper.Unmarshal(&config)
	return
}