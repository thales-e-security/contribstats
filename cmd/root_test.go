// Copyright © 2018 NAME HERE <EMAIL ADDRESS>
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"testing"

	"github.com/pkg/errors"
	"github.com/thales-e-security/contribstats/pkg/config"
	"github.com/thales-e-security/contribstats/pkg/server"
)

type modkStatServer struct {
	server.StatServer
	wantErr bool
}

func (mss *modkStatServer) Start() (err error) {
	if mss.wantErr {
		return errors.New("Expected Error")
	}
	return
}

func TestExecute(t *testing.T) {
	tests := []struct {
		name    string
		wantErr bool
		debug   bool
		s       server.Server
	}{
		{
			name:    "OK",
			wantErr: false,
			debug:   true,
			s: &modkStatServer{
				wantErr: false,
			},
		}, {
			name:    "Error",
			wantErr: true,
			debug:   true,
			s: &modkStatServer{
				wantErr: true,
			},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {

			debug = tt.debug
			serverNewStatServer = func(huh config.Constants) (ss server.Server) {
				ss = tt.s

				return
			}
			Execute()
		})
	}
}
