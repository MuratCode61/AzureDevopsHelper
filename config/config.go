package config

import (
	"azure_devops_helper/models"
	"log"
	"os"

	"github.com/spf13/viper"
)

var appConfig *models.AppConfig

func ReadAppConfig() {
	log.Println("ReadAppConfig is called.")

	workingdir, err := os.Getwd()
	if err != nil {
		log.Println("Error has been occurred while getting working dir.", err)
	}
	viper.AddConfigPath(workingdir)
	viper.SetConfigName("application")
	viper.SetConfigType("yaml")
	viper.ReadInConfig()

	err = viper.Unmarshal(&appConfig)
	if err != nil {
		log.Println("Error has been ocurred while unmarshalling config.", err)
	}
}

func GetAppConfig() *models.AppConfig {
	if appConfig == nil {
		ReadAppConfig()
	}
	return appConfig
}
