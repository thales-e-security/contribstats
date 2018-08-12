package config

import (
	"github.com/mitchellh/go-homedir"
	"github.com/pkg/errors"
	"github.com/spf13/viper"
	"path/filepath"
	"testing"
)

func TestInitConfig(t *testing.T) {
	home, _ := homedir.Dir()
	type args struct {
		cfgFile string
	}
	tests := []struct {
		name    string
		args    args
		want    string
		homeErr bool
		readErr bool
		wantErr bool
	}{
		{
			name: "OK default",
			args: args{
				cfgFile: "",
			},
			want: filepath.Join(home, ".contribstats.yml"),
		},
		{
			name: "OK override",
			args: args{
				cfgFile: "/test/.contribstats.yml",
			},
			want: "/test/.contribstats.yml",
		}, {
			name: "Error HomeDir",
			args: args{
				cfgFile: "",
			},
			want:    "",
			homeErr: true,
			wantErr: true,
		}, {
			name: "Error ReadConfig",
			args: args{
				cfgFile: "",
			},
			want:    filepath.Join(home, ".contribstats.yml"),
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
			got, err := InitConfig(tt.args.cfgFile)
			if (err != nil) != tt.wantErr {
				t.Errorf("InitConfig() error = %v, wantErr %v", err, tt.wantErr)
				return
			}
			if got != tt.want {
				t.Errorf("InitConfig() = %v, want %v", got, tt.want)
			}
		})
	}
}
