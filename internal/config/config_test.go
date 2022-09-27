/*
 * Copyright © 2022 Durudex
 *
 * This file is part of Durudex: you can redistribute it and/or modify
 * it under the terms of the GNU Affero General Public License as
 * published by the Free Software Foundation, either version 3 of the
 * License, or (at your option) any later version.
 *
 * Durudex is distributed in the hope that it will be useful,
 * but WITHOUT ANY WARRANTY; without even the implied warranty of
 * MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE. See the
 * GNU Affero General Public License for more details.
 *
 * You should have received a copy of the GNU Affero General Public License
 * along with Durudex. If not, see <https://www.gnu.org/licenses/>.
 */

package config_test

import (
	"os"
	"reflect"
	"testing"
	"time"

	"github.com/durudex/durudex-auth-service/internal/config"
)

// Testing creating a new config.
func TestConfig_NewConfig(t *testing.T) {
	// Environment configurations.
	type env struct{ configPath, postgresUrl string }

	// Testing args.
	type args struct{ env env }

	// Set environments configurations.
	setEnv := func(env env) {
		os.Setenv("CONFIG_PATH", env.configPath)
		os.Setenv("POSTGRES_URL", env.postgresUrl)
	}

	// Tests structures.
	tests := []struct {
		name    string
		args    args
		want    *config.Config
		wantErr bool
	}{
		{
			name: "OK",
			args: args{env: env{
				configPath:  "fixtures/main",
				postgresUrl: "postgres://localhost:1",
			}},
			want: &config.Config{
				GRPC: config.GRPCConfig{
					Host: "auth.service.durudex.local",
					Port: "8001",
					TLS: config.TLSConfig{
						Enable: true,
						CACert: "./certs/rootCA.pem",
						Cert:   "./certs/auth.service.durudex.local-cert.pem",
						Key:    "./certs/auth.service.durudex.local-key.pem",
					},
				},
				Database: config.DatabaseConfig{
					Postgres: config.PostgresConfig{
						MaxConns: 20,
						MinConns: 5,
						URL:      "postgres://localhost:1",
					},
				},
				Auth: config.AuthConfig{
					Session: config.SessionConfig{TTL: time.Hour * 720},
					JWT:     config.JWTConfig{TTL: time.Minute * 15},
				},
				Service: config.ServiceConfig{
					User: config.Service{
						Addr: "user.service.durudex.local:8004",
						TLS: config.TLSConfig{
							Enable: true,
							CACert: "./certs/rootCA.pem",
							Cert:   "./certs/client-cert.pem",
							Key:    "./certs/client-key.pem",
						},
					},
					Code: config.Service{
						Addr: "code.service.durudex.local:8003",
						TLS: config.TLSConfig{
							Enable: true,
							CACert: "./certs/rootCA.pem",
							Cert:   "./certs/client-cert.pem",
							Key:    "./certs/client-key.pem",
						},
					},
					Email: config.Service{
						Addr: "email.service.durudex.local:8002",
						TLS: config.TLSConfig{
							Enable: true,
							CACert: "./certs/rootCA.pem",
							Cert:   "./certs/client-cert.pem",
							Key:    "./certs/client-key.pem",
						},
					},
				},
			},
		},
	}

	// Conducting tests in various structures.
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Set environments configurations.
			setEnv(tt.args.env)

			// Creating a new config.
			got, err := config.NewConfig()
			if (err != nil) != tt.wantErr {
				t.Errorf("error initialize config: %s", err.Error())
			}

			// Check for similarity of a config.
			if !reflect.DeepEqual(got, tt.want) {
				t.Errorf("error config are not similar")
			}
		})
	}
}
