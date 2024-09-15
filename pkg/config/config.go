package config

import (
	"fmt"

	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants"
	"github.com/Ritwiksrivastava0809/go-bank/pkg/constants/errorLogs"
	"github.com/rs/zerolog/log"
	"github.com/spf13/viper"
)

var config *viper.Viper

// Init initializes the configuration
func Init(env string) {
	var err error
	config = viper.New()

	config.SetConfigType(constants.DefaultConfigurationType)
	config.SetConfigName(env)                                // this should be the YAML filename without extension (e.g. "config" for config.yaml)
	config.AddConfigPath(constants.DefaultConfigurationPath) // Adjust the config path as needed

	err = config.ReadInConfig()
	if err != nil {
		log.Fatal().Msg(fmt.Sprintf(errorLogs.ParsingError, err.Error()))
	}
}

// DBConfig returns the DB configuration
func DBConfig() DB {
	return DB{
		Username: config.GetString("db.username"),
		Password: config.GetString("db.password"),
		Host:     config.GetString("db.host"),
		Port:     config.GetInt("db.port"),
		Name:     config.GetString("db.name"),
		SslMode:  config.GetString("db.sslmode"),
	}
}

func GetConfig() *viper.Viper {
	return config
}

func GetInternalToken() string {
	return config.GetString("token.internal")
}

func GetSymmetricKey() string {
	return config.GetString("token.symmetric")
}

func GetAccessTokenDuration() int {
	return config.GetInt("token.access.duration")
}
