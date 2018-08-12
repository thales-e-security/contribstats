package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"io/ioutil"
	"os"
	"testing"
)

func TestInitConfig(t *testing.T) {
	// Temp Files
	f, _ := ioutil.TempFile("", "")
	f.Write([]byte(`domains:
- thalesesecurity.com
- thalesesec.net
- thales-e-security.com
interval: 60
organizations:
- unorepo`))
	type args struct {
		cfgFile string
	}
	tests := []struct {
		name     string
		args     args
		override bool
		homeErr  bool
		readErr  bool
		wantErr  bool
	}{
		{
			name: "OK default",
			args: args{
				cfgFile: "",
			},
		},
		{
			name: "OK override",
			args: args{
				cfgFile: f.Name(),
			},
			override: true,
		}, {
			name: "Error HomeDir",
			args: args{
				cfgFile: "",
			},
			homeErr: true,
			wantErr: true,
		}, {
			name: "Error ReadConfig",
			args: args{
				cfgFile: "",
			},
			homeErr: false,
			readErr: true,
			wantErr: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.homeErr {
				homedirDir = func() (string, error) {
					return "", errors.New("expected error")
				}
			} else {
				homedirDir = homedir.Dir
			}
			if tt.readErr {
				readConfig = func() error {
					return viper.ConfigFileNotFoundError{}
				}
			} else {
				readConfig = viper.ReadInConfig
			}
			if tt.override {
				defer os.Remove(f.Name())
			}
			err := InitConfig(tt.args.cfgFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}

		})
	}
}
