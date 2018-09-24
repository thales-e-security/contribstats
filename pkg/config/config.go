package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/spf13/viper"
	"path/filepath"
	"strings"
)

var cfgName = ".contribstats"
var homedirDir = homedir.Dir
var readConfig = viper.ReadInConfig

//Config stores the viper config results after loading for layer passing around
type Config struct {
	Interval      int
	Token         string
	Cache         string
	Organizations []string
	Domains       []string
	Origins       []string
	Members       []string
	Blacklist     []string
}

//InitConfig reads in config file and ENV variables if set.
func InitConfig(in string) (constants Config, err error) {
	// Find home directory.
	var home string
	home, err = homedirDir()
	if err != nil {
		return
	}
	viper.SetConfigType("yaml")

	if in != "" {
		// Use config file from the flag.
		viper.SetConfigFile(in)
	} else {
		// Search config in home directory with name ".contribstats" (without extension).
		viper.AddConfigPath(home)
		viper.AddConfigPath(".")
		viper.AddConfigPath("/config")
		viper.SetConfigName(cfgName)

	}
	viper.SetEnvPrefix("CONTRIBSTATS")
	viper.AutomaticEnv() // read in environment variables that match

	// If a config file is found, read it in.
	if err := readConfig(); err == nil {

	} else {
		if _, ok := err.(viper.ConfigFileNotFoundError); ok {
			viper.SetDefault("domains", []string{"thalesesecurity.com", "thalesesec.net", "thales-e-security.com"})
			viper.SetDefault("organizations", []string{"unorepo"})
			viper.SetDefault("interval", 60)
			viper.SetDefault("origins", []string{"*"})
			viper.WriteConfigAs(filepath.Join(home, strings.Join([]string{cfgName, "yml"}, ".")))
		}
	}
	err = viper.Unmarshal(&constants)
	return
}
