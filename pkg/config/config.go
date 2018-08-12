package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var cfgName = ".contribstats"
var homedirDir = homedir.Dir
var readConfig = viper.ReadInConfig

//InitConfig reads in config file and ENV variables if set.
func InitConfig(cfgFile string) (string, error) {
	if cfgFile != "" {
		// Use config file from the flag.
		viper.SetConfigFile(cfgFile)
	} else {
		// Find home directory.
		home, err := homedirDir()
		if err != nil {
			return "", err
		}

		// Search config in home directory with name ".contribstats" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config")
		viper.SetConfigName(cfgName)
		cfgFile = filepath.Join(home, strings.Join([]string{cfgName, "yml"}, "."))
	}
	viper.SetEnvPrefix("CONTRIBSTATS")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := readConfig(); err == nil {
		//fmt.Println("Using config file:", viper.ConfigFileUsed())
	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetDefault("domains", []string{"thalesesecurity.com", "thalesesec.net", "thales-e-security.com"})
			viper.SetDefault("organizations", []string{"unorepo"})
			viper.SetDefault("interval", 60)
			viper.WriteConfigAs(cfgFile)
		} else {
			logrus.Info(err)

		}
	}
	return cfgFile, nil
}
